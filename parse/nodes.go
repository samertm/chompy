package parse

import (
	"fmt"

	"github.com/samertm/chompy/semantic/stable"
)

type Node interface {
	String() string
	Valid() bool
	// Added interface to make accessing the parent node more
	// convenient.
	Up() Node
	SetUp(Node)
	// gets the immediate children (no grandchildren) of the Node
	// used for walking the tree
	Children(chan<- Node)
}

type grammarFn func(*parser) Node

type Tree struct {
	Up         Node
	RootStable *stable.Stable
	Kids       []Node
}

func (t *Tree) Up() Node {
	return t.Up
}

func (t *Tree) SetUp(n Node) {
	t.Up = n
}

// NOTE do i need this?
// protects the program from runtime errors if the channel is closed
// func protectChildren() {
// 	recover()
// }

func (t *Tree) Children(c chan<- Node) {
	defer close(c)
	for _, k := range t.Kids {
		if k != nil {
			c <- k
		}
	}
}

func (t *Tree) Valid() bool {
	// I believe the tree is valid if it has no kids
	if len(t.Kids) == 0 {
		return true
	}
	for _, k := range t.Kids {
		if k.Valid() == false {
			return false
		}
	}
	return true
}

func (t *Tree) String() (s string) {
	for _, k := range t.Kids {
		s += k.String()
	}
	return
}

type Pkg struct {
	Up   Node
	Name string
}

func (p *Pkg) Up() Node {
	return p.Up
}

func (p *Pkg) SetUp(n Node) {
	p.Up = n
}

func (p *Pkg) Children(c chan<- Node) {
	defer close(c)
	return
}

func (p *Pkg) Valid() bool {
	return true
}

func (p *Pkg) String() string {
	return fmt.Sprintln("in package ", p.Name)
}

type Impts struct {
	Up      Node
	Imports []Node
}

func (i *Impts) Up() Node {
	return i.Up
}

func (i *Impts) SetUp(n Node) {
	i.Up = n
}

func (i *Impts) Children(c chan<- Node) {
	defer close(c)
	for _, im := range i.Imports {
		if im != nil {
			c <- im
		}
	}
}

func (i *Impts) Valid() bool {
	for _, im := range i.Imports {
		if im.Valid() == false {
			return false
		}
	}
	return true
}

func (i *Impts) String() (s string) {
	s += fmt.Sprintln("start imports")
	for _, im := range i.Imports {
		s += im.String()
	}
	s += fmt.Sprintln("end imports")
	return
}

type Impt struct {
	Up       Node
	PkgName  string
	ImptName string
}

func (i *Impt) Up() Node {
	return i.Up
}

func (i *Impt) SetUp(n Node) {
	i.Up = n
}

func (i *Impt) Children(c chan<- Node) {
	defer close(c)
	return
}

func (i *Impt) Valid() bool {
	return true
}

func (i *Impt) String() string {
	return fmt.Sprintln("import: pkgName: " + i.PkgName + " imptName: " + i.ImptName)
}

type Erro struct {
	Up   Node
	Desc string
}

func (e *Erro) Up() Node {
	return e.Up
}

func (e *Erro) SetUp(n Node) {
	e.Up = n
}

func (e *Erro) Children(c chan<- Node) {
	defer close(c)
	return
}

func (e *Erro) Valid() bool {
	return false
}

func (e *Erro) String() string {
	return fmt.Sprintln("error: ", e.Desc)
}

type Consts struct {
	Up Node
	Cs []Node // consts
}

func (con *Consts) Up() Node {
	return con.Up
}

func (con *Consts) SetUp(n Node) {
	con.Up = n
}

func (con *Consts) Children(c chan<- Node) {
	defer close(c)
	for _, cn := range con.Cs {
		if cn != nil {
			c <- cn
		}
	}
}

func (c *Consts) Valid() bool {
	for _, cn := range c.Cs {
		if cn.Valid() == false {
			return false
		}
	}
	return false
}

func (c *Consts) String() (s string) {
	s += "start const decl\n"
	for _, con := range c.Cs {
		s += con.String()
	}
	s += "end const decl\n"
	return
}

// const
type Cnst struct {
	Up Node
	Is Node // idents
	T  Node
	Es Node // expressions
}

func (con *Cnst) Up() Node {
	return con.Up
}

func (con *Cnst) SetUp(n Node) {
	con.Up = n
}

func (con *Cnst) Children(c chan<- Node) {
	defer close(c)
	if con.Is != nil {
		c <- con.Is
	}
	if con.T != nil {
		c <- con.T
	}
	if con.Es != nil {
		c <- con.Es
	}
}

func (c *Cnst) Valid() bool {
	return c.Is != nil && c.T != nil && c.Es != nil &&
		c.Is.Valid() && c.T.Valid() && c.Es.Valid()
}

func (c *Cnst) String() (s string) {
	s += "start const spec\n"
	// subtle cisgendering
	s += c.Is.String()
	if c.T != nil {
		s += c.T.String()
	}
	if c.Es != nil {
		s += c.Es.String()
	}
	s += "end const spec\n"
	return
}

