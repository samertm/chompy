package semantic

import (
	"fmt"

	"github.com/samertm/chompy/parse"
)

// rewriteTree runs all of the rewrite functions in order.
func treeWalks(t *parse.Tree) []string {
	rewrites := []func(*parse.Tree) []string{
		collectErrors,
		checkMain,
		// rewriteDeclAssign,
		// createStables,
	}
	for _, fn := range rewrites {
		s := fn(t)
		if len(s) != 0 {
			return s
		}
	}
	return nil
}

func checkMain(t *parse.Tree) []string {
	for _, kid := range t.Kids {
		switch f := kid.(type) {
		case *parse.Funcdecl:
			i, ok := f.Name.(*parse.Ident)

			if !ok {
				log.Fatal("Expected Ident")
			}
			
			if i.Name != "main" {
				continue
			}
			
			fn, ok := f.FuncOrSig.(*parse.Func)
			
			if !ok {
				log.Fatal("Expected Func");
			}
			
			sig := fn.Sig.(*parse.Sig)
			
			params := sig.Params.(*parse.Params)
			if len(params.Params) != 0 {
				log.Fatal("Expected no arguments");
			}
			
			if sig.Result != nil {
				log.Fatal("Expected no result");
			}
			
			return nil
		}
	}

	return []string{"Did not find main"}
}

// Collects all the errors starting at node from Erro nodes.
func collectErrors(t *parse.Tree) []string {
	// Preamble.
	errors := make([]string, 0, 5)
	nodes := make(chan parse.Node)
	go walkAll(t, node)
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

// // This phase adds the "up" pointer to every node in the tree.
// func addUp(t *parse.Tree) []string {
// 	// Preamble.
// 	var allFn walkFn = func(n parse.Node) bool {
// 		ch := make(chan parse.Node)
// 		go n.Children(ch)
// 		for kid := range ch {
// 			kid.SetUp(n)
// 		}
// 		return true
// 	}
// 	all := map[string]walkFn{
// 		"all": allFn,
// 	}
// 	walkAllHooks(t, nil, all)
// 	return nil
// }

// func createStables(t *parse.Tree) []string {
// 	// Let's initialize errors so we can report any we see
// 	errs := make([]string, 0)
// 	var treeFn walkFn = func(n parse.Node) bool {
// 		t := n
// 		t.RootStable = stable.New(nil)
// 		kids := make(chan parse.Node)
// 		go t.Children(kids)
// 		for kid := range kids {
// 			switch k := kid.(type) {
// 			case *parse.Pkg:
// 				break
// 			case *parse.Impts:
// 				break
// 			case *parse.Funcdecl:
// 				fmt.Println("FOUND", k)
// 			case *parse.Consts:

// 			case *parse.Vars:
// 			default:
// 				errs = append(errs, k.String())
// 			}
// 		}
// 		return true
// 	}
// 	hooks := map[string]walkFn {
// 		"*parse.Tree": treeFn,
// 	}
// 	return errs
// }
