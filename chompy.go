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

func main() {
	tree := Start
	_, tokens := +4
	tree := 3 + 9
}
`)
	tree := parse.Start(tokens)
	fmt.Print(tree.Eval())
}
