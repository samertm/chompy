package parse

import "github.com/samertm/chompy/semantic/stable"

type Node interface {
	// String() string
	// Added interface to make accessing the parent node more
	// convenient.
	Up() Node
	SetUp(Node)
	// gets the immediate children (no grandchildren) of the Node
	// used for walking the tree
	// Children(chan<- Node)
}

type grammarFn func(*parser) Node

type Tree struct {
	RootStable *stable.Stable
	Kids       []Node
	up         Node
}

func (t *Tree) Up() Node {
	return t.up
}

func (t *Tree) SetUp(n Node) {
	t.up = n
}

// NOTE do i need this?
// protects the program from runtime errors if the channel is closed
// func protectChildren() {
// 	recover()
// }

type Pkg struct {
	Name string
	up   Node
}

func (p *Pkg) Up() Node {
	return p.up
}

func (p *Pkg) SetUp(n Node) {
	p.up = n
}

type Impts struct {
	Imports []*Impt
	up      Node
}

func (i *Impts) Up() Node {
	return i.up
}

func (i *Impts) SetUp(n Node) {
	i.up = n
}

type Impt struct {
	PkgName  string
	ImptName string
	up       Node
}

func (i *Impt) Up() Node {
	return i.up
}

func (i *Impt) SetUp(n Node) {
	i.up = n
}

type Erro struct {
	Desc string
}

// Up and SetUp are nops for the error type, because they get removed
// in the first semantic pass, and it breaks too much of the grammar.
func (e *Erro) Up() Node {
	return nil
}

// nop
func (e *Erro) SetUp(n Node) {
}

type Consts struct {
	Cs []*Cnst // consts
	up Node
}

func (con *Consts) Up() Node {
	return con.up
}

func (con *Consts) SetUp(n Node) {
	con.up = n
}

// const
type Cnst struct {
	Is []*Ident // idents
	T  *Typ
	Es []*Expr // expressions
	up Node
}

func (con *Cnst) Up() Node {
	return con.up
}

func (con *Cnst) SetUp(n Node) {
	con.up = n
}

type Lit struct {
	Typ string
	Val string
	up  Node
}

func (l *Lit) Up() Node {
	return l.up
}

func (l *Lit) SetUp(n Node) {
	l.up = n
}

// expression
type Expr struct {
	BinOp   string
	FirstN  Node
	SecondN Node
	up      Node
}

func (e *Expr) Up() Node {
	return e.up
}

func (e *Expr) SetUp(n Node) {
	e.up = n
}

type UnaryE struct {
	Op   string // Operand
	Expr Node
	up   Node
}

func (u *UnaryE) Up() Node {
	return u.up
}

func (u *UnaryE) SetUp(n Node) {
	u.up = n
}

// PrimaryExprPrimes are also PrimaryExprs
type PrimaryE struct {
	Expr  Node
	Prime *PrimaryE
	up    Node
}

func (p *PrimaryE) Up() Node {
	return p.up
}

func (p *PrimaryE) SetUp(n Node) {
	p.up = n
}

type Typ struct {
	T  Node
	up Node
}

func (t *Typ) Up() Node {
	return t.up
}

func (t *Typ) SetUp(n Node) {
	t.up = n
}

type Ident struct {
	Name string
	Pkg  string
	up   Node
}

func (i *Ident) Up() Node {
	return i.up
}

func (i *Ident) SetUp(n Node) {
	i.up = n
}

type Types struct {
	Typspecs []*Typespec
	up       Node
}

func (t *Types) Up() Node {
	return t.up
}

func (t *Types) SetUp(n Node) {
	t.up = n
}

type Typespec struct {
	I   *Ident //ident
	Typ *Typ   //type
	up  Node
}

func (t *Typespec) Up() Node {
	return t.up
}

func (t *Typespec) SetUp(n Node) {
	t.up = n
}

type Vars struct {
	Vs []*Varspec
	up Node
}

func (v *Vars) Up() Node {
	return v.up
}

func (v *Vars) SetUp(n Node) {
	v.up = n
}