type Idents struct {
	Up Node
	Is []Node
}

func (i *Idents) Up() Node {
	return i.Up
}

func (i *Idents) SetUp(n Node) {
	i.Up = n
}

func (i *Idents) Children(c chan<- Node) {
	defer close(c)
	for _, id := range i.Is {
		if id != nil {
			c <- id
		}
	}
}

func (i *Idents) Valid() bool {
	for _, id := range i.Is {
		if id.Valid() == false {
			return false
		}
	}
	return true
}

func (i *Idents) String() (s string) {
	for _, ident := range i.Is {
		s += "ident: " + ident.String() + "\n"
	}
	return
}

type Lit struct {
	Up  Node
	Typ string
	Val string
}

func (l *Lit) Up() Node {
	return l.Up
}

func (l *Lit) SetUp(n Node) {
	l.Up = n
}

func (l *Lit) Children(c chan<- Node) {
	defer close(c)
	return
}

func (l *Lit) Valid() bool {
	return true
}

func (l *Lit) String() string {
	return "lit: type: " + l.Typ + " val: " + l.Val + "\n"
}

type OpName struct {
	Up Node
	Id Node
}

func (o *OpName) Up() Node {
	return o.Up
}

func (o *OpName) SetUp(n Node) {
	o.Up = n
}

func (o *OpName) Children(c chan<- Node) {
	defer close(c)
	if o.Id != nil {
		c <- o.Id
	}
}

func (o *OpName) Valid() bool {
	return o.Id != nil && o.Id.Valid()
}

func (o *OpName) String() string {
	return "opname: " + o.Id.String() + "\n"
}

// expression list
type Exprs struct {
	Up Node
	Es []Node
}

func (e *Exprs) Up() Node {
	return e.Up
}

func (e *Exprs) SetUp(n Node) {
	e.Up = n
}

func (e *Exprs) Children(c chan<- Node) {
	defer close(c)
	for _, ex := range e.Es {
		if ex != nil {
			c <- ex
		}
	}
}

func (e *Exprs) Valid() bool {
	for _, ex := range e.Es {
		if ex.Valid() == false {
			return false
		}
	}
	return true
}

func (e *Exprs) String() (s string) {
	for _, ex := range e.Es {
		s += ex.String()
	}
	return
}

// expression list
type Expr struct {
	Up      Node
	BinOp   string
	FirstN  Node
	SecondN Node
}

func (e *Expr) Up() Node {
	return e.Up
}

func (e *Expr) SetUp(n Node) {
	e.Up = n
}

func (e *Expr) Children(c chan<- Node) {
	defer close(c)
	if e.FirstN != nil {
		c <- e.FirstN
	}
	if e.SecondN != nil {
		c <- e.SecondN
	}
}

// SecondN can be nil
func (e *Expr) Valid() bool {
	t := e.FirstN != nil && e.FirstN.Valid()
	if e.SecondN != nil {
		t = t && e.SecondN.Valid()
	}
	return t
}

func (e *Expr) String() (s string) {
	if e.BinOp != "" {
		s += "binary_op: " + e.BinOp + "\n"
	}
	if e.FirstN != nil {
		s += e.FirstN.String()
	}
	if e.SecondN != nil {
		s += e.SecondN.String()
	}
	return
}

type UnaryE struct {
	Up   Node
	Op   string // Operand
	Expr Node
}

func (u *UnaryE) Up() Node {
	return u.Up
}

func (u *UnaryE) SetUp(n Node) {
	u.Up = n
}

func (u *UnaryE) Children(c chan<- Node) {
	defer close(c)
	if u.Expr != nil {
		c <- u.Expr
	}
}

func (u *UnaryE) Valid() bool {
	return u.Expr != nil && u.Expr.Valid()
}

func (u *UnaryE) String() (s string) {
	s += "unary_op: " + u.Op + "\n"
	s += u.Expr.String()
	return
}

// PrimaryExprPrimes are also PrimaryExprs
type PrimaryE struct {
	Up    Node
	Expr  Node
	Prime Node
}

func (p *PrimaryE) Up() Node {
	return p.Up
}

func (p *PrimaryE) SetUp(n Node) {
	p.Up = n
}

func (p *PrimaryE) Children(c chan<- Node) {
	defer close(c)
	if p.Expr != nil {
		c <- p.Expr
	}
	if p.Prime != nil {
		c <- p.Prime
	}
}

func (p *PrimaryE) Valid() bool {
	t := p.Expr != nil && p.Expr.Valid()
	if p.Prime != nil {
		t = t && p.Prime.Valid()
	}
	return t
}

func (p *PrimaryE) String() (s string) {
	s += p.Expr.String()
	if p.Prime != nil {
		s += p.Prime.String()
	}
	return s
}

type Typ struct {
	Up Node
	T  Node
}

func (t *Typ) Up() Node {
	return t.Up
}

