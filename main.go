package main

import (
	"fmt"

	"github.com/samertm/chompy/lex"
	"github.com/samertm/chompy/parse"
)

var _ = fmt.Print // debugging

func main() {
	_, tokens := lex.Lex("bro", `
package main

import (
	"fmt"
	"github.com/samertm/chompy/lex"
	"github.com/samertm/chompy/parse"
)

var mex = 4

func main() {
}
`)
	tree := parse.Start(tokens)
	fmt.Println(tree)
}
