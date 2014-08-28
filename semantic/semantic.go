package semantic

import (
	"errors"
	"fmt"
	"log"
	"reflect"
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
	code = append(code, "\tmov\tr0, #0\n"+
		"\tbl\tmain\n"+
		"\tmov\tr7, #1\n"+
		"\tswi\t#0\n"...)
	return code
}

func emitFuncBody(stmts []parse.Node) []byte {
	table := stable.New(nil)
	var code []byte
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
		case *parse.Assign:
			if len(s.LeftExpr) != len(s.RightExpr) {
				log.Fatal("args must match, and you can only have one argument on each side.")
			}
			// Hack: closure turns switch statement into an expression
			code = append(code, func(s *parse.Assign) []byte {
				switch s.Op {
				case "=":
					return emitFuncAssignment(table, s)
				}
				return []byte("")
			}(s)...)
		case *parse.ReturnStmt:
			if len(s.Exprs) == 0 {
				code = append(code, "\tmov\tr0, #0\n"...)
				code = append(code, emitFuncReturn()...)
				continue
			} else if len(s.Exprs) > 1 {
				log.Fatalf("I don't handle more than one return value: %s\n", s)
			}
			code = append(code, emitEvalExpr(table, s.Exprs[0])...)
			code = append(code, "\tmov\tr0, r6\n"...)
			code = append(code, emitFuncReturn()...)
		case *parse.Expr:
			log.Fatalf("here: %s", s)
		default:
			log.Fatalf("I don't handle %s yet\n", reflect.TypeOf(s))
		}
	}
	// Add stack setup to the beginning
	code = append(emitFuncStackSetup(stackOffset), code...)
	return code
}

func bprintf(format string, a ...interface{}) []byte {
	return []byte(fmt.Sprintf(format, a...))
}

func emitFuncAssignment(t *stable.Stable, a *parse.Assign) []byte {
	// First, we need to check to see that the expressions on the left are all idents
	// TODO: Make this work for more than one variable. [Issue: https://github.com/samertm/chompy/issues/3]
	if len(a.LeftExpr) == 0 {
		log.Fatal("Expected idents on the left of the assignment")
	}
	id, ok := a.LeftExpr[0].FirstN.Expr.(*parse.PrimaryE).Expr.(*parse.Ident)
	if !ok {
		log.Fatalf("Expected left of assignment to be ident: %s", a)
	}
	n, ok := t.Get(id.Name)
	if !ok {
		log.Fatalf("Ident %s not in scope", id)
	}
	// TODO: handle more than one expression [Issue: https://github.com/samertm/chompy/issues/5]
	if len(a.RightExpr) != 1 {
		log.Fatal("Expected one expression to the right of the assignment")
	}
	code := emitEvalExpr(t, a.RightExpr[0])
	// TODO: Evaluate the expression on the right. For now, we will assume that it [Issue: https://github.com/samertm/chompy/issues/2]
	// is 5.
	return append(code, bprintf("\tstr\tr6, [r7, #%d]\n", n.Offset)...)
}

func emitEvalExpr(t *stable.Stable, ex *parse.Expr) []byte {
	// TODO: do something with e.Op [Issue: https://github.com/samertm/chompy/issues/4]
	if ex.FirstN == nil {
		log.Fatalf("failed on %s", ex)
	}
	var exprs []*parse.Expr
	for e := ex; e != nil; e = e.SecondN {
		exprs = append([]*parse.Expr{e}, exprs...)
	}
	var result []byte
	for _, exp := range exprs {
		switch e := exp.FirstN.Expr.(type) {
		case *parse.PrimaryE:
			switch n := e.Expr.(type) {
			case *parse.Ident:
				ni, ok := t.Get(n.Name)
				if !ok {
					log.Fatalf("Ident %s not in scope", n)
				}
				result = append(result, bprintf("\tldr\tr5, [r7, #%d]\n", ni.Offset)...)
			case *parse.Lit:
				if n.Typ != "Int" {
					log.Fatal("The only type available are ints")
				}
				result = append(result, bprintf("\tmovs\tr5, #%s\n", n.Val)...)
			default:
				log.Fatalf("I don't handle %s yet\n", reflect.TypeOf(e.Expr))
			}
			switch exp.BinOp {
			case "+":
				result = append(result, "\tadd\tr6, r6, r5\n"...)
			case "":
				result = append(result, "\tmovs\tr6, r5\n"...)
			default:
				log.Fatalf("Unknown binop: %s", exp)
			}
		case *parse.UnaryE:
			log.Fatalf("Found unarye: %s", exp)
		default:
			log.Fatalf("Only deals with primary and unary errors: %s", exp)
		}
	}
	return result
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
