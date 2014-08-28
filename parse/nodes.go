package parse

import (
	"fmt"

	"github.com/samertm/chompy/semantic/stable"
)

type Node interface {
	String() string
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

func (t *Tree) String() (s string) {
	for _, k := range t.Kids {
		s += k.String()
	}
	return
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

func (p *Pkg) String() string {
	return fmt.Sprintln("in package ", p.Name)
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

func (i *Impts) String() (s string) {
	s += fmt.Sprintln("start imports")
	for _, im := range i.Imports {
		s += im.String()
	}
	s += fmt.Sprintln("end imports")
	return
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

func (i *Impt) String() string {
	return fmt.Sprintln("import: pkgName: " + i.PkgName + " imptName: " + i.ImptName)
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

func (e *Erro) String() string {
	return fmt.Sprintln("error: ", e.Desc)
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

func (c *Cnst) String() (s string) {
	s += "start const spec\n"
	// subtle cisgendering
	for _, id := range c.Is {
		s += id.String()
	}
	if c.T != nil {
		s += c.T.String()
	}
	for _, ex := range c.Es {
		s += ex.String()
	}
	s += "end const spec\n"
	return
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

func (l *Lit) String() string {
	return "lit: type: " + l.Typ + " val: " + l.Val + "\n"
}

// expression
type Expr struct {
	BinOp   string
	FirstN  *UnaryE
	SecondN *Expr
	up      Node
}

func (e *Expr) Up() Node {
	return e.up
}

func (e *Expr) SetUp(n Node) {
	e.up = n
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

func (u *UnaryE) String() (s string) {
	s += "unary_op: " + u.Op + "\n"
	s += u.Expr.String()
	return
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

func (p *PrimaryE) String() (s string) {
	s += p.Expr.String()
	if p.Prime != nil {
		s += p.Prime.String()
	}
	return s
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

func (t *Typ) String() string {
	return "type: " + t.T.String() + "\n"
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

func (i *Ident) String() string {
	return "pkg: " + i.Pkg + " ident: " + i.Name
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

func (t *Types) String() (s string) {
	s += "start typedecl\n"
	for _, ty := range t.Typspecs {
		s += ty.String()
	}
	s += "end typedecl\n"
	return
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
	Vs []*Varspec
	up Node
}

func (v *Vars) Up() Node {
	return v.up
}

func (v *Vars) SetUp(n Node) {
	v.up = n
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

func (v *Varspec) String() (s string) {
	s += "start varspec\n"
	for _, id := range v.Idents {
		s += id.String()
	}
	if v.T != nil {
		s += v.T.String()
	}
	for _, ex := range v.Exprs {
		s += ex.String()
	}
	s += "end varspec\n"
	return
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

func (f *Funcdecl) String() (s string) {
	s += "start funcdecl\n"
	if f.Name != nil {
		s += "ident: " + f.Name.String() + "\n"
	}
	if f.Func != nil {
		s += f.Func.String()
	}
	s += "end funcdecl\n"
	return
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

func (sig *Sig) String() (s string) {
	for _, p := range sig.Params {
		s += p.String()
	}
	if sig.Result != nil {
		s += sig.Result.String()
	}
	return
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

func (s *Stmt) String() string {
	if s.S != nil {
		return s.S.String()
	}
	return ""
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

func (r *Result) String() (s string) {
	s += "start result\n"
	for _, p := range r.Params {
		s += p.String()
	}
	if r.Typ != nil {
		s += r.Typ.String()
	}
	s += "end result\n"
	return s
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

func (ps *Params) String() (s string) {
	s += "start parameters\n"
	for _, p := range ps.Params {
		s += p.String()
	}
	s += "end parameters\n"
	return
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

func (p *Param) String() (s string) {
	s += "start parameterdecl\n"
	for _, id := range p.Idents {
		s += id.String()
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
	Stmts []Node
	up    Node
}

func (b *Block) Up() Node {
	return b.up
}

func (b *Block) SetUp(n Node) {
	b.up = n
}

func (b *Block) String() (s string) {
	s += "start block\n"
	for _, st := range b.Stmts {
		s += st.String()
	}
	s += "end block\n"
	return
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

func (l *LabeledStmt) String() string {
	return "label: " + l.Label.String() + " stmt: " + l.Stmt.String() + "\n"
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

func (e *ExprStmt) String() string {
	return e.Expr.String()
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

func (s *SendStmt) String() string {
	return "chan: " + s.Chan.String() + " expr: " + s.Expr.String() + "\n"
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

func (i *IncDecStmt) String() string {
	return "expr: " + i.Expr.String() + " " + i.Postfix + "\n"
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

func (a *Assign) String() (s string) {
	s += "assign_op: " + a.Op + "\n"
	s += "left: "
	for _, ex := range a.LeftExpr {
		s += ex.String()
	}
	s += "right: "
	for _, ex := range a.RightExpr {
		s += ex.String()
	}
	return
}

type IfStmt struct {
	SimpleStmt Node
	Expr       *Expr
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

func (i *IfStmt) String() (s string) {
	if i.SimpleStmt != nil {
		s += i.SimpleStmt.String()
	}
	s += i.Expr.String()
	s += i.Body.String()
	if i.Else != nil {
		s += i.Else.String()
	}
	return
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

func (f *ForStmt) String() (s string) {
	s += f.Clause.String()
	s += f.Body.String()
	return
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

func (r *RangeClause) String() (s string) {
	for _, ex := range r.Exprs {
		s += ex.String()
	}
	for _, id := range r.Idents {
		s += id.String()
	}
	s += "op :" + r.Op + "\n"
	s += r.Expr.String()
	return
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

func (g *GoStmt) String() string {
	return "go: " + g.Expr.String()
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

func (r *ReturnStmt) String() (s string) {
	s += "start return\n"
	for _, ex := range r.Exprs {
		s += ex.String()
	}
	s += "end return\n"
	return
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

func (b *BreakStmt) String() (s string) {
	s += "break: "
	if b.Label != nil {
		s += b.Label.String()
	}
	s += "\n"
	return
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

func (c *ContinueStmt) String() (s string) {
	s += "continue: "
	if c.Label != nil {
		s += c.Label.String()
	}
	s += "\n"
	return
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

func (g *GotoStmt) String() string {
	return "goto: " + g.Label.String() + "\n"
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

func (f *Fallthrough) String() string {
	return "fallthrough\n"
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

func (d *DeferStmt) String() string {
	return d.Expr.String()
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

func (s *ShortVarDecl) String() (str string) {
	str += "start shortvardecl\n"
	for _, id := range s.Idents {
		str += id.String()
	}
	for _, ex := range s.Exprs {
		str += ex.String()
	}
	str += "end shortvardecl\n"
	return
}

type EmptyStmt struct{}

func (e *EmptyStmt) Up() Node {
	return nil
}

func (e *EmptyStmt) SetUp(n Node) {
}

func (e *EmptyStmt) String() string {
	return "empty statement\n"
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

func (c *Conversion) String() (s string) {
	s += "start conversion\n"
	s += c.Typ.String()
	s += c.Expr.String()
	s += "end conversion\n"
	return
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
	Ident *Ident
	up    Node
}

func (s *Selector) Up() Node {
	return s.up
}

func (s *Selector) SetUp(n Node) {
	s.up = n
}

func (s *Selector) String() string {
	return s.Ident.String()
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

func (i *Index) String() string {
	return "index: " + i.Expr.String()
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
	Typ *Typ
	up  Node
}

func (t *TypeAssertion) Up() Node {
	return t.up
}

func (t *TypeAssertion) SetUp(n Node) {
	t.up = n
}

func (t *TypeAssertion) String() string {
	return "type assert: " + t.Typ.String()
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

func (c *Call) String() (s string) {
	s += "start call\n"
	if c.Args != nil {
		s += c.Args.String()
	}
	s += "end call\n"
	return
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

func (a *Args) String() (s string) {
	for _, ex := range a.Exprs {
		s += ex.String()
	}
	if a.DotDotDot {
		s += "...\n"
	}
	return
}
