package semantic

import (
	"log"

	"github.com/samertm/chompy/parse"
)

// expects tree node
func Semantic(n parse.Node) string {
	s := NewStable()
	t, err := parse.AssertTree(n)
	if err != nil {
		log.Fatal(err)
	}
	// require a package statement
	if len(t.Kids) == 0 {
		log.Fatal("Requires a package statement")
	}
	return root(s, t)
}

// guaranteed to have at least one kid
func root(s *stable, t *parse.Tree) []string {
	msgs := make([]string, 0, 1)
	p, err := AssertPkg(t.Kids[0])
	if err != nil {
		return []string{"error: package statement must be first in file"}
	}
	// get all import statements
	i := 1
	for ; i < len(t.Kids); i++ {
		k := t.Kids[i]
		switch k.(type) {
		case *parse.Impts:
			// do something with imports
		default:
			break
		}
	}
	for i := 1; i < len(t.Kids); i++ {
		k := t.Kids[i]
		switch k.(type) {
		case *parse.Pkg:
			return "error: only one package statement per file"
		case *parse.Impts:
			return "error: imports must be declared before all other declarations"
		case *parse.Consts:
			consts(s, k.(*parse.Consts))
		case *parse.Typ:
			typ(s, k.(*parse.Typ))
		case *parse.Vars:
			vars(s, k.(*parse.Vars))

		}
	}
}
