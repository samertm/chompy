package parse

import (
	"fmt"
	"github.com/samertm/chompy/lex"
)

var _ = fmt.Println // debugging

func Start(toks chan lex.Token) Node {
	p := &parser{
		toks:    toks,
		oldToks: make([]*lex.Token, 0),
	}
	t := sourceFile(p)
	return t
}

// SourceFile = PackageClause ";" { ImportDecl ";" } { TopLevelDecl ";" } .
// Every nonterminal function assumes that it is in the correct
// starting state, except for sourceFile.
func sourceFile(p *parser) *Tree {
	tr := &Tree{Kids: make([]Node, 0)}
	if !p.accept(topPackageClause) {
		tr.Kids = append(tr.Kids, &Erro{"PackageClause not found"})
		return tr
	}
	pkg := packageClause(p)
	tr.Kids = append(tr.Kids, pkg)
	if err := p.expect(tokSemicolon); err != nil {
		tr.Kids = append(tr.Kids, err)
		return tr
	}
	p.next() // eat semicolon
	for p.accept(topImportDecl) {
		impts := importDecl(p)
		tr.Kids = append(tr.Kids, impts)
		if err := p.expect(tokSemicolon); err != nil {
			tr.Kids = append(tr.Kids, err)
		}
		p.next() // eat semicolon
	}
	for p.accept(topTopLevelDecl...) {
		topDecl := topLevelDecl(p)
		tr.Kids = append(tr.Kids, topDecl)
		if err := p.expect(tokSemicolon); err != nil {
			tr.Kids = append(tr.Kids, err)
		}
		p.next() // eat semicolon
	}
	return tr
}

// PackageClause  = "package" PackageName .
func packageClause(p *parser) Node {
	p.next() // eat "package"
	if err := p.expect(topPackageName); err != nil {
		return err
	}
	return packageName(p)
}

// PackageName    = identifier .
func packageName(p *parser) Node {
	t := p.next()
	// should I sanity-check t?
	return &Pkg{Name: t.Val}
}

// ImportDecl       = "import" ( ImportSpec | "(" { ImportSpec ";" } ")" ) .
func importDecl(p *parser) Node {
	p.next() // eat "import"
	i := &Impts{Imports: make([]Node, 0)}
	if p.accept(tokOpenParen) {
		p.next() // eat "("
		for p.accept(topImportSpec...) {
			i.Imports = append(i.Imports, importSpec(p))
			if err := p.expect(tokSemicolon); err != nil {
				return err
			}
			p.next() // eat ";"
		}
		if err := p.expect(tokCloseParen); err != nil {
			return err
		}
		p.next() // eat ")"
		return i
	}
	// a single importSpec
	if !p.accept(topImportSpec...) {
		return &Erro{"expected importSpec"}
	}
	i.Imports = append(i.Imports, importSpec(p))
	return i
}

// ImportSpec       = [ "." | PackageName ] ImportPath .
func importSpec(p *parser) Node {
	i := &Impt{}
	if p.accept(tokDot) {
		p.next() // eat dot
		i.PkgName = "."
	}
	if p.accept(topPackageName) {
		t := p.next() // t is the package name
		if i.PkgName == "." {
			// a dot was already processed
			return &Erro{"expected tokString"}
		}
		i.PkgName = t.Val
	}
	// ImportPath       = string_lit .
	if !p.accept(topImportPath) {
		return &Erro{"expected tokString"}
	}
	t := p.next()
	i.ImptName = t.Val
	return i
}

// TopLevelDecl  = Declaration | FunctionDecl .
func topLevelDecl(p *parser) Node {
	if p.accept(topDeclaration...) {
		decl := declaration(p)
		return decl
	}
	if p.accept(topFunctionDecl) {
		fun := functionDecl(p)
		return fun
	}
	return &Erro{"Expected declaration or function declaration"}
}

// Declaration   = ConstDecl | TypeDecl | VarDecl .
func declaration(p *parser) Node {
	if p.accept(topConstDecl) {
		consts := constDecl(p)
		return consts
	}
	if p.accept(topTypeDecl) {
		types := typeDecl(p)
		return types
	}
	if p.accept(topVarDecl) {
		vars := varDecl(p)
		return vars
	}
	return &Erro{"expected const"}
}

// ConstDecl = "const" ( ConstSpec | "(" { ConstSpec ";" } ")" ) .
func constDecl(p *parser) Node {
	p.next() // eat "const"
	cs := &Consts{}
	if p.accept(topConstSpec) {
		cs.Cs = append(cs.Cs, constSpec(p))
		return cs
	}
	if p.accept(tokOpenParen) {
		p.next() // eat "("
		for p.accept(topConstSpec) {
			cs.Cs = append(cs.Cs, constSpec(p))
			if err := p.expect(tokSemicolon); err != nil {
				return err
			}
			p.next() // eat ";"
		}
		if err := p.expect(tokCloseParen); err != nil {
			return err
		}
		p.next() // eat ")"
		return cs
	}
	return &Erro{"expected ConstSpec"}
}

