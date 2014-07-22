package semantic

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/samertm/chompy/parse"
	"github.com/samertm/chompy/semantic/stable"
)

// Check is the "main" method for the semanic package. It runs all
// of the semanic checks and generates the IR for the backend.
func Check(n parse.Node) string {
	t, ok := n.(*parse.Tree)
	if !ok {
		log.Fatal("Needed a tree.")
	}
	if t.Valid() != true {
		errs := collectErrors(t)
		for _, s := range errs {
			fmt.Println(s)
		}
		log.Fatal("Tree is not valid ):")
	}
	errs := createStables(t)
	if len(errs) != 0 {
		log.Fatal(errs)
	}
	return "yay"
}

// Walks through all children
func walkAll(node parse.Node, kids chan<- parse.Node) {
	walkAllHooks(node, kids, nil)
}

// walkAllHooks will iterate through every single child for a node
// and thrown them down "kids", and does dispatch based on the type
// of the node from "hooks". Does a depth-first search of the tree.
// Closes kids when done.
//
// hooks is a map from strings (in the form "*parse.TYPE" to match
// the node types, which are all pointers and mostly from the package
// "parse") to functions. The function returns a bool indicating
// whether we should walk through the
//
// NOTE on the api, should it be (node, kids) or (kids, node)?
func walkAllHooks(node parse.Node, kids chan<- parse.Node, hooks map[string]func(parse.Node) bool) {
	// Preamble: set up first channel in stack.
	// chans is a stack. New channels are pushed and popped on
	// the right.
	defer close(kids)
	chans := make([]chan parse.Node, 0, 1)
	chans = append(chans, make(chan parse.Node))
	go node.Children(chans[0])
	for len(chans) != 0 {
		next, ok := <-chans[len(chans)-1]
		if !ok {
			// The channel has closed, pop it off the
			// stack.
			chans = chans[:len(chans)-1]
			continue
		}
		// next is a node. Create a new channel to recieve
		// its children and push it into the stack.
		if hooks == nil {
			goto APPEND
		} else {
			// We need to look inside hooks and dispatch
			// based on the type of next. I assume that
			// reflect is expensive, so we don't take
			// this path if hooks is nil.
			typ := reflect.TypeOf(node).String()
			fn, ok := hooks[typ]
			if !ok {
				goto APPEND
			} else {
				val := fn(next)
				if val {
					goto APPEND
				}
				goto NEXT
			}
		}
	APPEND:
		chans = append(chans, make(chan parse.Node))
		go next.Children(chans[len(chans)-1])
	NEXT:
		// Pass next to kids.
		kids <- next
	}
}

// Collects all the errors starting at node from Erro nodes.
func collectErrors(node parse.Node) []string {
	// Preamble.
	errors := make([]string, 0, 5)
	nodes := make(chan parse.Node)
	go walkAll(node, nodes)
	// Collect errors.
	for n := range nodes {
		switch n.(type) {
		case *parse.Erro:
			e := n.(*parse.Erro)
			errors = append(errors, e.Desc)
		default:
			if n.Valid() == false {
				fmt.Println("INVALID:", n)
			}
		}
	}
	return errors
}