func (t *Typ) SetUp(n Node) {
	t.Up = n
}

func (t *Typ) Children(c chan<- Node) {
	defer close(c)
	if t.T != nil {
		c <- t.T
	}
}

func (t *Typ) Valid() bool {
	return t.T != nil && t.T.Valid()
}

func (t *Typ) String() string {
	return "type: " + t.T.String() + "\n"
}

type Ident struct {
	Up   Node
	Name string
}

func (i *Ident) Up() Node {
	return i.Up
}

func (i *Ident) SetUp(n Node) {
	i.Up = n
}

func (i *Ident) Children(c chan<- Node) {
	defer close(c)
	return
}

func (i *Ident) Valid() bool {
	return true
}

func (i *Ident) String() string {
	return i.Name
}

type QualifiedIdent struct {
	Up    Node
	Pkg   string
	Ident string
}

func (q *QualifiedIdent) Up() Node {
	return q.Up
}

func (q *QualifiedIdent) SetUp(n Node) {
	q.Up = n
}

func (q *QualifiedIdent) Children(c chan<- Node) {
	defer close(c)
	return
}

func (q *QualifiedIdent) Valid() bool {
	return true
}

func (q *QualifiedIdent) String() string {
	return "pkg: " + q.Pkg + " ident: " + q.Ident
}

type Types struct {
	Up       Node
	Typspecs []Node
}

func (t *Types) Up() Node {
	return t.Up
}

func (t *Types) SetUp(n Node) {
	t.Up = n
}

func (t *Types) Children(c chan<- Node) {
	defer close(c)
	for _, ty := range t.Typspecs {
		if ty != nil {
			c <- ty
		}
	}
}

func (t *Types) Valid() bool {
	for _, ty := range t.Typspecs {
		if ty.Valid() == false {
			return false
		}
	}
	return true
}

func (t *Types) String() (s string) {
	s += "start typedecl\n"
	for _, ty := range t.Typspecs {
		s += ty.String()
	}
	s += "end typedecl\n"
	return
}

type Typespec struct {
	Up  Node
	I   Node //ident
	Typ Node //type
}

func (t *Typespec) Up() Node {
	return t.Up
}

func (t *Typespec) SetUp(n Node) {
	t.Up = n
}

func (t *Typespec) Children(c chan<- Node) {
	defer close(c)
	if t.I != nil {
		c <- t.I
	}
	if t.Typ != nil {
		c <- t.Typ
	}
}

func (t *Typespec) Valid() bool {
	return t.I != nil && t.Typ != nil && t.I.Valid() && t.Typ.Valid()
}

func (t *Typespec) String() (s string) {
	s += "start typespec\n"
	if t.I != nil {
		s += "ident: " + t.I.String() + "\n"
	}
	if t.Typ != nil {
		s += t.Typ.String()
	}
	s += "end typespec\n"
	return
}

type Vars struct {
	Up Node
	Vs []Node
}

func (v *Vars) Up() Node {
	return v.Up
}

func (v *Vars) SetUp(n Node) {
	v.Up = n
}

func (v *Vars) Children(c chan<- Node) {
	defer close(c)
	for _, va := range v.Vs {
		if va != nil {
			c <- va
		}
	}
}

func (v *Vars) Valid() bool {
	for _, va := range v.Vs {
		if va.Valid() == false {
			return false
		}
	}
	return true
}

func (v *Vars) String() (s string) {
	s += "start vardecl\n"
	for _, va := range v.Vs {
		s += va.String()
	}
	s += "end vardecl\n"
	return
}

type Varspec struct {
	Up     Node
	Idents Node
	T      Node // type
	Exprs  Node
}

func (v *Varspec) Up() Node {
	return v.Up
}

func (v *Varspec) SetUp(n Node) {
	v.Up = n
}

func (v *Varspec) Children(c chan<- Node) {
	defer close(c)
	if v.Idents != nil {
		c <- v.Idents
	}
	if v.T != nil {
		c <- v.T
	}
	if v.Exprs != nil {
		c <- v.Exprs
	}
}

func (v *Varspec) Valid() bool {
	t := true
	t = t && v.Idents != nil && v.Idents.Valid() &&
		v.Exprs != nil && v.Exprs.Valid()
	if v.T != nil {
		t = t && v.T.Valid()
	}
	return t
}

func (v *Varspec) String() (s string) {
	s += "start varspec\n"
	if v.Idents != nil {
		s += v.Idents.String()
	}
	if v.T != nil {
		s += v.T.String()
	}
	if v.Exprs != nil {
		s += v.Exprs.String()
	}
	s += "end varspec\n"
	return
}

type Funcdecl struct {
	Up        Node
	Name      Node //ident
	FuncOrSig Node
}

func (f *Funcdecl) Up() Node {
	return f.Up
}

func (f *Funcdecl) SetUp(n Node) {
	f.Up = n
}

func (f *Funcdecl) Children(c chan<- Node) {
	defer close(c)
	if f.Name != nil {
		c <- f.Name
	}
	if f.FuncOrSig != nil {
		c <- f.FuncOrSig
	}
}