// ConstSpec = IdentifierList [ [ Type ] "=" ExpressionList ] .
func constSpec(p *parser) Node {
	c := &Cnst{}
	c.Is = identifierList(p)
	// type is allowed only if the statement has an expression list
	typeAccepted := false
	if p.accept(topType...) {
		typeAccepted = true
		c.T = typeGrammar(p)
	}
	exprAccepted := false
	if p.accept(tokEqual) {
		exprAccepted = true
		p.next() // eat "="
		c.Es = expressionList(p)
	}
	if typeAccepted == true && exprAccepted == false {
		return &Erro{"Type allowed only if followed by expression"}
	}
	return c
}

// IdentifierList = identifier { "," identifier } .
func identifierList(p *parser) Node {
	idnts := &Idents{}
	id := p.next() // first identifier
	idnts.Is = append(idnts.Is, &Ident{Name: id.Val})
	// look for form: "," identifier
	for p.accept(tokComma) {
		p.next() // throw away ","
		if !p.accept(tokIdentifier) {
			return &Erro{"expected identifier"}
		}
		id = p.next() // identifier
		idnts.Is = append(idnts.Is, &Ident{Name: id.Val})
	}
	return idnts
}

// ExpressionList = Expression { "," Expression } .
func expressionList(p *parser) Node {
	exs := &Exprs{}
	exs.Es = append(exs.Es, expression(p))
	for p.accept(tokComma) {
		p.next() // eat comma
		exs.Es = append(exs.Es, expression(p))
	}
	return exs
}

// Expression = UnaryExpr | Expression binary_op UnaryExpr .
// (equvialent to: Expression = UnaryExpr {binary_op Expression})
func expression(p *parser) Node {
	e := &Expr{}
	firstE := e
	if !p.accept(topUnaryExpr...) {
		return &Erro{"Expected unary expression"}
	}
	e.FirstN = unaryExpr(p)
	for p.accept(tokBinaryOp...) {
		bOp := p.next() // grab binary operator
		e.BinOp = bOp.Val
		if !p.accept(topUnaryExpr...) {
			return &Erro{"Expected unary expression recursed"}
		}
		nextE := &Expr{FirstN: unaryExpr(p)}
		e.SecondN = nextE
		e = nextE
	}
	return firstE
}

// UnaryExpr  = PrimaryExpr | unary_op UnaryExpr .
func unaryExpr(p *parser) Node {
	un := &UnaryE{}
	if p.accept(topPrimaryExpr...) {
		un.Expr = primaryExpr(p)
		return un
	}
	if p.accept(tokUnaryOp...) {
		uOp := p.next() // grab unary operator
		un.Op = uOp.Val
		un.Expr = unaryExpr(p)
		return un
	}
	return &Erro{"expected primary exp or unary_op"}
}

// Operand    = Literal | OperandName .
func operand(p *parser) Node {
	if p.accept(topLiteral...) {
		return literal(p)
	}
	if p.accept(topOperandName) {
		return operandName(p)
	}
	return &Erro{"Expected literal or operand name"}
}

// Literal    = BasicLit .
func literal(p *parser) Node {
	// BasicLit   = int_lit | string_lit .
	if p.accept(topBasicLit...) {
		l := p.next() // int_lit or string_lit
		return &Lit{Typ: l.Typ.String(), Val: l.Val}
	}
	return &Erro{"Expected basic literal"}
}

// OperandName = identifier .
func operandName(p *parser) Node {
	id := p.next() // get identifier
	if p.accept(tokDot) {
		// operand name did not include "."
		// so that type assertions parse correctly
		p.hookTracker()
		p.next()           // eat "."
		if p.accept(tokIdentifier) {
			nextid := p.next() // get identifier
			return &OpName{Id: &QualifiedIdent{Pkg: id.Val, Ident: nextid.Val}}
		}
		p.backtrack()
		p.unhookTracker()
	}
	// fmt.Println("OPERAND NAME", id.Val)
	return &OpName{Id: &Ident{Name: id.Val}}
}

// Type      = TypeName | "(" Type ")" .
func typeGrammar(p *parser) Node {
	if p.accept(topTypeName) {
		t := &Typ{}
		t.T = typeName(p)
		return t
	}
	if p.accept(tokOpenParen) {
		p.next() // eat "("
		t := typeGrammar(p)
		p.next() // eat ")"
		return t
	}
	return &Erro{"Expected type"}
}

// TypeName  = identifier | QualifiedIdent .
func typeName(p *parser) Node {
	id := p.next() // ident
	if p.accept(tokDot) {
		// is qualified ident
		p.next()           // eat "."
		nextid := p.next() // get identifier
		return &QualifiedIdent{Pkg: id.Val, Ident: nextid.Val}
	}
	return &Ident{Name: id.Val}
}

