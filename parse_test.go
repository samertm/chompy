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
	`package main

import (
	"fmt"
)

const ribs
`,
	`package main

import (
	"fmt"
)

const ribs, mibs
`,
	`package main

import (
	"fmt"
)

const (
	ribs, mibs
	tibs
)
`,
	`package main

import "fmt"

const ribs = 4
`,
	`package main

import "fmt"

const (
	ribs = 4
	mibs = "hi there"
)
`,
	`package main

import "fmt"

const ribs = 4 + 4
`,
	`package main

import "fmt"

const ribs int = 4 + 4
`,
	`package main

import "fmt"

const ribs fmt.Int = 4 + 4
`,
	`package main

import "fmt"

type thangs int
`,
	`package main

import "fmt"

var meow int
`,
	`package main

import "fmt"

var meow int = 4 + 4
`,
	`package main

import "fmt"

func meow() {
	var meow int = 4 + 4
}
`,
	`package main

import "fmt"

func meow(a thing, b otherthing) string {
	var meow int = 4 + 4
}
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
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start const decl
start const spec
ident: ribs
end const spec
end const decl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start const decl
start const spec
ident: ribs
ident: mibs
end const spec
end const decl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start const decl
start const spec
ident: ribs
ident: mibs
end const spec
start const spec
ident: tibs
end const spec
end const decl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start const decl
start const spec
ident: ribs
op: 
lit: type: Int val: 4
end const spec
end const decl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start const decl
start const spec
ident: ribs
op: 
lit: type: Int val: 4
end const spec
start const spec
ident: mibs
op: 
lit: type: String val: hi there
end const spec
end const decl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start const decl
start const spec
ident: ribs
binary_op: +
op: 
lit: type: Int val: 4
op: 
lit: type: Int val: 4
end const spec
end const decl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start const decl
start const spec
ident: ribs
type: int
binary_op: +
op: 
lit: type: Int val: 4
op: 
lit: type: Int val: 4
end const spec
end const decl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start const decl
start const spec
ident: ribs
type: pkg: fmt ident: Int
binary_op: +
op: 
lit: type: Int val: 4
op: 
lit: type: Int val: 4
end const spec
end const decl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start typedecl
start typespec
ident: thangs
type: int
end typespec
end typedecl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start vardecl
start varspec
ident: meow
type: int
end varspec
end vardecl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start vardecl
start varspec
ident: meow
type: int
binary_op: +
op: 
lit: type: Int val: 4
op: 
lit: type: Int val: 4
end varspec
end vardecl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start funcdecl
ident: meow
start block
start vardecl
start varspec
ident: meow
type: int
binary_op: +
op: 
lit: type: Int val: 4
op: 
lit: type: Int val: 4
end varspec
end vardecl
end block
end funcdecl
`,
	`in package  main
start imports
import: pkgName:  imptName: fmt
end imports
start funcdecl
ident: meow
start parameters
start parameterdecl
ident: a
type: thing
end parameterdecl
start parameterdecl
ident: b
type: otherthing
end parameterdecl
end parameters
start result
type: string
end result
start block
start vardecl
start varspec
ident: meow
type: int
binary_op: +
op: 
lit: type: Int val: 4
op: 
lit: type: Int val: 4
end varspec
end vardecl
end block
end funcdecl
`,
}

func TestParse(t *testing.T) {
	if len(inputs) != len(outputs) {
		t.Errorf("len(inputs) != len(outputs)")
		return
	}
	for i, _ := range inputs {
		_, tokens := lex.Lex("bro", inputs[i])
		tree := parse.Start(tokens)
		result := tree.Eval()
		if outputs[i] != result {
			t.Errorf("\n---expected---\n%s---recieved---\n%s---end---", outputs[i], result)
		}
	}
}
