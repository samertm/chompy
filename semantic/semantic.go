package semantic

import (
	"errors"
	"log"

	"github.com/samertm/chompy/parse"
)

var _ = log.Fatal // debugging

type sErrors []string

func (e sErrors) Error() string {
	if len(e) == 0 {
		return "No errors."
	}
	str := make([]byte, 0)
	for i, err := range e {
		if i != 0 {
			str = append(str, '\n')
		}
		str = append(str, err...)
	}
	return string(str)
}

func Gen(n parse.Node) ([]byte, error) {
	t, err := check(n)
	if err != nil {
		return nil, err
	}
	return genCode(t), nil
}

// check is the "main" method for the semanic package. It runs all
// of the semanic checks and generates the IR for the backend.
func check(n parse.Node) (*parse.Tree, error) {
	t, ok := n.(*parse.Tree)
	if !ok {
		return nil, errors.New("Needed a tree.")
	}
	errs := treeWalks(t)
	if len(errs) != 0 {
		return nil, errs
	}
	return t, nil
}

// Only deals with main.main right now.
func genCode(t *parse.Tree) []byte {
	code := emitStart()
	for _, n := range t.Kids {
		var f *parse.Funcdecl
		var ok bool
		if f, ok = n.(*parse.Funcdecl); !ok {
			continue
		}

		name := f.Name.Name
		code = append(code, emitFuncHeader(name)...)
		if name == "main" {
			code = append(code, emitFuncBody()...)
		}
	}
	return code
}

func emitStart() []byte {
	code := emitFuncHeader("_start")
	code = append(code, "\tbl\tmain\n"+
		"\tmov\tr0, #0\n"+
		"\tmov\tr7, #1\n"+
		"\tswi\t#0\n"...)
	return code
}

func emitFuncBody() []byte {
	return []byte("\tbx\tlr\n")
}

func emitFuncHeader(name string) []byte {
	return []byte("\t.align\t2\n" +
		"\t.global\t" + name + "\n" +
		name + ":\n")
}

//func emitMov(dest int, src int)