type Varspec struct {
	Idents []*Ident
	T      *Typ // type
	Exprs  []*Expr
	up     Node
}

func (v *Varspec) Up() Node {
	return v.up
}

func (v *Varspec) SetUp(n Node) {
	v.up = n
}

type Funcdecl struct {
	Name *Ident //ident
	Func *Func
	up   Node
}

func (f *Funcdecl) Up() Node {
	return f.up
}

func (f *Funcdecl) SetUp(n Node) {
	f.up = n
}

type Func struct {
	Sig  *Sig
	Body *Block
	up   Node
}

func (f *Func) Up() Node {
	return f.up
}

func (f *Func) SetUp(n Node) {
	f.up = n
}

type Sig struct {
	Params []*Param
	Result *Result
	up     Node
}

func (s *Sig) Up() Node {
	return s.up
}

func (s *Sig) SetUp(n Node) {
	s.up = n
}

type Stmt struct {
	S  Node
	up Node
}

func (s *Stmt) Up() Node {
	return s.up
}

func (s *Stmt) SetUp(n Node) {
	s.up = n
}

type Result struct {
	// <HACK>pretend this is a union
	Params []*Param
	Typ    *Typ
	// </HACK>
	up Node
}

func (r *Result) Up() Node {
	return r.up
}

func (r *Result) SetUp(n Node) {
	r.up = n
}

type Params struct {
	Params []*Param
	up     Node
}

func (p *Params) Up() Node {
	return p.up
}

func (p *Params) SetUp(n Node) {
	p.up = n
}

type Param struct {
	Idents    []*Ident
	DotDotDot bool // if true, apply "..." to type
	Typ       *Typ
	up        Node
}

func (p *Param) Up() Node {
	return p.up
}

func (p *Param) SetUp(n Node) {
	p.up = n
}

type Block struct {
	Stmts []Node
	up    Node
}

func (b *Block) Up() Node {
	return b.up
}

func (b *Block) SetUp(n Node) {
	b.up = n
}

type LabeledStmt struct {
	Label *Ident // identifier
	Stmt  Node
	up    Node
}

func (l *LabeledStmt) Up() Node {
	return l.up
}

func (l *LabeledStmt) SetUp(n Node) {
	l.up = n
}

type ExprStmt struct {
	Expr Node
	up   Node
}

func (e *ExprStmt) Up() Node {
	return e.up
}

func (e *ExprStmt) SetUp(n Node) {
	e.up = n
}

type SendStmt struct {
	Chan Node
	Expr Node
	up   Node
}

func (s *SendStmt) Up() Node {
	return s.up
}

func (s *SendStmt) SetUp(n Node) {
	s.up = n
}

type IncDecStmt struct {
	Expr    Node
	Postfix string // either "++" or "--"
	up      Node
}

func (i *IncDecStmt) Up() Node {
	return i.up
}

func (i *IncDecStmt) SetUp(n Node) {
	i.up = n
}

// Assignment = ExpressionList assign_op ExpressionList .
type Assign struct {
	Op        string // add_op, mul_op, or "="
	LeftExpr  []*Expr
	RightExpr []*Expr
	up        Node
}

func (a *Assign) Up() Node {
	return a.up
}

func (a *Assign) SetUp(n Node) {
	a.up = n
}

type IfStmt struct {
	SimpleStmt Node
	Expr       Node
	Body       *Block
	Else       Node
	up         Node
}

func (i *IfStmt) Up() Node {
	return i.up
}

func (i *IfStmt) SetUp(n Node) {
	i.up = n
}

type ForStmt struct {
	Clause Node // ForClause or Condition
	Body   *Block
	up     Node
}

func (f *ForStmt) Up() Node {
	return f.up
}

func (f *ForStmt) SetUp(n Node) {
	f.up = n
}

type ForClause struct {
	InitStmt  Node
	Condition Node
	PostStmt  Node
	up        Node
}

func (f *ForClause) Up() Node {
	return f.up
}

func (f *ForClause) SetUp(n Node) {
	f.up = n
}

