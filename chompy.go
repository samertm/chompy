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
	_ "fmt"
	f "meow"
	. "cat"
)
`)
	tree := parse.Start(tokens)
	tree.Eval()
}