// TypeDecl = "type" ( TypeSpec | "(" { TypeSpec ";" } ")" ) .
func typeDecl(p *parser) Node {
	p.next() // eat "type"
	types := &Types{}
	if p.accept(topTypeSpec) {
		types.Typspecs = append(types.Typspecs, typeSpec(p))
		return types
	}
	if err := p.expect(tokOpenParen); err != nil {
		return err
	}
	p.next() // eat "("
	for p.accept(topTypeSpec) {
		types.Typspecs = append(types.Typspecs, typeSpec(p))
		if err := p.expect(tokSemicolon); err != nil {
			return err
		}
		p.next() // eat ";"
	}
	if err := p.expect(tokCloseParen); err != nil {
		return err
	}
	p.next() // eat ")"
	return types
}

// TypeSpec     = identifier Type .
func typeSpec(p *parser) Node {
	spec := &Typespec{}
	spec.I = &Ident{Name: p.next().Val} // ident
	if !p.accept(topType...) {
		return &Erro{"Expected type"}
	}
	spec.Typ = typeGrammar(p)
	return spec
}

// VarDecl     = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
func varDecl(p *parser) Node {
	p.next() // eat "var"
	vs := &Vars{}
	if p.accept(topVarSpec) {
		vs.Vs = append(vs.Vs, varSpec(p))
		return vs
	}
	if err := p.expect(tokOpenParen); err != nil {
		return err
	}
	p.next() // eat "("
	for p.accept(topVarSpec) {
		vs.Vs = append(vs.Vs, varSpec(p))
		if err := p.expect(tokSemicolon); err != nil {
			return err
		}
		p.next() // eat ";"
	}
	if err := p.expect(tokCloseParen); err != nil {
		return err
	}
	p.next() // eat ")"
	return vs
}

// VarSpec = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
func varSpec(p *parser) Node {
	spec := &Varspec{}
	spec.Idents = identifierList(p)
	if p.accept(topType...) {
		spec.T = typeGrammar(p)
		if p.accept(tokEqual) {
			p.next() // eat "="
			if !p.accept(topExpressionList...) {
				return &Erro{"Expected expression list"}
			}
			spec.Exprs = expressionList(p)
		}
		return spec
	}
	if p.accept(tokEqual) {
		p.next() // eat "="
		if !p.accept(topExpressionList...) {
			return &Erro{"Expected expression list"}
		}
		spec.Exprs = expressionList(p)
		return spec
	}
	return &Erro{"Expected type or expression list"}
}

// ParameterDecl  = [ IdentifierList ] [ "..." ] Type .
func parameterDecl(p *parser) Node {
	par := &Param{}
	if p.accept(topIdentifierList) {
		par.Idents = identifierList(p)
	}
	if p.accept(tokDotDotDot) {
		par.DotDotDot = true
	}
	if !p.accept(topType...) {
		return &Erro{"Expected type"}
	}
	par.Typ = typeGrammar(p)
	return par
}

// ParameterList  = ParameterDecl { "," [ ParameterDecl ] } .
// slightly modified from grammar.txt, so that it will grab a lone ","
func parameterList(p *parser) Node {
	ps := &Params{}
	ps.Params = append(ps.Params, parameterDecl(p))
	for p.accept(tokComma) {
		p.next() // eat ","
		// makes ParameterDecl optional
		if !p.accept(topParameterDecl) {
			return ps
		}
		ps.Params = append(ps.Params, parameterDecl(p))
	}
	return ps
}

// Parameters     = "(" [ ParameterList [ "," ] ] ")" .
func parameters(p *parser) Node {
	p.next() // eat "("
	var ps Node
	if p.accept(topParameterList) {
		ps = parameterList(p)
	}
	if err := p.expect(tokCloseParen); err != nil {
		return err
	}
	p.next() // eat ")"
	return ps
}

// Result         = Parameters | Type .
// Type can start with a (, so we need to check for the couple corner cases first.
// These are checked in order:
// () is an empty type
// any more than one opening paren is a type
// a single paren starts a parameter (topParameters)
// otherwise, check to see if it satisfies topType
func result(p *parser) Node {
	if p.accept(tokOpenParen) {
		save := p.next() // grab "("
		if p.accept(tokCloseParen) || p.accept(tokOpenParen) {
			// saw "()" or "((", assume type
			p.push(save)
			return &Result{ParamsOrTyp: typeGrammar(p)}
		}
		// saw something other than "(" or ")", assume parameters
		p.push(save)
		return &Result{ParamsOrTyp: parameters(p)}
	}
	if !p.accept(topType...) {
		return &Erro{"Expected type or parameters"}
	}
	return &Result{ParamsOrTyp: typeGrammar(p)}
}

