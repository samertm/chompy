package parse

import (
	"errors"
	"fmt"
)

type Node interface {
	String() string
	Valid() bool
}

type grammarFn func(*parser) Node

type Tree struct {
	Kids []Node
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
	Name string
}

func (p *Pkg) Valid() bool {
	return true
}

func (p *Pkg) String() string {
	return fmt.Sprintln("in package ", p.Name)
}

type Impts struct {
	Imports []Node
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
	PkgName  string
	ImptName string
}

func (i *Impt) Valid() bool {
	return true
}

func (i *Impt) String() string {
	return fmt.Sprintln("import: pkgName: " + i.PkgName + " imptName: " + i.ImptName)
}

type Erro struct {
	Desc string
}

func (e *Erro) Valid() bool {
	return false
}

func (e *Erro) String() string {
	return fmt.Sprintln("error: ", e.Desc)
}

type Consts struct {
	Cs []Node // consts
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
	Is Node // idents
	T  Node
	Es Node // expressions
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
	Is []Node
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
	Typ string
	Val string
}

func (l *Lit) Valid() bool {
	return true
}

func (l *Lit) String() string {
	return "lit: type: " + l.Typ + " val: " + l.Val + "\n"
}

type OpName struct {
	Id Node
}

func (o *OpName) Valid() bool {
	return o.Id != nil && o.Id.Valid()
}

func (o *OpName) String() string {
	return "opname: " + o.Id.String() + "\n"
}

// expression list
type Exprs struct {
	Es []Node
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
	BinOp   string
	FirstN  Node
	SecondN Node
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
	Op   string // Operand
	Expr Node
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
	Expr  Node
	Prime Node
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
	T Node
}

func (t *Typ) Valid() bool {
	return t.T != nil && t.T.Valid()
}

func (t *Typ) String() string {
	return "type: " + t.T.String() + "\n"
}

type Ident struct {
	Name string
}

func (i *Ident) Valid() bool {
	return true
}

func (i *Ident) String() string {
	return i.Name
}

type QualifiedIdent struct {
	Pkg   string
	Ident string
}

func (q *QualifiedIdent) Valid() bool {
	return true
}

func (q *QualifiedIdent) String() string {
	return "pkg: " + q.Pkg + " ident: " + q.Ident
}

type Types struct {
	Typspecs []Node
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
	I   Node //ident
	Typ Node //type
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
	Vs []Node
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
	Idents Node
	T      Node // type
	Exprs  Node
}

func (v *Varspec) Valid() bool {
	return v.Idents != nil && v.T != nil && v.Exprs != nil &&
		v.Idents.Valid() && v.T.Valid() && v.Exprs.Valid()
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
	Name      Node //ident
	FuncOrSig Node
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
	Sig  Node
	Body Node
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
	Params Node
	Result Node
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
	Stmts []Node
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
	S Node
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
	ParamsOrTyp Node
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
	Params []Node
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
	Idents    Node
	DotDotDot bool // if true, apply "..." to type
	Typ       Node
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
	Stmts Node
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
	Label Node // identifier
	Stmt  Node
}

func (l *LabeledStmt) Valid() bool {
	return l.Label != nil && l.Stmt != nil && l.Label.Valid() && l.Stmt.Valid()
}

func (l *LabeledStmt) String() string {
	return "label: " + l.Label.String() + " stmt: " + l.Stmt.String() + "\n"
}

type ExprStmt struct {
	Expr Node
}

func (e *ExprStmt) Valid() bool {
	return e.Expr != nil && e.Expr.Valid()
}

func (e *ExprStmt) String() string {
	return e.Expr.String()
}

type SendStmt struct {
	Chan Node
	Expr Node
}

func (s *SendStmt) Valid() bool {
	return s.Chan != nil && s.Expr != nil && s.Chan.Valid() && s.Expr.Valid()
}

func (s *SendStmt) String() string {
	return "chan: " + s.Chan.String() + " expr: " + s.Expr.String() + "\n"
}

type IncDecStmt struct {
	Expr    Node
	Postfix string // either "++" or "--"
}

func (i *IncDecStmt) Valid() bool {
	return i.Expr != nil && i.Expr.Valid()
}

func (i *IncDecStmt) String() string {
	return "expr: " + i.Expr.String() + " " + i.Postfix + "\n"
}

// Assignment = ExpressionList assign_op ExpressionList .
type Assign struct {
	Op        string // add_op, mul_op, or "="
	LeftExpr  Node
	RightExpr Node
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
	SimpleStmt Node
	Expr       Node
	Block      Node
	Else       Node
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
	Clause Node // ForClause or Condition
	Block  Node
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
	InitStmt  Node
	Condition Node
	PostStmt  Node
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
	ExprsOrIdents Node
	Op            string // "=" or ":="
	Expr          Node   // that comes after the op... need a better nayme
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
	Expr Node
}

func (g *GoStmt) Valid() bool {
	return g.Expr != nil && g.Expr.Valid()
}

func (g *GoStmt) String() string {
	return "go: " + g.Expr.String()
}

type ReturnStmt struct {
	Exprs Node
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
	Label Node
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
	Label Node
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
	Label Node
}

func (g *GotoStmt) Valid() bool {
	return g.Label != nil && g.Label.Valid()
}

func (g *GotoStmt) String() string {
	return "goto: " + g.Label.String() + "\n"
}

type Fallthrough struct {
}

func (f *Fallthrough) Valid() bool {
	return true
}

func (f *Fallthrough) String() string {
	return "fallthrough\n"
}

type DeferStmt struct {
	Expr Node
}

func (d *DeferStmt) Valid() bool {
	return d.Expr != nil && d.Expr.Valid()
}

func (d *DeferStmt) String() string {
	return d.Expr.String()
}

type ShortVarDecl struct {
	Idents Node // identifier list
	Exprs  Node // expression list
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

func (e *EmptyStmt) Valid() bool {
	return true
}

func (e *EmptyStmt) String() string {
	return "empty statement\n"
}

type Conversion struct {
	Typ  Node
	Expr Node
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
	Name Node
	Typ  Node
	Args Node
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
	Ident Node
}

func (s *Selector) Valid() bool {
	return s.Ident != nil && s.Ident.Valid()
}

func (s *Selector) String() string {
	return s.Ident.String()
}

type Index struct {
	Expr Node
}

func (i *Index) Valid() bool {
	return i.Expr != nil && i.Expr.Valid()
}

func (i *Index) String() string {
	return "index: " + i.Expr.String()
}

type Slice struct {
	Start Node
	End   Node
	Cap   Node
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

type Call struct {
	Args Node
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
	Exprs     Node
	DotDotDot bool
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
