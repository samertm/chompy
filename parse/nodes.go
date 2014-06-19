package parse

import (
	"fmt"
)

type Node interface {
	Eval() string
}

type grammarFn func(*parser) Node

type Tree struct {
	Kids []Node
}

func (t *Tree) Eval() (s string) {
	for _, k := range t.Kids {
		s += k.Eval()
	}
	return
}

type Pkg struct {
	Name string
}

func (p *Pkg) Eval() string {
	return fmt.Sprintln("in package ", p.Name)
}

type Impts struct {
	Imports []Node
}

func (i *Impts) Eval() (s string) {
	s += fmt.Sprintln("start imports")
	for _, im := range i.Imports {
		s += im.Eval()
	}
	s += fmt.Sprintln("end imports")
	return
}

type Impt struct {
	PkgName  string
	ImptName string
}

func (i *Impt) Eval() string {
	return fmt.Sprintln("import: pkgName: " + i.PkgName + " imptName: " + i.ImptName)
}

type Erro struct {
	Desc string
}

func (e *Erro) Eval() string {
	return fmt.Sprintln("error: ", e.Desc)
}

type Consts struct {
	Cs []Node // consts
}

func (c *Consts) Eval() (s string) {
	s += "start const decl\n"
	for _, con := range c.Cs {
		s += con.Eval()
	}
	s += "end const decl\n"
	return
}

// const
type Cnst struct {
	Is Node // idents
	T  Node
	Es Node // expressions
}

func (c *Cnst) Eval() (s string) {
	s += "start const spec\n"
	// subtle cisgendering
	s += c.Is.Eval()
	if c.T != nil {
		s += c.T.Eval()
	}
	if c.Es != nil {
		s += c.Es.Eval()
	}
	s += "end const spec\n"
	return
}

type Idents struct {
	Is []Node
}

func (i *Idents) Eval() (s string) {
	for _, ident := range i.Is {
		s += "ident: " + ident.Eval() + "\n"
	}
	return
}

type Lit struct {
	Typ string
	Val string
}

func (l *Lit) Eval() string {
	return "lit: type: " + l.Typ + " val: " + l.Val + "\n"
}

type OpName struct {
	Id string
}

func (o *OpName) Eval() string {
	return "opname: " + o.Id + "\n"
}

type UnaryE struct {
	Op   string // Operand
	Expr Node
}

func (u *UnaryE) Eval() (s string) {
	s += "op: " + u.Op + "\n"
	s += u.Expr.Eval()
	return
}

// expression list
type Exprs struct {
	Es []Node
}

func (e *Exprs) Eval() (s string) {
	for _, ex := range e.Es {
		s += ex.Eval()
	}
	return
}

// expression list
type Expr struct {
	BinOp   string
	FirstN  Node
	SecondN Node
}

func (e *Expr) Eval() (s string) {
	if e.BinOp != "" {
		s += "binary_op: " + e.BinOp + "\n"
	}
	if e.FirstN != nil {
		s += e.FirstN.Eval()
	}
	if e.SecondN != nil {
		s += e.SecondN.Eval()
	}
	return
}

type Typ struct {
	T Node
}

func (t *Typ) Eval() string {
	return "type: " + t.T.Eval() + "\n"
}

type Ident struct {
	Name string
}

func (i *Ident) Eval() string {
	return i.Name
}

type QualifiedIdent struct {
	Pkg   string
	Ident string
}

func (q *QualifiedIdent) Eval() string {
	return "pkg: " + q.Pkg + " ident: " + q.Ident
}

type Types struct {
	Typspecs []Node
}

func (t *Types) Eval() (s string) {
	s += "start typedecl\n"
	for _, ty := range t.Typspecs {
		s += ty.Eval()
	}
	s += "end typedecl\n"
	return
}

type Typespec struct {
	I   Node //ident
	Typ Node //type
}

func (t *Typespec) Eval() (s string) {
	s += "start typespec\n"
	if t.I != nil {
		s += "ident: " + t.I.Eval() + "\n"
	}
	if t.Typ != nil {
		s += t.Typ.Eval()
	}
	s += "end typespec\n"
	return
}

type Vars struct {
	Vs []Node
}

func (v *Vars) Eval() (s string) {
	s += "start vardecl\n"
	for _, va := range v.Vs {
		s += va.Eval()
	}
	s += "end vardecl\n"
	return
}

type Varspec struct {
	Idents Node
	T      Node // type
	Exprs  Node
}

func (v *Varspec) Eval() (s string) {
	s += "start varspec\n"
	if v.Idents != nil {
		s += v.Idents.Eval()
	}
	if v.T != nil {
		s += v.T.Eval()
	}
	if v.Exprs != nil {
		s += v.Exprs.Eval()
	}
	s += "end varspec\n"
	return
}

type Funcdecl struct {
	Name      Node //ident
	FuncOrSig Node
}

func (f *Funcdecl) Eval() (s string) {
	s += "start funcdecl\n"
	if f.Name != nil {
		s += "ident: " + f.Name.Eval() + "\n"
	}
	if f.FuncOrSig != nil {
		s += f.FuncOrSig.Eval()
	}
	s += "end funcdecl\n"
	return
}

type Func struct {
	Sig  Node
	Body Node
}

func (f *Func) Eval() (s string) {
	if f.Sig != nil {
		s += f.Sig.Eval()
	}
	if f.Body != nil {
		s += f.Body.Eval()
	}
	return
}

type Sig struct {
	Params Node
	Result Node
}

func (sig *Sig) Eval() (s string) {
	if sig.Params != nil {
		s += sig.Params.Eval()
	}
	if sig.Result != nil {
		s += sig.Result.Eval()
	}
	return
}

type Stmts struct {
	Stmts []Node
}

func (ss *Stmts) Eval() (s string) {
	for _, st := range ss.Stmts {
		s += st.Eval()
	}
	return
}

type Stmt struct {
	S Node
}

func (s *Stmt) Eval() string {
	if s.S != nil {
		return s.S.Eval()
	}
	return ""
}

type Result struct {
	ParamsOrTyp Node
}

func (r *Result) Eval() (s string) {
	s += "start result\n"
	if r.ParamsOrTyp != nil {
		s += r.ParamsOrTyp.Eval()
	}
	s += "end result\n"
	return s
}

type Params struct {
	Params []Node
}

func (ps *Params) Eval() (s string) {
	s += "start parameters\n"
	for _, p := range ps.Params {
		s += p.Eval()
	}
	s += "end parameters\n"
	return
}

type Param struct {
	Idents    Node
	DotDotDot bool // if true, apply "..." to type
	Typ       Node
}

func (p *Param) Eval() (s string) {
	s += "start parameterdecl\n"
	if p.Idents != nil {
		s += p.Idents.Eval()
	}
	if p.DotDotDot {
		s += "...\n"
	}
	if p.Typ != nil {
		s += p.Typ.Eval()
	}
	s += "end parameterdecl\n"
	return
}

type Block struct {
	Stmts Node
}

func (b *Block) Eval() (s string) {
	s += "start block\n"
	s += b.Stmts.Eval()
	s += "end block\n"
	return
}