// Signature      = Parameters [ Result ] .
func signature(p *parser) Node {
	s := &Sig{}
	s.Params = parameters(p)
	if p.accept(topResult...) {
		s.Result = result(p)
	}
	return s
}

// StatementList = { Statement ";" } .
func statementList(p *parser) Node {
	ss := &Stmts{}
	for p.accept(topStatement...) {
		// fmt.Println("peeking: ", p.peek())
		// s := statement(p)
		// fmt.Println("next stmt: ", s)
		// ss.Stmts = append(ss.Stmts, s)
		ss.Stmts = append(ss.Stmts, statement(p))
		if err := p.expect(tokSemicolon); err != nil {
			return err
		}
		p.next() // eat ";"
	}
	// fmt.Println(ss)
	return ss
}

// Block = "{" StatementList "}" .
func block(p *parser) Node {
	p.next() // eat "{"
	b := &Block{}
	// I don't think I need this check, because I need to allow empty statements
	// if !p.accept(topStatementList...) {
	// 	return &Erro{"Expected statement list, found " + p.peek().String()}
	// }
	b.Stmts = statementList(p)
	if err := p.expect(tokCloseSquiggly); err != nil {
		return err
	}
	p.next() // eat "}"
	return b
}

// FunctionBody = Block .
func functionBody(p *parser) Node {
	// this error check is probably redundant
	if err := p.expect(tokOpenSquiggly); err != nil {
		return err
	}
	return block(p)
}

// Function     = Signature FunctionBody .
func function(p *parser) Node {
	if err := p.expect(topSignature); err != nil {
		return err
	}
	f := &Func{}
	f.Sig = signature(p)
	if err := p.expect(topFunctionBody); err != nil {
		return err
	}
	f.Body = functionBody(p)
	return f
}

// FunctionName = identifier .
func functionName(p *parser) Node {
	i := p.next() // grab ident
	return &Ident{Name: i.Val}
}

// FunctionDecl = "func" FunctionName Function .
func functionDecl(p *parser) Node {
	p.next() // eat "func"
	f := &Funcdecl{}
	if err := p.expect(topFunctionName); err != nil {
		return err
	}
	f.Name = functionName(p)
	if p.accept(topFunction) {
		// only stores funcs for now...
		f.FuncOrSig = function(p)
		return f
	}
	return &Erro{"Expected function"}
}

// DeferStmt = "defer" Expression .
func deferStmt(p *parser) Node {
	p.next() // eat "defer"
	if !p.accept(topExpression...) {
		return &Erro{"Expected expression"}
	}
	return &DeferStmt{Expr: expression(p)}
}

// FallthroughStmt = "fallthrough" .
func fallthroughStmt(p *parser) Node {
	p.next() // eat "fallthrough"
	return &Fallthrough{}
}

// GotoStmt = "goto" Label .
func gotoStmt(p *parser) Node {
	p.next() // eat "goto"
	return &GotoStmt{Label: label(p)}
}

// ContinueStmt = "continue" [ Label ] .
func continueStmt(p *parser) Node {
	p.next() // eat "continue"
	c := &ContinueStmt{}
	if p.accept(topLabel) {
		c.Label = label(p)
	}
	return c
}

// BreakStmt = "break" [ Label ] .
func breakStmt(p *parser) Node {
	p.next() // eat "break"
	b := &BreakStmt{}
	if p.accept(topLabel) {
		b.Label = label(p)
	}
	return b
}

// ReturnStmt = "return" [ ExpressionList ] .
func returnStmt(p *parser) Node {
	p.next() // eat "return"
	r := &ReturnStmt{}
	if p.accept(topExpressionList...) {
		r.Exprs = expressionList(p)
	}
	return r
}

// GoStmt = "go" Expression .
func goStmt(p *parser) Node {
	p.next() // eat "go"
	g := &GoStmt{}
	if !p.accept(topExpression...) {
		return &Erro{"Expected expression"}
	}
	g.Expr = expression(p)
	return g
}

// RangeClause = ( ExpressionList "=" | IdentifierList ":=" ) "range" Expression .
func rangeClause(p *parser) Node {
	// prepare for backtracking
	p.hookTracker()
	defer p.unhookTracker()

	var exprsOrIdents Node
	isIdentList := false
	// check for identifier list first
	if p.accept(topIdentifierList) {
		exprsOrIdents = identifierList(p)
		if exprsOrIdents.Valid() {
			// identifier didn't crap out
			if p.accept(tokColonEqual) {
				p.next() // eat ":="
				isIdentList = true
			}
		}
	}
	if !isIdentList {
		// not an identifier list, try expression list
		p.backtrack()
		if !p.accept(topExpressionList...) {
			return &Erro{"Expected expression list or identifier list"}
		}
		// in expression list
		exprsOrIdents = expressionList(p)
		if !exprsOrIdents.Valid() {
			return &Erro{"Expected expression list or identifier list"}
		}
		// expression list was valid
		if err := p.expect(tokEqual); err != nil {
			return err
		}
		p.next() // eat "="
	}
	// exprsOrIdents is now an expression list or an ident list
	if err := p.expect(tokRange); err != nil {
		return err
	}
	p.next() // eat "range"
	if !p.accept(topExpression...) {
		return &Erro{"expected expression"}
	}
	expr := expression(p)
	return &RangeClause{
		ExprsOrIdents: exprsOrIdents,
		Op:            ":=",
		Expr:          expr,
	}
}