func (f *Funcdecl) Valid() bool {
	return f.Name != nil && f.FuncOrSig != nil &&
		f.Name.Valid() && f.FuncOrSig.Valid()
}

func (f *Funcdecl) String() (s string) {
	s += "start funcdecl\n"
	if f.Name != nil {
		s += "ident: " + f.Name.String() + "\n"
	}
	if f.FuncOrSig != nil {
		s += f.FuncOrSig.String()
	}
	s += "end funcdecl\n"
	return
}

type Func struct {
	Up   Node
	Sig  Node
	Body Node
}

func (f *Func) Up() Node {
	return f.Up
}

func (f *Func) SetUp(n Node) {
	f.Up = n
}

func (f *Func) Children(c chan<- Node) {
	defer close(c)
	if f.Sig != nil {
		c <- f.Sig
	}
	if f.Body != nil {
		c <- f.Body
	}
}

func (f *Func) Valid() bool {
	return f.Sig != nil && f.Body != nil &&
		f.Sig.Valid() && f.Body.Valid()
}

func (f *Func) String() (s string) {
	if f.Sig != nil {
		s += f.Sig.String()
	}
	if f.Body != nil {
		s += f.Body.String()
	}
	return
}

type Sig struct {
	Up     Node
	Params Node
	Result Node
}

func (s *Sig) Up() Node {
	return s.Up
}

func (s *Sig) SetUp(n Node) {
	s.Up = n
}

func (s *Sig) Children(c chan<- Node) {
	defer close(c)
	if s.Params != nil {
		c <- s.Params
	}
	if s.Result != nil {
		c <- s.Result
	}
}

func (sig *Sig) Valid() bool {
	t := true
	if sig.Params != nil {
		t = t && sig.Params.Valid()
	}
	if sig.Result != nil {
		t = t && sig.Result.Valid()
	}
	return t
}

func (sig *Sig) String() (s string) {
	if sig.Params != nil {
		s += sig.Params.String()
	}
	if sig.Result != nil {
		s += sig.Result.String()
	}
	return
}

type Stmts struct {
	Up    Node
	Stmts []Node
}

func (s *Stmts) Up() Node {
	return s.Up
}

func (s *Stmts) SetUp(n Node) {
	s.Up = n
}

func (s *Stmts) Children(c chan<- Node) {
	defer close(c)
	for _, ss := range s.Stmts {
		if ss != nil {
			c <- ss
		}
	}
}

func (ss *Stmts) Valid() bool {
	for _, s := range ss.Stmts {
		if s.Valid() == false {
			return false
		}
	}
	return true
}

func (ss *Stmts) String() (s string) {
	for _, st := range ss.Stmts {
		s += st.String()
	}
	return
}

type Stmt struct {
	Up Node
	S  Node
}

func (s *Stmt) Up() Node {
	return s.Up
}

func (s *Stmt) SetUp(n Node) {
	s.Up = n
}

func (s *Stmt) Children(c chan<- Node) {
	defer close(c)
	if s.S != nil {
		c <- s.S
	}
}

func (s *Stmt) Valid() bool {
	return s.S != nil && s.S.Valid()
}

func (s *Stmt) String() string {
	if s.S != nil {
		return s.S.String()
	}
	return ""
}

type Result struct {
	Up          Node
	ParamsOrTyp Node
}

func (r *Result) Up() Node {
	return r.Up
}

func (r *Result) SetUp(n Node) {
	r.Up = n
}

func (r *Result) Children(c chan<- Node) {
	defer close(c)
	if r.ParamsOrTyp != nil {
		c <- r.ParamsOrTyp
	}
}

func (r *Result) Valid() bool {
	return r.ParamsOrTyp != nil && r.ParamsOrTyp.Valid()
}

func (r *Result) String() (s string) {
	s += "start result\n"
	if r.ParamsOrTyp != nil {
		s += r.ParamsOrTyp.String()
	}
	s += "end result\n"
	return s
}

type Params struct {
	Up     Node
	Params []Node
}

func (p *Params) Up() Node {
	return p.Up
}

func (p *Params) SetUp(n Node) {
	p.Up = n
}

func (p *Params) Children(c chan<- Node) {
	defer close(c)
	for _, pa := range p.Params {
		if pa != nil {
			c <- pa
		}
	}
}

func (ps *Params) Valid() bool {
	for _, p := range ps.Params {
		if p.Valid() == false {
			return false
		}
	}
	return true
}

func (ps *Params) String() (s string) {
	s += "start parameters\n"
	for _, p := range ps.Params {
		s += p.String()
	}
	s += "end parameters\n"
	return
}

type Param struct {
	Up        Node
	Idents    Node
	DotDotDot bool // if true, apply "..." to type
	Typ       Node
}

func (p *Param) Up() Node {
	return p.Up
}

func (p *Param) SetUp(n Node) {
	p.Up = n
}

