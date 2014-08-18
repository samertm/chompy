package semantic

import (
	"log"

	"github.com/samertm/chompy/parse"
)

var _ = log.Fatal // debugging

func Gen(n parse.Node) []byte {
	t := check(n)
	return genCode(t)
}

// check is the "main" method for the semanic package. It runs all
// of the semanic checks and generates the IR for the backend.
func check(n parse.Node) *parse.Tree {
	t, ok := n.(*parse.Tree)
	if !ok {
		log.Fatal("Needed a tree.")
	}
	errs := treeWalks(t)
	if len(errs) != 0 {
		log.Fatal(errs)
	}
	return t
}

// Only deals with main.main right now.
func genCode(t *parse.Tree) []byte {
	code := emitStart()
	hooks := map[string]walkFn{
		"*parse.Funcdecl": func(n parse.Node) bool {
			f := n.(*parse.Funcdecl)
			name := f.Name.(*parse.Ident)
			code = append(code, emitFuncHeader(name.Name)...)
			if name.Name == "main" {
				code = append(code, emitFuncBody()...)
			}
			return false
		},
	}
	walkAllHooks(t, nil, hooks)
	return code
}

func emitStart() []byte {
	code := emitFuncHeader("_start")
	code = append(code, "\tbl\tmain\n" +
		"\tmov\tr0, #0\n" +
		"\tmov\tr7, #1\n" +
		"\tswi\t#0\n"...)
	return code
}

func emitFuncBody() []byte {
	return []byte("\tbx\tlr\n")
}

func emitFuncHeader(name string) []byte {
	return []byte("\t.align\t2\n"+
		"\t.global\t" + name + "\n" +
		name + ":\n")
}

//func emitMov(dest int, src int)
