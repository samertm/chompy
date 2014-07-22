package semantic

import (
	"fmt"

	"github.com/samertm/chompy/parse"
	"github.com/samertm/chompy/semantic/stable"
)

// rewriteTree runs all of the rewrite functions in order.
func treeWalks(t *parse.Tree) []string {
	rewrites := []func(*parse.Tree) []string{
		collectErrors,
		rewriteDeclAssign,
		createStables,
	}
	for _, fn := range rewrites {
		s := fn(t)
		if len(s) != 0 {
			return s
		}
	}
	return nil
}

// Collects all the errors starting at node from Erro nodes.
func collectErrors(t *parse.Tree) []string {
	// Preamble.
	errors := make([]string, 0, 5)
	nodes := make(chan parse.Node)
	go walkAll(t, nodes)
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

// This phase adds the "up" pointer to every node in the tree.
func addUp(t *parse.Tree) []string {
	// stumpyyyyyyyyyy
	return make([]string, 0)
}

// This phase rewrites Varspec, (add nodes) into Assign and Decl
// nodes. I'm not sure if I should name it rewrite01 (to show the
// order these rewrites should be done in)...
func rewriteDeclAssign(t *parse.Tree) []string {
	// stumpyyyyyyyyyy
	return make([]string, 0)
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
