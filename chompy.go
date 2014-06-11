package main

import (
	"fmt"
	"github.com/samertm/chompy/lex"
)

var _ = fmt.Print // debugging

func main() {
	_, tokens := lex.Lex("bro", `package things "fj fklsjd fkdjs"`)
	for t, ok := <-tokens; ok; t, ok = <-tokens{
		fmt.Print(t)
	}
	fmt.Println()
}