// PostStmt = SimpleStmt .
func postStmt(p *parser) Node {
	return simpleStmt(p)
}

// InitStmt = SimpleStmt .
func initStmt(p *parser) Node {
	return simpleStmt(p)
}

// ForClause = [ InitStmt ] ";" [ Condition ] ";" [ PostStmt ] .
func forClause(p *parser) Node {
	var init, cond, post Node
	if p.accept(topInitStmt...) {
		init = initStmt(p)
	}
	if err := p.expect(tokSemicolon); err != nil {
		return err
	}
	p.next() // eat ";"
	if p.accept(topCondition...) {
		cond = condition(p)
	}
	if err := p.expect(tokSemicolon); err != nil {
		return err
	}
	p.next() // eat ";"
	if p.accept(topPostStmt...) {
		post = postStmt(p)
	}
	return &ForClause{
		InitStmt:  init,
		Condition: cond,
		PostStmt:  post,
	}
}

// Condition = Expression .
func condition(p *parser) Node {
	return expression(p)
}

// ForStmt = "for" [ Condition | ForClause | RangeClause ] Block .
func forStmt(p *parser) Node {
	p.next() // eat "for"
	// see if it accepts no optional stmts
	if !p.accept(topCondition...) &&
		!p.accept(topForClause...) &&
		!p.accept(topRangeClause...) {
		if !p.accept(topBlock) {
			return &Erro{"Expected block"}
		}
		return &ForStmt{Clause: nil, Block: block(p)}
	}
	var clause Node
	if p.accept(topCondition...) {
		clause = condition(p)
	}
	if !clause.Valid() && p.accept(topForClause...) {
		clause = forClause(p)
	}
	if !clause.Valid() && p.accept(topRangeClause...) {
		clause = rangeClause(p)
	}
	if !clause.Valid() {
		// maybe I should just attach this to ForStmt
		return &Erro{"Invalid clause"}
	}
	if !p.accept(topBlock) {
		return &Erro{"Expected block"}
	}
	return &ForStmt{Clause: clause, Block: block(p)}
}

// IfStmt = "if" [ SimpleStmt ";" ] Expression Block [ "else" ( IfStmt | Block ) ] .
func ifStmt(p *parser) Node {
	p.next() // eat "if"

	p.hookTracker()

	ifstmt := &IfStmt{}
	// next expr may be simplestmt or expression

	// check to see if it's a simple statement
	// if we don't see a semicolon, we'll assume that
	// they meant to use an expression and backtrack
	if p.accept(topSimpleStmt...) {
		s := simpleStmt(p)
		if !p.accept(tokSemicolon) {
			// we ate an expression
			p.backtrack()
		}
		p.next() // eat ";"
		ifstmt.SimpleStmt = s
	}
	// stop backtracking
	p.unhookTracker()

	// Expression
	if !p.accept(topExpression...) {
		return &Erro{"Expected expression"}
	}
	ifstmt.Expr = expression(p)
	// Block
	if !p.accept(topBlock) {
		return &Erro{"Expected block"}
	}
	ifstmt.Block = block(p)
	// else
	if p.accept(tokElse) {
		var els Node
		if p.accept(topIfStmt) {
			els = ifStmt(p)
		} else if p.accept(topBlock) {
			els = block(p)
		} else {
			els = &Erro{"Expected if stmt or block"}
		}
		ifstmt.Else = els
	}
	return ifstmt
}

// Assignment = ExpressionList assign_op ExpressionList .
func assignment(p *parser) Node {
	assign := &Assign{}
	assign.LeftExpr = expressionList(p)
	if !p.accept(tokAssignOp...) {
		return &Erro{"Expected assignment operator"}
	}
	op := p.next() // grab operator
	assign.Op = op.Val
	if !p.accept(topExpressionList...) {
		return &Erro{"Expected expression list"}
	}
	assign.RightExpr = expressionList(p)
	return assign
}

// IncDecStmt = Expression ( "++" | "--" )
func incDecStmt(p *parser) Node {
	e := expression(p)
	if !p.accept(tokIncDec...) {
		return &Erro{"Expected '++' or '--'"}
	}
	op := p.next() // grab operator
	return &IncDecStmt{Expr: e, Postfix: op.Val}
}

