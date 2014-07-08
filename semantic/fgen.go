// generates the three address code internal representation

package semantic

import "github.com/samertm/chompy/parse"

func Fgen(t parse.Node) string {
	// first: check that it's a tree
	tree, err := parse.AssertTree(t)
	if err != nil {
		return err.Error()
	}
	if tree.Kids[2].Valid() == false {
		f, _ := parse.AssertFuncdecl(tree.Kids[2])
		if f.FuncOrSig.Valid() == false {
			fun, err := parse.AssertFunc(f.FuncOrSig)
			if err != nil {
				return err.Error()
			}
			if fun.Sig.(*parse.Sig).Result == nil {
				return "HI"
			}
			return fun.String()
		}
		return f.String()
	}
	// second: check for a package statement
	return "yay :D"
}
