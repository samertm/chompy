// Defines nodes that are created when re-writing the AST
package semantic

import "github.com/samertm/chompy/parse"

type Decl struct {
	Up    Node
	I     *parse.Ident
	T     *parse.Typ
	Const bool
}

func (d *Decl) Up() Node {
	return d.Up
}

func (d *Decl) SetUp(n Node) {
	d.Up = n
}

func (d *Decl) Children(c chan<- parse.Node) {
	defer close(c)
	if d.I != nil {
		c <- d.I
	}
	if d.T != nil {
		c <- d.T
	}
}

func (d *Decl) Valid() bool {
	t := d.I != nil && d.I.Valid()
	if d.Const {
		return t && d.T != nil && d.T.Valid()
	}
	return t
}

func (d *Decl) String() string {
	s := "start decl\n"
	if d.Const {
		s += "const\n"
	}
	s += "ident: " + d.I.String()
	if d.T != nil {
		s += " type: " + d.T.String()
	}
	s += "\nend decl\n"
	return s
}

type Assign struct {
	Up Node
	I  *parse.Ident
	E  *parse.Expr
}

func (a *Assign) Up() Node {
	return a.Up
}

func (a *Assign) SetUp(n Node) {
	a.Up = n
}

func (a *Assign) Children(c chan<- parse.Node) {
	defer close(c)
	if a.I != nil {
		c <- a.I
	}
	if a.E != nil {
		c <- a.E
	}
}

func (a *Assign) Valid() bool {
	return a.I != nil && a.I.Valid() && a.E != nil && a.E.Valid()
}

func (a *Assign) String() string {
	return "start assign\nident: " + a.I.String() +
		" expr: " + a.E.String() + "end assign\n"
}