func (p *Param) Children(c chan<- Node) {
	defer close(c)
	if p.Idents != nil {
		c <- p.Idents
	}
	if p.Typ != nil {
		c <- p.Typ
	}
}

func (p *Param) Valid() bool {
	return p.Idents != nil && p.Typ != nil && p.Idents.Valid() && p.Typ.Valid()
}

func (p *Param) String() (s string) {
	s += "start parameterdecl\n"
	if p.Idents != nil {
		s += p.Idents.String()
	}
	if p.DotDotDot {
		s += "...\n"
	}
	if p.Typ != nil {
		s += p.Typ.String()
	}
	s += "end parameterdecl\n"
	return
}

type Block struct {
	Up    Node
	Stmts Node
}

func (b *Block) Up() Node {
	return b.Up
}

func (b *Block) SetUp(n Node) {
	b.Up = n
}

func (b *Block) Children(c chan<- Node) {
	defer close(c)
	if b.Stmts != nil {
		c <- b.Stmts
	}
}

func (b *Block) Valid() bool {
	return b.Stmts != nil && b.Stmts.Valid()
}

func (b *Block) String() (s string) {
	s += "start block\n"
	s += b.Stmts.String()
	s += "end block\n"
	return
}

type LabeledStmt struct {
	Up    Node
	Label Node // identifier
	Stmt  Node
}

func (l *LabeledStmt) Up() Node {
	return l.Up
}

func (l *LabeledStmt) SetUp(n Node) {
	l.Up = n
}

func (l *LabeledStmt) Children(c chan<- Node) {
	defer close(c)
	if l.Label != nil {
		c <- l.Label
	}
	if l.Stmt != nil {
		c <- l.Stmt
	}
}

func (l *LabeledStmt) Valid() bool {
	return l.Label != nil && l.Stmt != nil && l.Label.Valid() && l.Stmt.Valid()
}

func (l *LabeledStmt) String() string {
	return "label: " + l.Label.String() + " stmt: " + l.Stmt.String() + "\n"
}

type ExprStmt struct {
	Up   Node
	Expr Node
}

func (e *ExprStmt) Up() Node {
	return e.Up
}

func (e *ExprStmt) SetUp(n Node) {
	e.Up = n
}

func (e *ExprStmt) Children(c chan<- Node) {
	defer close(c)
	if e.Expr != nil {
		c <- e.Expr
	}
}

func (e *ExprStmt) Valid() bool {
	return e.Expr != nil && e.Expr.Valid()
}

func (e *ExprStmt) String() string {
	return e.Expr.String()
}

type SendStmt struct {
	Up   Node
	Chan Node
	Expr Node
}

func (s *SendStmt) Up() Node {
	return s.Up
}

func (s *SendStmt) SetUp(n Node) {
	s.Up = n
}

func (s *SendStmt) Children(c chan<- Node) {
	defer close(c)
	if s.Chan != nil {
		c <- s.Chan
	}
	if s.Expr != nil {
		c <- s.Expr
	}
}

func (s *SendStmt) Valid() bool {
	return s.Chan != nil && s.Expr != nil && s.Chan.Valid() && s.Expr.Valid()
}

func (s *SendStmt) String() string {
	return "chan: " + s.Chan.String() + " expr: " + s.Expr.String() + "\n"
}

type IncDecStmt struct {
	Up      Node
	Expr    Node
	Postfix string // either "++" or "--"
}

func (i *IncDecStmt) Up() Node {
	return i.Up
}

func (i *IncDecStmt) SetUp(n Node) {
	i.Up = n
}

func (i *IncDecStmt) Children(c chan<- Node) {
	defer close(c)
	if i.Expr != nil {
		c <- i.Expr
	}
}

func (i *IncDecStmt) Valid() bool {
	return i.Expr != nil && i.Expr.Valid()
}

func (i *IncDecStmt) String() string {
	return "expr: " + i.Expr.String() + " " + i.Postfix + "\n"
}

// Assignment = ExpressionList assign_op ExpressionList .
type Assign struct {
	Up        Node
	Op        string // add_op, mul_op, or "="
	LeftExpr  Node
	RightExpr Node
}

func (a *Assign) Up() Node {
	return a.Up
}

func (a *Assign) SetUp(n Node) {
	a.Up = n
}

func (a *Assign) Children(c chan<- Node) {
	defer close(c)
	if a.LeftExpr != nil {
		c <- a.LeftExpr
	}
	if a.RightExpr != nil {
		c <- a.RightExpr
	}
}

func (a *Assign) Valid() bool {
	return a.LeftExpr != nil && a.RightExpr != nil &&
		a.LeftExpr.Valid() && a.RightExpr.Valid()
}

func (a *Assign) String() (s string) {
	s += "assign_op: " + a.Op + "\n"
	s += "left: " + a.LeftExpr.String()
	s += "right: " + a.RightExpr.String()
	return
}

type IfStmt struct {
	Up         Node
	SimpleStmt Node
	Expr       Node
	Block      Node
	Else       Node
}

