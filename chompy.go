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
	_, tokens := lex.Lex("bro", "hey bro")
	tree := parse.Start(tokens)
	fmt.Print(tree.Eval())
}
`)
	tree := parse.Start(tokens)
	fmt.Print(tree.Eval())
}