// Channel  = Expression .
func channel(p *parser) Node {
	return expression(p)
}

// SendStmt = Channel "<-" Expression .
func sendStmt(p *parser) Node {
	c := channel(p)
	if err := p.expect(tokLeftArrow); err != nil {
		return err
	}
	p.next() // eat "<-"
	if !p.accept(topExpression...) {
		return &Erro{"Expected expression"}
	}
	e := expression(p)
	return &SendStmt{
		Chan: c,
		Expr: e,
	}
}

// ExpressionStmt = Expression .
func expressionStmt(p *parser) Node {
	return expression(p)
}

// Label       = identifier .
func label(p *parser) Node {
	i := p.next() // grab ident
	return &Ident{Name: i.Val}
}

// LabeledStmt = Label ":" Statement .
func labeledStmt(p *parser) Node {
	l := label(p)
	if err := p.expect(tokColon); err != nil {
		return err
	}
	p.next() // eat ":"
	if !p.accept(topStatement...) {
		return &Erro{"Expected statement"}
	}
	s := statement(p)
	return &LabeledStmt{
		Label: l,
		Stmt:  s,
	}
}

// EmptyStmt = .
// TODO ...do I need this function?
func emptyStmt(p *parser) Node {
	return &EmptyStmt{}
}

// ShortVarDecl = IdentifierList ":=" ExpressionList .
func shortVarDecl(p *parser) Node {
	// fmt.Println("SHORTVARDECL")
	ids := identifierList(p)
	if err := p.expect(tokColonEqual); err != nil {
		return err
	}
	p.next() // eat ":="
	if !p.accept(topExpressionList...) {
		return &Erro{"Expected expression"}
	}
	e := expressionList(p)
	return &ShortVarDecl{
		Idents: ids,
		Exprs:  e,
	}
}

// SimpleStmt = EmptyStmt | ExpressionStmt | SendStmt | IncDecStmt | Assignment | ShortVarDecl .
func simpleStmt(p *parser) Node {
	// fmt.Println("in simplestmt, peek: ", p.peek())
	var stmt Node
	// set up backtracking
	p.hookTracker()
	defer p.unhookTracker()
	// this is a super inefficient way of doing this lol
	// look at statements in reverse order
	if p.accept(topShortVarDecl) {
		stmt = shortVarDecl(p)
		if stmt.Valid() {
			// fmt.Println("accepted shortvardecl: ", stmt)
			return stmt
		}
		// fmt.Println("BEFORE BACKTRACK: ", p.peek())
		p.backtrack()
		// fmt.Println("AFTER BACKTRACK: ", p.peek())
		// ... my backtracking might not be working...
	}
	// Assignment, IncDecStmt, SendStmt, ExpressionStmt all start with expressions
	if p.accept(topExpression...) {
		// check in order
		// fmt.Println("BEFORE ASSIGNMENT", p.peek())
		stmt = assignment(p)
		if stmt.Valid() {
			return stmt
		}
		// fmt.Println("AFTER ASSIGNMENT")
		// fmt.Println("BEFORE BACKTRACK: ", p.peek())
		p.backtrack()
		// fmt.Println("AFTER BACKTRACK: ", p.peek())
		stmt = incDecStmt(p)
		if stmt.Valid() {
			return stmt
		}
		p.backtrack()
		stmt = sendStmt(p)
		if stmt.Valid() {
			return stmt
		}
		p.backtrack()
		// fmt.Println("BEFORE EXPRESSIONSTMT")
		stmt = expressionStmt(p)
		// fmt.Println("AFTER EXPRESSIONSTMT: ", stmt)
		if stmt.Valid() {
			return stmt
		}
		p.backtrack()

		// none were valid, return error
		return &Erro{"Expected statement meow"}
	}
	// nothing accepted, return empty statement i.e. nil
	return emptyStmt(p)
}

// Statement =
// 	Declaration | LabeledStmt | SimpleStmt |
// 	GoStmt | ReturnStmt | BreakStmt | ContinueStmt | GotoStmt |
// 	FallthroughStmt | Block | IfStmt |  ForStmt |
// 	DeferStmt .
func statement(p *parser) Node {
	// least general first: !LabeledStmt, !SimpleStmt
	// Declaration, GoStmt, ReturnStmt, BreakStmt, ContinueStmt,
	// GotoStmt, FallthroughStmt, Block, IfStmt, ForStmt, DeferStmt
	if p.accept(topDeclaration...) {
		return declaration(p)
	} else if p.accept(topGoStmt) {
		return goStmt(p)
	} else if p.accept(topReturnStmt) {
		return returnStmt(p)
	} else if p.accept(topBreakStmt) {
		return breakStmt(p)
	} else if p.accept(topContinueStmt) {
		return continueStmt(p)
	} else if p.accept(topGotoStmt) {
		return gotoStmt(p)
	} else if p.accept(topFallthroughStmt) {
		return fallthroughStmt(p)
	} else if p.accept(topBlock) {
		return block(p)
	} else if p.accept(topIfStmt) {
		return ifStmt(p)
	} else if p.accept(topForStmt) {
		return forStmt(p)
	} else if p.accept(topDeferStmt) {
		return deferStmt(p)
	}
	// now, LabeledStmt, then SimpleStmt. Accept SimpleStmt as the default,
	// because it can be an EmptyStmt
	// start backtracking
	p.hookTracker()
	defer p.unhookTracker()
	if p.accept(topLabeledStmt) {
		l := labeledStmt(p)
		if l.Valid() {
			return l
		}
		p.backtrack()
	}
	return simpleStmt(p)
}

