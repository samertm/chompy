package parse

import (
	"fmt"
)

type Node interface {
	Eval() string
	Valid() bool
}

type grammarFn func(*parser) Node

type Tree struct {
	Kids []Node
}

func (t *Tree) Valid() bool {
	for _, k := range t.Kids {
		if k.Valid() == false {
			return false
		}
	}
	return true
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

func (p *Pkg) Valid() bool {
	return true
}

func (p *Pkg) Eval() string {
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

func (i *Impt) Valid() bool {
	return true
}

func (i *Impt) Eval() string {
	return fmt.Sprintln("import: pkgName: " + i.PkgName + " imptName: " + i.ImptName)
}

type Erro struct {
	Desc string
}

func (e *Erro) Valid() bool {
	return false
}

func (e *Erro) Eval() string {
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

func (c *Cnst) Valid() bool {
	return c.Is != nil && c.T != nil && c.Es != nil &&
		c.Is.Valid() && c.T.Valid() && c.Es.Valid()
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

func (i *Idents) Valid() bool {
	for _, id := range i.Is {
		if id.Valid() == false {
			return false
		}
	}
	return true
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

func (l *Lit) Valid() bool {
	return true
}

func (l *Lit) Eval() string {
	return "lit: type: " + l.Typ + " val: " + l.Val + "\n"
}

type OpName struct {
	Id string
}

func (o *OpName) Valid() bool {
	return true
}

func (o *OpName) Eval() string {
	return "opname: " + o.Id + "\n"
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

// SecondN can be nil
func (e *Expr) Valid() bool {
	t := e.FirstN != nil && e.FirstN.Valid()
	if e.SecondN != nil {
		t = t && e.SecondN.Valid()
	}
	return t
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

type UnaryE struct {
	Op   string // Operand
	Expr Node
}

func (u *UnaryE) Valid() bool {
	return u.Expr != nil && u.Expr.Valid()
}

func (u *UnaryE) Eval() (s string) {
	s += "unary_op: " + u.Op + "\n"
	s += u.Expr.Eval()
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

func (p *PrimaryE) Eval() (s string) {
	s += p.Expr.Eval()
	if p.Prime != nil {
		s += p.Prime.Eval()
	}
	return s
}

type Typ struct {
	T Node
}

func (t *Typ) Valid() bool {
	return t.T != nil && t.T.Valid()
}

func (t *Typ) Eval() string {
	return "type: " + t.T.Eval() + "\n"
}

type Ident struct {
	Name string
}

func (i *Ident) Valid() bool {
	return true
}

func (i *Ident) Eval() string {
	return i.Name
}

type QualifiedIdent struct {
	Pkg   string
	Ident string
}

func (q *QualifiedIdent) Valid() bool {
	return true
}

func (q *QualifiedIdent) Eval() string {
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

func (t *Typespec) Valid() bool {
	return t.I != nil && t.Typ != nil && t.I.Valid() && t.Typ.Valid()
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

func (v *Vars) Valid() bool {
	for _, va := range v.Vs {
		if va.Valid() == false {
			return false
		}
	}
	return true
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

func (v *Varspec) Valid() bool {
	return v.Idents != nil && v.T != nil && v.Exprs != nil &&
		v.Idents.Valid() && v.T.Valid() && v.Exprs.Valid()
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

func (f *Funcdecl) Valid() bool {
	return f.Name != nil && f.FuncOrSig != nil &&
		f.Name.Valid() && f.FuncOrSig.Valid()
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

func (f *Func) Valid() bool {
	return f.Sig != nil && f.Body != nil &&
		f.Sig.Valid() && f.Body.Valid()
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

func (sig *Sig) Valid() bool {
	return sig.Params != nil && sig.Result != nil &&
		sig.Params.Valid() && sig.Result.Valid()
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

func (ss *Stmts) Valid() bool {
	for _, s := range ss.Stmts {
		if s.Valid() == false {
			return false
		}
	}
	return true
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

func (s *Stmt) Valid() bool {
	return s.S != nil && s.S.Valid()
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

func (r *Result) Valid() bool {
	return r.ParamsOrTyp != nil && r.ParamsOrTyp.Valid()
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

func (ps *Params) Valid() bool {
	for _, p := range ps.Params {
		if p.Valid() == false {
			return false
		}
	}
	return true
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

func (p *Param) Valid() bool {
	return p.Idents != nil && p.Typ != nil && p.Idents.Valid() && p.Typ.Valid()
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

func (b *Block) Valid() bool {
	return b.Stmts != nil && b.Stmts.Valid()
}

func (b *Block) Eval() (s string) {
	s += "start block\n"
	s += b.Stmts.Eval()
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

func (l *LabeledStmt) Eval() string {
	return "label: " + l.Label.Eval() + " stmt: " + l.Stmt.Eval() + "\n"
}

type ExprStmt struct {
	Expr Node
}

func (e *ExprStmt) Valid() bool {
	return e.Expr != nil && e.Expr.Valid()
}

func (e *ExprStmt) Eval() string {
	return e.Expr.Eval()
}

type SendStmt struct {
	Chan Node
	Expr Node
}

func (s *SendStmt) Valid() bool {
	return s.Chan != nil && s.Expr != nil && s.Chan.Valid() && s.Expr.Valid()
}

func (s *SendStmt) Eval() string {
	return "chan: " + s.Chan.Eval() + " expr: " + s.Expr.Eval() + "\n"
}

type IncDecStmt struct {
	Expr    Node
	Postfix string // either "++" or "--"
}

func (i *IncDecStmt) Valid() bool {
	return i.Expr != nil && i.Expr.Valid()
}

func (i *IncDecStmt) Eval() string {
	return "expr: " + i.Expr.Eval() + " " + i.Postfix + "\n"
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

func (a *Assign) Eval() (s string) {
	s += "assign_op: " + a.Op + "\n"
	s += "left: " + a.LeftExpr.Eval()
	s += "right: " + a.RightExpr.Eval()
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

func (i *IfStmt) Eval() (s string) {
	if i.SimpleStmt != nil {
		s += i.SimpleStmt.Eval()
	}
	s += i.Expr.Eval()
	s += i.Block.Eval()
	if i.Else != nil {
		s += i.Else.Eval()
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

func (f *ForStmt) Eval() (s string) {
	s += f.Clause.Eval()
	s += f.Block.Eval()
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

func (f *ForClause) Eval() (s string) {
	if f.InitStmt != nil {
		s += f.InitStmt.Eval()
	}
	if f.Condition != nil {
		s += f.Condition.Eval()
	}
	if f.PostStmt != nil {
		s += f.PostStmt.Eval()
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

func (r *RangeClause) Eval() (s string) {
	s += r.ExprsOrIdents.Eval()
	s += "op :" + r.Op + "\n"
	s += r.Expr.Eval()
	return
}

type GoStmt struct {
	Expr Node
}

func (g *GoStmt) Valid() bool {
	return g.Expr != nil && g.Expr.Valid()
}

func (g *GoStmt) Eval() string {
	return "go: " + g.Expr.Eval()
}

type ReturnStmt struct {
	Exprs Node
}

func (r *ReturnStmt) Valid() bool {
	return r.Exprs != nil && r.Exprs.Valid()
}

func (r *ReturnStmt) Eval() (s string) {
	s += "start return\n"
	if r.Exprs != nil {
		s += r.Exprs.Eval()
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

func (b *BreakStmt) Eval() (s string) {
	s += "break: "
	if b.Label != nil {
		s += b.Label.Eval()
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

func (c *ContinueStmt) Eval() (s string) {
	s += "continue: "
	if c.Label != nil {
		s += c.Label.Eval()
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

func (g *GotoStmt) Eval() string {
	return "goto: " + g.Label.Eval() + "\n"
}

type Fallthrough struct {
}

func (f *Fallthrough) Valid() bool {
	return true
}

func (f *Fallthrough) Eval() string {
	return "fallthrough\n"
}

type DeferStmt struct {
	Expr Node
}

func (d *DeferStmt) Valid() bool {
	return d.Expr != nil && d.Expr.Valid()
}

func (d *DeferStmt) Eval() string {
	return d.Expr.Eval()
}

type ShortVarDecl struct {
	Idents Node // identifier list
	Exprs  Node // expression list
}

func (s *ShortVarDecl) Valid() bool {
	return s.Idents != nil && s.Exprs != nil &&
		s.Idents.Valid() && s.Exprs.Valid()
}

func (s *ShortVarDecl) Eval() (str string) {
	str += "start shortvardecl\n"
	str += s.Idents.Eval()
	str += s.Exprs.Eval()
	str += "end shortvardecl\n"
	return
}

type EmptyStmt struct{}

func (e *EmptyStmt) Valid() bool {
	return true
}

func (e *EmptyStmt) Eval() string {
	return "empty statement\n"
}

type Conversion struct {
	Typ  Node
	Expr Node
}

func (c *Conversion) Valid() bool {
	return c.Typ != nil && c.Expr != nil && c.Typ.Valid() && c.Expr.Valid()
}

func (c *Conversion) Eval() (s string) {
	s += "start conversion\n"
	s += c.Typ.Eval()
	s += c.Expr.Eval()
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

func (b *Builtin) Eval() (s string) {
	s += "start builtin\n"
	s += b.Name.Eval()
	if b.Typ != nil {
		s += b.Typ.Eval()
	}
	s += b.Args.Eval()
	return
}

type Selector struct {
	Ident Node
}

func (s *Selector) Valid() bool {
	return s.Ident != nil && s.Ident.Valid()
}

func (s *Selector) Eval() string {
	return s.Ident.Eval()
}

type Index struct {
	Expr Node
}

func (i *Index) Valid() bool {
	return i.Expr != nil && i.Expr.Valid()
}

func (i *Index) Eval() string {
	return "index: " + i.Expr.Eval()
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

func (s *Slice) Eval() (str string) {
	str += "start slice\n"
	if s.Start != nil {
		str += "start: " + s.Start.Eval()
	}
	if s.End != nil {
		str += "end: " + s.End.Eval()
	}
	if s.Cap != nil {
		str += "cap: " + s.Cap.Eval()
	}
	str += "end slice\n"
	return
}

type TypeAssertion struct {
	Typ Node
}

func (t *TypeAssertion) Valid() bool {
	return t.Typ != nil && t.Typ.Valid()
}

func (t *TypeAssertion) Eval() string {
	return "type assert: " + t.Typ.Eval()
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

func (c *Call) Eval() (s string) {
	s += "start call\n"
	s += c.Args.Eval()
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

func (a *Args) Eval() (s string) {
	s += a.Exprs.Eval()
	if a.DotDotDot {
		s += "...\n"
	}
	return
}
