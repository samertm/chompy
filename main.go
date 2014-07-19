package main

import (
	"fmt"

	"github.com/samertm/chompy/lex"
	"github.com/samertm/chompy/parse"
	"github.com/samertm/chompy/semantic"
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

var _ = fmt.Print // debugging

func main() {
	tree := parse.Start(tokens)
	fmt.Print(tree)
}
`)
	tree := parse.Start(tokens)
	semantic.Check(tree)
}
