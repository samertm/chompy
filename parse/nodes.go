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
	t  string
	es []Node // expressions
}

func (c *cnst) Eval() (s string) {
	s += "start const spec\n"
	// subtle cisgendering
	s += c.is.Eval()
	s += "type: " + c.t + "\n"
	for _, e := range c.es {
		s += e.Eval()
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