// ArgumentList   = ExpressionList [ "..." ] .
func argumentList(p *parser) Node {
	a := &Args{}
	a.Exprs = expressionList(p)
	if p.accept(tokDotDotDot) {
		p.next() // eat "..."
		a.DotDotDot = true
	}
	return a
}

// Call           = "(" [ ArgumentList [ "," ] ] ")" .
// TODO right now, conversions are processed as calls
// which means we don't accept the full conversion grammar
// I'm not sure how to fix this yet... because we would need
// to know if something is a type, and right now we only see
// identifiers.
func call(p *parser) Node {
	p.next() // eat "("
	c := &Call{}
	if p.accept(topArgumentList...) {
		c.Args = argumentList(p)
		if p.accept(tokComma) {
			p.next() // eat ","
		}
	}
	if err := p.expect(tokCloseParen); err != nil {
		return err
	}
	p.next() // eat ")"
	return c
}

// TypeAssertion  = "." "(" Type ")" .
func typeAssertion(p *parser) Node {
	// fmt.Println("Enter typeassertion", p.peek())
	p.next() // eat "."
	if err := p.expect(tokOpenParen); err != nil {
		return err
	}
	p.next() // eat "("
	if !p.accept(topType...) {
		return &Erro{"Expected type"}
	}
	t := &TypeAssertion{}
	t.Typ = typeGrammar(p)
	if err := p.expect(tokCloseParen); err != nil {
		return err
	}
	p.next() // eat ")"
	return t
}

// Slice          = "[" ( [ Expression ] ":" [ Expression ] ) |
//                      ( [ Expression ] ":" Expression ":" Expression )
//                  "]" .
// the logic for determining if a slice has the non-optional expressions
// is in Slice.Valid. So, we do not care about that in this function.
func slice(p *parser) Node {
	p.next() // eat "["
	s := &Slice{}
	if p.accept(topExpression...) {
		// fmt.Println("IN START")
		s.Start = expression(p)
	}
	if err := p.expect(tokColon); err != nil {
		return err
	}
	p.next() // eat ":"
	if p.accept(topExpression...) {
		// fmt.Println("IN END")
		s.End = expression(p)
		// fmt.Println("END: ", s.End.Eval())
	}
	if p.accept(tokColon) {
		// on the second ":"
		p.next() // eat ":"
		// fmt.Println("IN CAP")
		if !p.accept(topExpression...) {
			// fmt.Println("OH NOOOO")
			return &Erro{"Expected expression"}
		}
		// fmt.Println("BEFORE CAP")
		s.Cap = expression(p)
		// fmt.Println("AFTER CAP")
		// fmt.Println("CAP: ", s.Cap.Eval())
	}
	if err := p.expect(tokCloseSquareBrace); err != nil {
		return err
	}
	p.next() // eat "]"
	return s
}

// Index          = "[" Expression "]" .
func index(p *parser) Node {
	p.next() // eat "["
	if !p.accept(topExpression...) {
		return &Erro{"Expected expression"}
	}
	i := &Index{}
	i.Expr = expression(p)
	if err := p.expect(tokCloseSquareBrace); err != nil {
		return err
	}
	p.next() // eat "]"
	return i
}

// Selector       = "." identifier .
func selector(p *parser) Node {
	p.next() // eat "."
	s := &Selector{}
	if !p.accept(tokIdentifier) {
		return &Erro{"Expected identifier"}
	}
	ident := p.next() // get identifier
	s.Ident = &Ident{Name: ident.Val}
	return s
}

