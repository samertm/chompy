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
		nodes:   make(chan Node),
	}
	t := sourceFile(p)
	return t
}

// should the states return their list?... probably but not rn
// every nonterminal function assumes that it is in the correct starting state,
// except for sourceFile
func sourceFile(p *parser) *Tree {
	defer close(p.nodes)
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

func packageClause(p *parser) Node {
	p.next() // eat "package"
	if err := p.expect(topPackageName); err != nil {
		return err
	}
	return packageName(p)
}

func packageName(p *parser) Node {
	t := p.next()
	// should I sanity-check t?
	return &Pkg{Name: t.Val}
}

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
	if !p.accept(topImportPath) {
		return &Erro{"expected tokString"}
	}
	// process importPath here.
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

func expressionList(p *parser) Node {
	exs := &Exprs{}
	exs.Es = append(exs.Es, expression(p))
	for p.accept(tokComma) {
		p.next() // eat comma
		exs.Es = append(exs.Es, expression(p))
	}
	return exs
}

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

func primaryExpr(p *parser) Node {
	if p.accept(topOperand...) {
		return operand(p)
	}
	return &Erro{"expected operand"}
}

func operand(p *parser) Node {
	if p.accept(topLiteral...) {
		return literal(p)
	}
	if p.accept(topOperandName) {
		return operandName(p)
	}
	return &Erro{"Expected literal or operand name"}
}

func literal(p *parser) Node {
	if p.accept(topBasicLit...) {
		l := p.next() // int_lit or string_lit
		return &Lit{Typ: l.Typ.String(), Val: l.Val}
	}
	return &Erro{"Expected basic literal"}
}

func operandName(p *parser) Node {
	id := p.next() // get identifier
	return &OpName{Id: id.String()}
}

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

func typeName(p *parser) Node {
	i := p.next() // ident
	if p.accept(tokDot) {
		// is qualified ident
		p.next() // eat "."
		nexti := p.next()
		return &QualifiedIdent{Pkg: i.Val, Ident: nexti.Val}
	}
	return &Ident{Name: i.Val}
}

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

func typeSpec(p *parser) Node {
	spec := &Typespec{}
	spec.I = &Ident{Name: p.next().Val} // ident
	if !p.accept(topType...) {
		return &Erro{"Expected type"}
	}
	spec.Typ = typeGrammar(p)
	return spec
}

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
			return &Result{typeGrammar(p)}
		}
		// saw something other than "(" or ")", assume parameters
		p.push(save)
		return &Result{parameters(p)}
	}
	if !p.accept(topType...) {
		return &Erro{"Expected type or parameters"}
	}
	return &Result{typeGrammar(p)}
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

// Statement =
// 	Declaration .
func statement(p *parser) Node {
	s := &Stmt{}
	if p.accept(topDeclaration...) {
		s.S = declaration(p)
		return s
	}
	return &Erro{"Expected declaration"}
}

// StatementList = { Statement ";" } .
func statementList(p *parser) Node {
	ss := &Stmts{}
	for p.accept(topStatement...) {
		ss.Stmts = append(ss.Stmts, statement(p))
		if err := p.expect(tokSemicolon); err != nil {
			return err
		}
		p.next() // eat ";"
	}
	return ss
}

// Block = "{" StatementList "}" .
func block(p *parser) Node {
	p.next() // eat "{"
	b := &Block{}
	if !p.accept(topStatementList...) {
		return &Erro{"Expected statement list"}
	}
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
