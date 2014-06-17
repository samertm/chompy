package parse

import (
	"fmt"
)

type Node interface {
	Eval() string
}

type grammarFn func(*parser) Node

type tree struct {
	kids []Node
}

func (t *tree) Eval() (s string) {
	for _, k := range t.kids {
		s += k.Eval()
	}
	return
}

type pkg struct {
	name string
}

func (p *pkg) Eval() string {
	return fmt.Sprintln("in package ", p.name)
}

type impts struct {
	imports []Node
}

func (i *impts) Eval() (s string) {
	s += fmt.Sprintln("start imports")
	for _, im := range i.imports {
		s += im.Eval()
	}
	s += fmt.Sprintln("end imports")
	return
}

type impt struct {
	pkgName  string
	imptName string
}

func (i *impt) Eval() string {
	return fmt.Sprintln("import: pkgName: " + i.pkgName + " imptName: " + i.imptName)
}

type erro struct {
	desc string
}

func (e *erro) Eval() string {
	return fmt.Sprintln("error: ", e.desc)
}

type consts struct {
	cs []Node // consts
}

func (c *consts) Eval() (s string) {
	s += "start const decl\n"
	for _, con := range c.cs {
		s += con.Eval()
	}
	s += "end const decl\n"
	return
}

// const
type cnst struct {
	is Node // idents
	t  Node
	es Node // expressions
}

func (c *cnst) Eval() (s string) {
	s += "start const spec\n"
	// subtle cisgendering
	s += c.is.Eval()
	if c.t != nil {
		s += c.t.Eval()
	}
	if c.es != nil {
		s += c.es.Eval()
	}
	s += "end const spec\n"
	return
}

type idents struct {
	is []string
}

func (i *idents) Eval() (s string) {
	for _, ident := range i.is {
		s += "ident: " + ident + "\n"
	}
	return
}

type lit struct {
	typ string
	val string
}

func (l *lit) Eval() string {
	return "lit: type: " + l.typ + " val: " + l.val + "\n"
}

type opName struct {
	id string
}

func (o *opName) Eval() string {
	return "opname: " + o.id + "\n"
}

type unaryE struct {
	op   string // Operand
	expr Node
}

func (u *unaryE) Eval() (s string) {
	s += "op: " + u.op + "\n"
	s += u.expr.Eval()
	return
}

// expression list
type exprs struct {
	es []Node
}

func (e *exprs) Eval() (s string) {
	for _, ex := range e.es {
		s += ex.Eval()
	}
	return
}

// expression list
type expr struct {
	binOp   string
	firstN  Node
	secondN Node
}

func (e *expr) Eval() (s string) {
	if e.binOp != "" {
		s += "binary_op: " + e.binOp + "\n"
	}
	if e.firstN != nil {
		s += e.firstN.Eval()
	}
	if e.secondN != nil {
		s += e.secondN.Eval()
	}
	return
}

type typ struct {
	t Node
}

func (t *typ) Eval() string {
	return "type: " + t.t.Eval() + "\n"
}

type ident struct {
	name string
}

func (i *ident) Eval() string {
	return i.name
}

type qualifiedIdent struct {
	pkg   string
	ident string
}

func (q *qualifiedIdent) Eval() string {
	return "pkg: " + q.pkg + " ident: " + q.ident
}

type types struct {
	typspecs []Node
}

func (t *types) Eval() (s string) {
	s += "start typedecl\n"
	for _, ty := range t.typspecs {
		s += ty.Eval()
	}
	s += "end typedecl\n"
	return
}

type typespec struct {
	i   Node //ident
	typ Node //type
}

func (t *typespec) Eval() (s string) {
	s += "start typespec\n"
	if t.i != nil {
		s += "ident: " + t.i.Eval() + "\n"
	}
	if t.typ != nil {
		s += t.typ.Eval()
	}
	s += "end typespec\n"
	return
}