// PrimaryExprPrime =
//              Selector      [ PrimaryExprPrime ] |
//              Index         [ PrimaryExprPrime ] |
//              Slice         [ PrimaryExprPrime ] |
//              TypeAssertion [ PrimaryExprPrime ] |
//              Call          [ PrimaryExprPrime ] .
func primaryExprPrime(p *parser) Node {
	e := &PrimaryE{}
	p.hookTracker()
	defer p.unhookTracker()
	// fmt.Println("in primaryexprprime, peek: ", p.peek())
	if p.accept(topSelector) {
		s := selector(p)
		if s.Valid() {
			e.Expr = s
			if p.accept(topPrimaryExprPrime...) {
				e.Prime = primaryExprPrime(p)
			}
			return e
		}
		p.backtrack()
	}
	if p.accept(topIndex) {
		i := index(p)
		if i.Valid() {
			e.Expr = i
			if p.accept(topPrimaryExprPrime...) {
				e.Prime = primaryExprPrime(p)
			}
			return e
		}
		p.backtrack()
	}
	if p.accept(topSlice) {
		// fmt.Println("STARTING SLICE")
		s := slice(p)
		if s.Valid() {
			e.Expr = s
			if p.accept(topPrimaryExprPrime...) {
				e.Prime = primaryExprPrime(p)
			}
			return e
		}
		// fmt.Println("ENDED SLICE")
		p.backtrack()
	}
	if p.accept(topTypeAssertion) {
		ta := typeAssertion(p)
		if ta.Valid() {
			e.Expr = ta
			if p.accept(topPrimaryExprPrime...) {
				e.Prime = primaryExprPrime(p)
			}
			return e
		}
		p.backtrack()
	}
	if p.accept(topCall) {
		c := call(p)
		//fmt.Println("HERE", c)
		if c.Valid() {
			e.Expr = c
			if p.accept(topPrimaryExprPrime...) {
				e.Prime = primaryExprPrime(p)
			}
			return e
		}
		// fmt.Println("call invalid")
		p.backtrack()
	}
	p.backtrack()
	return &Erro{"Expected primary expression"}
}

// PrimaryExpr =
// 	Operand     [ PrimaryExprPrime ] |
// 	Conversion  [ PrimaryExprPrime ] |
// 	BuiltinCall [ PrimaryExprPrime ] .
func primaryExpr(p *parser) Node {
	e := &PrimaryE{}
	p.hookTracker()
	defer p.unhookTracker()
	// fmt.Println("in primaryexpr, peek: ", p.peek())
	if p.accept(topBuiltinCall) {
		b := builtinCall(p)
		if b.Valid() {
			e.Expr = b
			if p.accept(topPrimaryExprPrime...) {
				e.Prime = primaryExprPrime(p)
			}
			return e
		}
		p.backtrack()
	}
	// if p.accept(topConversion...) {
	// 	c := conversion(p)
	// 	if c.Valid() {
	// 		e.Expr = c
	// 		if p.accept(topPrimaryExprPrime...) {
	// 			e.Prime = primaryExprPrime(p)
	// 		}
	// 		return e
	// 	}
	// 	p.backtrack()
	// }
	if p.accept(topOperand...) {
		o := operand(p)
		if o.Valid() {
			e.Expr = o
			if p.accept(topPrimaryExprPrime...) {
				// fmt.Println("PRIME")
				e.Prime = primaryExprPrime(p)
				// fmt.Println("prime: ", e.Prime.Eval())
			}
			return e
		}
		p.backtrack()
	}
	p.backtrack()
	return &Erro{"Expected primary expression"}
}

// BuiltinCall = identifier "(" [ BuiltinArgs [ "," ] ] ")" .
// BuiltinArgs = Type [ "," ArgumentList ] | ArgumentList .
func builtinCall(p *parser) Node {
	b := &Builtin{}
	i := p.next() // get identifier
	b.Name = &Ident{Name: i.Val}
	if err := p.expect(tokOpenParen); err != nil {
		return err
	}
	if p.accept(topBuiltinArgs...) {
		// BuiltinArgs = Type [ "," ArgumentList ] | ArgumentList .
		if p.accept(topType...) {
			b.Typ = typeGrammar(p)
			if p.accept(tokComma) {
				p.next() // eat ","
				if !p.accept(topArgumentList...) {
					return &Erro{"Expected argument list"}
				}
				b.Args = argumentList(p)
			}
		} else {
			// argument list
			if !p.accept(topArgumentList...) {
				return &Erro{"Expected argument list"}
			}
			b.Args = argumentList(p)
		}
		// accept comma
		if p.accept(tokComma) {
			p.next() // eat ","
		}
	}
	if err := p.expect(tokCloseParen); err != nil {
		return err
	}
	p.next() // eat ")"
	return b
}

// Conversion = Type "(" Expression [ "," ] ")" .
func conversion(p *parser) Node {
	c := &Conversion{}
	c.Typ = typeGrammar(p)
	if err := p.expect(tokOpenParen); err != nil {
		return err
	}
	p.next() // eat "("
	if !p.accept(topExpression...) {
		return &Erro{"Expected expression"}
	}
	c.Expr = expression(p)
	if p.accept(tokComma) {
		p.next() // eat ","
	}
	if err := p.expect(tokCloseParen); err != nil {
		return err
	}
	p.next() // eat ")"
	return c
}
