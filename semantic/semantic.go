package semantic

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/samertm/chompy/parse"
	"github.com/samertm/chompy/semantic/stable"
)

var _ = log.Fatal   // debugging
var _ = fmt.Println // debugging

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
	for _, node := range t.Kids {
		switch n := node.(type) {
		case *parse.Funcdecl:
			name := n.Name.Name
			code = append(code, emitFuncHeader(name)...)
			code = append(code, emitFuncBody(n.Func.Body.Stmts)...)
			code = append(code, emitFuncReturn()...)
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

func emitFuncBody(stmts []parse.Node) []byte {
	table := stable.New(nil)
	var stackOffset int
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *parse.Vars:
			for _, v := range s.Vs {
				for _, id := range v.Idents {
					// Assume the type is an int
					stackOffset += 4
					t := &stable.Basic{Name: "int", Size: 4}
					table.Insert(id.Name, &stable.NodeInfo{T: t, Offset: stackOffset})
				}
			}
	}
	code := emitFuncStackSetup(stackOffset)
	for _, stmt := range stmts {
		
	}
	return code
}

func emitFuncStackSetup(offset int) []byte {
	return []byte("\tpush\t{r7}\n" +
		"\tsub\tsp, sp, #" + strconv.Itoa(offset) + "\n" +
		"\tadd\tr7, sp, #0\n")
}

func emitFuncReturn() []byte {
	return []byte("\tpop\t{r7}\n" +
		"\tbx\tlr\n")
}

func emitFuncHeader(name string) []byte {
	return []byte("\t.align\t2\n" +
		"\t.global\t" + name + "\n" +
		name + ":\n")
		
}

//func emitMov(dest int, src int)
