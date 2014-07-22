package semantic

import (
	"errors"

	"github.com/samertm/chompy/parse"
	"github.com/samertm/chompy/semantic/stable"
)

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