// createType creates a stable.Type from n. But you already knew that
// from the function signature, didn't you ;D
// NOTE this may be better off in the package stable, but I can't put
// it there because it accepts a parse.Node.
func createType(node parse.Node) (stable.Type, error) {
	if node == nil {
		return nil, errors.New("Recieved nil type")
	}
	switch n := node.(type) {
	case *parse.Typ:
		return createType(n.T)
	case *parse.Ident:
		return &stable.Basic{Pkg: "blank", Name: n.Name}, nil
	case *parse.QualifiedIdent:
		return &stable.Basic{Pkg: n.Pkg, Name: n.Ident}, nil
	case *parse.Cnst:
		return createType(n.T)
	case *parse.Lit:
		return &stable.Basic{Pkg: "blank", Name: n.Typ}, nil
	case *parse.Expr:
		f, err := createType(n.FirstN)
		if err != nil {
			return nil, err
		}
		s, err := createType(n.SecondN)
		if err != nil {
			return nil, err
		}
		if !f.Equal(s) {
			return nil, typeMismatch(f, s)
		}
		// Because the types match, we can return either one
		return s, nil
	case *parse.UnaryE:
		return createType(n.Expr)
	case *parse.PrimaryE:
		ex, err := createType(n.Expr)
		if err != nil {
			return nil, err
		}
		// If the expression has no prime (i.n. it does not
		// continue), then we can return ex.
		if n.Prime == nil {
			return ex, nil
		}
		// Otherwise, we need to check to see that prime has
		// the same type as ex.
		prime, err := createType(n.Prime)
		if err != nil {
			return nil, err
		}
		if !ex.Equal(prime) {
			return nil, typeMismatch(ex, prime)
		}
	case *parse.Typespec:
		return createType(n.Typ)
	case *parse.Funcdecl:
		// NOTE Might be able to break this into another case
		// statement (so that most of it gets handled by, say
		// case *parse.Func)
		fn := &stable.Func{}
		t, err := createType(n.Name)
		if err != nil {
			return nil, err
		}
		name, ok := t.(*stable.Basic)
		if !ok {
			return nil, errors.New("Expected a basic type")
		}
		fn.Name = name
		// We need to get the function signature so we can
		// iterate over it.
		var sig *parse.Sig
		switch s := n.FuncOrSig.(type) {
		case *parse.Func:
			sig, ok = s.Sig.(*parse.Sig)
			if !ok {
				return nil, errors.New("Expected signature")
			}
		case *parse.Sig:
			// We may have set sig in the previous block
			sig = s
		default:
			return nil, errors.New("Expected signature")
		}
		// Go through sig's params and create types for them
		p, ok := sig.Params.(*parse.Params)
		if !ok {
			return nil, errors.New("Expected params")
		}
		// Go through the params and turn them into types to
		// appends to args.
		args, err := makeTypes(p.Params)
		if err != nil {
			return nil, err
		}
		fn.Args = args
		// Now, get the result. The result might be params,
		// so we need to check for it manually.
		var result []stable.Type
		switch i := sig.Result.(type) {
		case *parse.Params:
			result, err = makeTypes(i.Params)
			if err != nil {
				return nil, err
			}
		case *parse.Typ:
			t, err := createType(i)
			if err != nil {
				return nil, err
			}
			result = []stable.Type{t}
		}
		fn.Result = result
		return fn, nil
	default:
		return nil, errors.New("Node has no type " + n.String())
	}
	return nil, errors.New("This should never happen")
}

// For use with types that hold multiple nodes, like Params.
func makeTypes(nodes []parse.Node) ([]stable.Type, error) {
	types := make([]stable.Type, 0)
	for _, n := range nodes {
		t, err := createType(n)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, nil
}

// Creates an error from any number of mismatching types.
func typeMismatch(types ...stable.Type) error {
	s := "Types do not match: "
	for _, t := range types {
		s += t.String() + "\n"
	}
	return errors.New(s)
}

func createStables(t *parse.Tree) []string {
	// Let's initialize errors so we can report any we see
	errs := make([]string, 0)
	// First, let's create an Stable to hold the information
	// about the root tree's children.
	t.RootStable = stable.New(nil)
	// We're going to iterate through t's children and add them
	// to rootStable.
	kids := make(chan parse.Node)
	go t.Children(kids)
	for kid := range kids {
		switch k := kid.(type) {
		case *parse.Pkg:
			break
		case *parse.Impts:
			break
		case *parse.Funcdecl:
			fmt.Println("FOUND", k)
		case *parse.Consts:
			fmt.Println("FOUND", k)
		case *parse.Vars:
			// for _, v := range k.Vs {
			// }
		default:
			errs = append(errs, k.String())
		}
	}
	return errs
}

// This phase rewrites Varspec, (add nodes) into Assign and Decl
// nodes. I'm not sure if I should name it rewrite01 (to show the
// order these rewrites should be done in)...
func rewriteDeclAssign(t *parse.Tree) []string {
	// stumpyyyyyyyyyy
	return make([]string, 0)
}