func (i *IfStmt) Up() Node {
	return i.Up
}

func (i *IfStmt) SetUp(n Node) {
	i.Up = n
}

func (i *IfStmt) Children(c chan<- Node) {
	defer close(c)
	if i.SimpleStmt != nil {
		c <- i.SimpleStmt
	}
	if i.Expr != nil {
		c <- i.Expr
	}
	if i.Block != nil {
		c <- i.Block
	}
	if i.Else != nil {
		c <- i.Else
	}
}

func (i *IfStmt) Valid() bool {
	return i.SimpleStmt != nil && i.Expr != nil && i.Block != nil &&
		i.Else != nil && i.SimpleStmt.Valid() && i.Expr.Valid() &&
		i.Block.Valid() && i.Else.Valid()
}

func (i *IfStmt) String() (s string) {
	if i.SimpleStmt != nil {
		s += i.SimpleStmt.String()
	}
	s += i.Expr.String()
	s += i.Block.String()
	if i.Else != nil {
		s += i.Else.String()
	}
	return
}

type ForStmt struct {
	Up     Node
	Clause Node // ForClause or Condition
	Block  Node
}

func (f *ForStmt) Up() Node {
	return f.Up
}

func (f *ForStmt) SetUp(n Node) {
	f.Up = n
}

func (f *ForStmt) Children(c chan<- Node) {
	defer close(c)
	if f.Clause != nil {
		c <- f.Clause
	}
	if f.Block != nil {
		c <- f.Block
	}
}

func (f *ForStmt) Valid() bool {
	return f.Clause != nil && f.Block != nil &&
		f.Clause.Valid() && f.Block.Valid()
}

func (f *ForStmt) String() (s string) {
	s += f.Clause.String()
	s += f.Block.String()
	return
}

type ForClause struct {
	Up        Node
	InitStmt  Node
	Condition Node
	PostStmt  Node
}

func (f *ForClause) Up() Node {
	return f.Up
}

func (f *ForClause) SetUp(n Node) {
	f.Up = n
}

func (f *ForClause) Children(c chan<- Node) {
	defer close(c)
	if f.InitStmt != nil {
		c <- f.InitStmt
	}
	if f.Condition != nil {
		c <- f.Condition
	}
	if f.PostStmt != nil {
		c <- f.PostStmt
	}
}

func (f *ForClause) Valid() bool {
	return f.InitStmt != nil && f.Condition != nil && f.PostStmt != nil &&
		f.InitStmt.Valid() && f.Condition.Valid() && f.PostStmt.Valid()
}

func (f *ForClause) String() (s string) {
	if f.InitStmt != nil {
		s += f.InitStmt.String()
	}
	if f.Condition != nil {
		s += f.Condition.String()
	}
	if f.PostStmt != nil {
		s += f.PostStmt.String()
	}
	return
}

type RangeClause struct {
	Up            Node
	ExprsOrIdents Node
	Op            string // "=" or ":="
	Expr          Node   // that comes after the op... need a better nayme
}

func (r *RangeClause) Up() Node {
	return r.Up
}

func (r *RangeClause) SetUp(n Node) {
	r.Up = n
}

func (r *RangeClause) Children(c chan<- Node) {
	defer close(c)
	if r.ExprsOrIdents != nil {
		c <- r.ExprsOrIdents
	}
	if r.Expr != nil {
		c <- r.Expr
	}
}

func (r *RangeClause) Valid() bool {
	return r.ExprsOrIdents != nil && r.Expr != nil &&
		r.ExprsOrIdents.Valid() && r.Expr.Valid()
}

func (r *RangeClause) String() (s string) {
	s += r.ExprsOrIdents.String()
	s += "op :" + r.Op + "\n"
	s += r.Expr.String()
	return
}

type GoStmt struct {
	Up   Node
	Expr Node
}

func (g *GoStmt) Up() Node {
	return g.Up
}

func (g *GoStmt) SetUp(n Node) {
	g.Up = n
}

func (g *GoStmt) Children(c chan<- Node) {
	defer close(c)
	if g.Expr != nil {
		c <- g.Expr
	}
}

func (g *GoStmt) Valid() bool {
	return g.Expr != nil && g.Expr.Valid()
}

func (g *GoStmt) String() string {
	return "go: " + g.Expr.String()
}

type ReturnStmt struct {
	Up    Node
	Exprs Node
}

func (r *ReturnStmt) Up() Node {
	return r.Up
}

func (r *ReturnStmt) SetUp(n Node) {
	r.Up = n
}

func (r *ReturnStmt) Children(c chan<- Node) {
	defer close(c)
	if r.Exprs != nil {
		c <- r.Exprs
	}
}

func (r *ReturnStmt) Valid() bool {
	return r.Exprs != nil && r.Exprs.Valid()
}

func (r *ReturnStmt) String() (s string) {
	s += "start return\n"
	if r.Exprs != nil {
		s += r.Exprs.String()
	}
	s += "end return\n"
	return
}

