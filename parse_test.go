package main

import (
	"github.com/samertm/chompy/lex"
	"github.com/samertm/chompy/parse"
	"testing"
)

var inputs = []string{
	`
package main

import (
	_ "fmt"
	f "meow"
	. "cat"
)
`,
}

var outputs = []string{
	`in package  main
start imports
import: pkgName: _ imptName: fmt
import: pkgName: f imptName: meow
import: pkgName: . imptName: cat
end imports
`,
}

func TestParse(t *testing.T) {
	for i, _ := range inputs {
		_, tokens := lex.Lex("bro", inputs[i])
		tree := parse.Start(tokens)
		result := tree.Eval()
		if outputs[i] != result {
			t.Errorf("---expected---\n%s\n---recieved---\n%s", outputs[i], result)
		}
	}
}
