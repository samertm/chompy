// Defines nodes that are created when re-writing the AST
package semantic

import "github.com/samertm/chompy/parse"

type Decl struct {
	I     *parse.Ident
	T     *parse.Typ
	Const bool
	up    parse.Node
}

// TODO: fix these stumps
func (d *Decl) Replace(old, new Node) {
}

func (d *Decl) Up() parse.Node {
	return d.up
}

func (d *Decl) SetUp(n parse.Node) {
	d.up = n
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
	I  *parse.Ident
	E  *parse.Expr
	up parse.Node
}

// TODO fix stump
func (a *Assign) Replace(old, new Node) {
}

func (a *Assign) Up() parse.Node {
	return a.up
}

func (a *Assign) SetUp(n parse.Node) {
	a.up = n
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