type RangeClause struct {
	// <HACK>Think of this as a union: it has one or the other.
	Exprs  []*Expr
	Idents []*Ident
	// </HACK>
	Op   string // "=" or ":="
	Expr Node   // that comes after the op... need a better nayme
	up   Node
}

func (r *RangeClause) Up() Node {
	return r.up
}

func (r *RangeClause) SetUp(n Node) {
	r.up = n
}

type GoStmt struct {
	Expr Node
	up   Node
}

func (g *GoStmt) Up() Node {
	return g.up
}

func (g *GoStmt) SetUp(n Node) {
	g.up = n
}

type ReturnStmt struct {
	Exprs []*Expr
	up    Node
}

func (r *ReturnStmt) Up() Node {
	return r.up
}

func (r *ReturnStmt) SetUp(n Node) {
	r.up = n
}

type BreakStmt struct {
	Label *Ident
	up    Node
}

func (b *BreakStmt) Up() Node {
	return b.up
}

func (b *BreakStmt) SetUp(n Node) {
	b.up = n
}

type ContinueStmt struct {
	Label *Ident
	up    Node
}

func (con *ContinueStmt) Up() Node {
	return con.up
}

func (con *ContinueStmt) SetUp(n Node) {
	con.up = n
}

type GotoStmt struct {
	Label *Ident
	up    Node
}

func (g *GotoStmt) Up() Node {
	return g.up
}

func (g *GotoStmt) SetUp(n Node) {
	g.up = n
}

type Fallthrough struct {
	up Node
}

func (f *Fallthrough) Up() Node {
	return f.up
}

func (f *Fallthrough) SetUp(n Node) {
	f.up = n
}

type DeferStmt struct {
	Expr Node
	up   Node
}

func (d *DeferStmt) Up() Node {
	return d.up
}

func (d *DeferStmt) SetUp(n Node) {
	d.up = n
}

type ShortVarDecl struct {
	Idents []*Ident // identifier list
	Exprs  []*Expr  // expression list
	up     Node
}

func (s *ShortVarDecl) Up() Node {
	return s.up
}

func (s *ShortVarDecl) SetUp(n Node) {
	s.up = n
}


type EmptyStmt struct{}

func (e *EmptyStmt) Up() Node {
	return nil
}

func (e *EmptyStmt) SetUp(n Node) {
}

type Conversion struct {
	Typ  *Typ
	Expr Node
	up   Node
}

func (con *Conversion) Up() Node {
	return con.up
}

func (con *Conversion) SetUp(n Node) {
	con.up = n
}

type Builtin struct {
	Name *Ident
	Typ  *Typ
	Args *Args
	up   Node
}

func (b *Builtin) Up() Node {
	return b.up
}

func (b *Builtin) SetUp(n Node) {
	b.up = n
}

type Selector struct {
	Ident *Ident
	up    Node
}

func (s *Selector) Up() Node {
	return s.up
}

func (s *Selector) SetUp(n Node) {
	s.up = n
}

type Index struct {
	Expr Node
	up   Node
}

func (i *Index) Up() Node {
	return i.up
}

func (i *Index) SetUp(n Node) {
	i.up = n
}

type Slice struct {
	Start Node
	End   Node
	Cap   Node
	up    Node
}

func (s *Slice) Up() Node {
	return s.up
}

func (s *Slice) SetUp(n Node) {
	s.up = n
}

type TypeAssertion struct {
	Typ *Typ
	up  Node
}

func (t *TypeAssertion) Up() Node {
	return t.up
}

func (t *TypeAssertion) SetUp(n Node) {
	t.up = n
}

type Call struct {
	Args *Args
	up   Node
}

func (con *Call) Up() Node {
	return con.up
}

func (con *Call) SetUp(n Node) {
	con.up = n
}

type Args struct {
	Exprs     []*Expr
	DotDotDot bool
	up        Node
}

func (a *Args) Up() Node {
	return a.up
}

func (a *Args) SetUp(n Node) {
	a.up = n
}
