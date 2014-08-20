package semantic

import "github.com/samertm/chompy/parse"

// rewriteTree runs all of the rewrite functions in order.
func treeWalks(t *parse.Tree) []string {
	rewrites := []func(*parse.Tree) []string{
		checkMain,
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
			if f.Name.Name == "main" {
				if f.Func.Sig.Params == nil &&
					f.Func.Sig.Result == nil {
					return nil
				}
			}
		}
	}
	return []string{"Did not find main"}
}