type BreakStmt struct {
	Up    Node
	Label Node
}

func (b *BreakStmt) Up() Node {
	return b.Up
}

func (b *BreakStmt) SetUp(n Node) {
	b.Up = n
}

func (b *BreakStmt) Children(c chan<- Node) {
	defer close(c)
	if b.Label != nil {
		c <- b.Label
	}
}

func (b *BreakStmt) Valid() bool {
	return b.Label != nil && b.Label.Valid()
}

func (b *BreakStmt) String() (s string) {
	s += "break: "
	if b.Label != nil {
		s += b.Label.String()
	}
	s += "\n"
	return
}

type ContinueStmt struct {
	Up    Node
	Label Node
}

func (con *ContinueStmt) Up() Node {
	return con.Up
}

func (con *ContinueStmt) SetUp(n Node) {
	con.Up = n
}

func (con *ContinueStmt) Children(c chan<- Node) {
	defer close(c)
	if con.Label != nil {
		c <- con.Label
	}
}

func (c *ContinueStmt) Valid() bool {
	return c.Label != nil && c.Label.Valid()
}

func (c *ContinueStmt) String() (s string) {
	s += "continue: "
	if c.Label != nil {
		s += c.Label.String()
	}
	s += "\n"
	return
}

type GotoStmt struct {
	Up    Node
	Label Node
}

func (g *GotoStmt) Up() Node {
	return g.Up
}

func (g *GotoStmt) SetUp(n Node) {
	g.Up = n
}

func (g *GotoStmt) Children(c chan<- Node) {
	defer close(c)
	if g.Label != nil {
		c <- g.Label
	}
}

func (g *GotoStmt) Valid() bool {
	return g.Label != nil && g.Label.Valid()
}

func (g *GotoStmt) String() string {
	return "goto: " + g.Label.String() + "\n"
}

type Fallthrough struct {
	Up Node
}

func (f *Fallthrough) Up() Node {
	return f.Up
}

func (f *Fallthrough) SetUp(n Node) {
	f.Up = n
}

func (f *Fallthrough) Children(c chan<- Node) {
	defer close(c)
	return
}

func (f *Fallthrough) Valid() bool {
	return true
}

func (f *Fallthrough) String() string {
	return "fallthrough\n"
}

type DeferStmt struct {
	Up   Node
	Expr Node
}

func (d *DeferStmt) Up() Node {
	return d.Up
}

func (d *DeferStmt) SetUp(n Node) {
	d.Up = n
}

func (d *DeferStmt) Children(c chan<- Node) {
	defer close(c)
	if d.Expr != nil {
		c <- d.Expr
	}
}

func (d *DeferStmt) Valid() bool {
	return d.Expr != nil && d.Expr.Valid()
}

func (d *DeferStmt) String() string {
	return d.Expr.String()
}

type ShortVarDecl struct {
	Up     Node
	Idents Node // identifier list
	Exprs  Node // expression list
}

func (s *ShortVarDecl) Up() Node {
	return s.Up
}

func (s *ShortVarDecl) SetUp(n Node) {
	s.Up = n
}

func (s *ShortVarDecl) Children(c chan<- Node) {
	defer close(c)
	if s.Idents != nil {
		c <- s.Idents
	}
	if s.Exprs != nil {
		c <- s.Exprs
	}
}

func (s *ShortVarDecl) Valid() bool {
	return s.Idents != nil && s.Exprs != nil &&
		s.Idents.Valid() && s.Exprs.Valid()
}

func (s *ShortVarDecl) String() (str string) {
	str += "start shortvardecl\n"
	str += s.Idents.String()
	str += s.Exprs.String()
	str += "end shortvardecl\n"
	return
}

type EmptyStmt struct{}

func (e *EmptyStmt) Up() Node {
	return e.Up
}

func (e *EmptyStmt) SetUp(n Node) {
	e.Up = n
}

func (e *EmptyStmt) Children(c chan<- Node) {
	defer close(c)
	return
}

func (e *EmptyStmt) Valid() bool {
	return true
}

func (e *EmptyStmt) String() string {
	return "empty statement\n"
}

type Conversion struct {
	Up   Node
	Typ  Node
	Expr Node
}

func (con *Conversion) Up() Node {
	return con.Up
}

func (con *Conversion) SetUp(n Node) {
	con.Up = n
}

func (con *Conversion) Children(c chan<- Node) {
	defer close(c)
	if con.Typ != nil {
		c <- con.Typ
	}
	if con.Expr != nil {
		c <- con.Expr
	}
}

func (c *Conversion) Valid() bool {
	return c.Typ != nil && c.Expr != nil && c.Typ.Valid() && c.Expr.Valid()
}

func (c *Conversion) String() (s string) {
	s += "start conversion\n"
	s += c.Typ.String()
	s += c.Expr.String()
	s += "end conversion\n"
	return
}

type Builtin struct {
	Up   Node
	Name Node
	Typ  Node
	Args Node
}

