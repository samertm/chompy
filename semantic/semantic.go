package semantic

import (
	"log"

	"github.com/samertm/chompy/parse"
)

// Check is the "main" method for the semanic package. It runs all
// of the semanic checks and generates the IR for the backend.
func Check(n parse.Node) string {
	t, ok := n.(*parse.Tree)
	if !ok {
		log.Fatal("Needed a tree.")
	}
	errs := treeWalks(t)
	if len(errs) != 0 {
		log.Fatal(errs)
	}
	return "yay"
}
