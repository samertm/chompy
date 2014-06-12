package main

import (
	"fmt"
	"github.com/samertm/chompy/lex"
)

var _ = fmt.Print // debugging

func main() {
	_, tokens := lex.Lex("bro", `
package main

import (
	"fmt"
	"github.com/samertm/chompy/lex"
)

var _ = fmt.Print // debugging

func main() {
	_, tokens := lex.Lex("bro", "string plz")
	for t, ok := <-tokens; ok; t, ok = <-tokens{
		fmt.Print(t)
	}
	fmt.Println()
}
package main

import "yo"

func bro() {
	brotato := myManShawn{}
}
 `)
	for t, ok := <-tokens; ok; t, ok = <-tokens{
		fmt.Print(t)
	}
	fmt.Println()
}