func (b *Builtin) Up() Node {
	return b.Up
}

func (b *Builtin) SetUp(n Node) {
	b.Up = n
}

func (b *Builtin) Children(c chan<- Node) {
	defer close(c)
	if b.Name != nil {
		c <- b.Name
	}
	if b.Typ != nil {
		c <- b.Typ
	}
	if b.Args != nil {
		c <- b.Args
	}
}

func (b *Builtin) Valid() bool {
	t := b.Name != nil && b.Name.Valid() && b.Args != nil && b.Args.Valid()
	if b.Typ != nil {
		t = t && b.Typ.Valid()
	}
	return t
}

func (b *Builtin) String() (s string) {
	s += "start builtin\n"
	s += b.Name.String()
	if b.Typ != nil {
		s += b.Typ.String()
	}
	s += b.Args.String()
	return
}

type Selector struct {
	Up    Node
	Ident Node
}

func (s *Selector) Up() Node {
	return s.Up
}

func (s *Selector) SetUp(n Node) {
	s.Up = n
}

func (s *Selector) Children(c chan<- Node) {
	defer close(c)
	if s.Ident != nil {
		c <- s.Ident
	}
}

func (s *Selector) Valid() bool {
	return s.Ident != nil && s.Ident.Valid()
}

func (s *Selector) String() string {
	return s.Ident.String()
}

type Index struct {
	Up   Node
	Expr Node
}

func (i *Index) Up() Node {
	return i.Up
}

func (i *Index) SetUp(n Node) {
	i.Up = n
}

func (i *Index) Children(c chan<- Node) {
	defer close(c)
	if i.Expr != nil {
		c <- i.Expr
	}
}

func (i *Index) Valid() bool {
	return i.Expr != nil && i.Expr.Valid()
}

func (i *Index) String() string {
	return "index: " + i.Expr.String()
}

type Slice struct {
	Up    Node
	Start Node
	End   Node
	Cap   Node
}

func (s *Slice) Up() Node {
	return s.Up
}

func (s *Slice) SetUp(n Node) {
	s.Up = n
}

func (s *Slice) Children(c chan<- Node) {
	defer close(c)
	if s.Start != nil {
		c <- s.Start
	}
	if s.End != nil {
		c <- s.End
	}
	if s.Cap != nil {
		c <- s.Cap
	}
}

func (s *Slice) Valid() (t bool) {
	if s.Cap != nil {
		// checking:
		// "[" ( [ Expression ] ":" Expression ":" Expression ) "]"
		t = s.End != nil && s.End.Valid() && s.Cap.Valid()
		if s.Start != nil {
			t = t && s.Start.Valid()
		}
	} else {
		// checking:
		// "[" ( [ Expression ] ":" [ Expression ] ) "]"
		t = true
		if s.Start != nil {
			t = t && s.Start.Valid()
		}
		if s.End != nil {
			t = t && s.End.Valid()
		}
	}
	return
}

func (s *Slice) String() (str string) {
	str += "start slice\n"
	if s.Start != nil {
		str += "start: " + s.Start.String()
	}
	if s.End != nil {
		str += "end: " + s.End.String()
	}
	if s.Cap != nil {
		str += "cap: " + s.Cap.String()
	}
	str += "end slice\n"
	return
}

type TypeAssertion struct {
	Up  Node
	Typ Node
}

func (t *TypeAssertion) Up() Node {
	return t.Up
}

func (t *TypeAssertion) SetUp(n Node) {
	t.Up = n
}

func (t *TypeAssertion) Children(c chan<- Node) {
	defer close(c)
	if t.Typ != nil {
		c <- t.Typ
	}
}

func (t *TypeAssertion) Valid() bool {
	return t.Typ != nil && t.Typ.Valid()
}

func (t *TypeAssertion) String() string {
	return "type assert: " + t.Typ.String()
}

type Call struct {
	Up   Node
	Args Node
}

func (con *Call) Up() Node {
	return con.Up
}

func (con *Call) SetUp(n Node) {
	con.Up = n
}

func (con *Call) Children(c chan<- Node) {
	defer close(c)
	if con.Args != nil {
		c <- con.Args
	}
}

func (c *Call) Valid() bool {
	if c.Args != nil {
		return c.Args.Valid()
	}
	return true
}

func (c *Call) String() (s string) {
	s += "start call\n"
	if c.Args != nil {
		s += c.Args.String()
	}
	s += "end call\n"
	return
}

type Args struct {
	Up        Node
	Exprs     Node
	DotDotDot bool
}

func (a *Args) Up() Node {
	return a.Up
}

func (a *Args) SetUp(n Node) {
	a.Up = n
}

func (a *Args) Children(c chan<- Node) {
	defer close(c)
	if a.Exprs != nil {
		c <- a.Exprs
	}
}

func (a *Args) Valid() bool {
	return a.Exprs.Valid()
}

func (a *Args) String() (s string) {
	s += a.Exprs.String()
	if a.DotDotDot {
		s += "...\n"
	}
	return
}
