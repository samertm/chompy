package parse

import (
	"github.com/samertm/chompy/lex"
	"fmt"
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
func sourceFile(p *parser) *tree {
	defer close(p.nodes)
	tr := &tree{kids: make([]Node, 0)}
	if !p.accept(topPackageClause) {
		tr.kids = append(tr.kids, &erro{"PackageClause not found"})
		return tr
	}
	pkg := packageClause(p)
	tr.kids = append(tr.kids, pkg)
	if err := p.expect(tokSemicolon); err != nil {
		tr.kids = append(tr.kids, err)
		return tr
	}
	p.next() // eat semicolon
	for p.accept(topImportDecl) {
		impts := importDecl(p)
		tr.kids = append(tr.kids, impts)
		if err := p.expect(tokSemicolon); err != nil {
			tr.kids = append(tr.kids, err)
		}
		p.next() // eat semicolon
	}
	for p.accept(topTopLevelDecl...) {
		topDecl := topLevelDecl(p)
		tr.kids = append(tr.kids, topDecl)
		if err := p.expect(tokSemicolon); err != nil {
			tr.kids = append(tr.kids, err)
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
	return &pkg{name: t.Val}
}

func importDecl(p *parser) Node {
	p.next() // eat "import"
	i := &impts{imports: make([]Node, 0)}
	if p.accept(tokOpenParen) {
		p.next() // eat "("
		for p.accept(topImportSpec...) {
			i.imports = append(i.imports, importSpec(p))
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
		return &erro{"expected importSpec"}
	}
	i.imports = append(i.imports, importSpec(p))
	return i
}

func importSpec(p *parser) Node {
	i := &impt{}
	if p.accept(tokDot) {
		p.next() // eat dot
		i.pkgName = "."
	}
	if p.accept(topPackageName) {
		t := p.next() // t is the package name
		if i.pkgName == "." {
			// a dot was already processed
			return &erro{"expected tokString"}
		}
		i.pkgName = t.Val
	}
	if !p.accept(topImportPath) {
		return &erro{"expected tokString"}
	}
	// process importPath here.
	t := p.next()
	i.imptName = t.Val
	return i
}

func topLevelDecl(p *parser) Node {
	if p.accept(topDeclaration...) {
		decl := declaration(p)
		return decl
	}
	return &erro{"expected const"}
}

func declaration(p *parser) Node {
	if p.accept(topConstDecl) {
		consts := constDecl(p)
		return consts
	}
	return &erro{"expected const"}
}

func constDecl(p *parser) Node {
	p.next() // eat "const"
	cs := &consts{}
	if p.accept(topConstSpec) {
		cs.cs = append(cs.cs, constSpec(p))
		return cs
	}
	if p.accept(tokOpenParen) {
		p.next() // eat "("
		for p.accept(topConstSpec) {
			cs.cs = append(cs.cs, constSpec(p))
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
	return &erro{"expected ConstSpec"}
}

func constSpec(p *parser) Node {
	c := &cnst{}
	c.is = identifierList(p)
	if p.accept(tokEqual) {
		p.next() // eat "="
		c.es = expressionList(p)
	}
	return c
}

func identifierList(p *parser) Node {
	idnts := &idents{}
	id := p.next() // first identifier
	idnts.is = append(idnts.is, id.Val)
	// look for form: "," identifier
	for p.accept(tokComma) {
		p.next() // throw away ","
		if !p.accept(tokIdentifier) {
			return &erro{"expected identifier"}
		}
		id = p.next() // identifier
		idnts.is = append(idnts.is, id.Val)
	}
	return idnts
}

func expressionList(p *parser) Node {
	exs := &exprs{}
	exs.es = append(exs.es, expression(p))
	for p.accept(tokComma) {
		p.next() // eat comma
		exs.es = append(exs.es, expression(p))
	}
	return exs
}

func expression(p *parser) Node {
	if p.accept(topUnaryExpr...) {
		return unaryExpr(p)
	}
	return &erro{"expected unary expr"}
}

func unaryExpr(p *parser) Node {
	un := &unaryE{}
	if p.accept(topPrimaryExpr...) {
		un.expr = primaryExpr(p)
		return un
	}
	if p.accept(tokUnaryOp...) {
		uOp := p.next() // grab unary operator
		un.op = uOp.Val
		un.expr = unaryExpr(p)
		return un
	}
	return &erro{"expected primary exp or unary_op"}
}

func primaryExpr(p *parser) Node {
	if p.accept(topOperand...) {
		return operand(p)
	}
	return &erro{"expected operand"}
}

func operand(p *parser) Node {
	if p.accept(topLiteral...) {
		return literal(p)
	}
	if p.accept(topOperandName) {
		return operandName(p)
	}
	return &erro{"Expected literal or operand name"}
}

func literal(p *parser) Node {
	if p.accept(topBasicLit...) {
		l := p.next() // int_lit or string_lit
		return &lit{typ: l.Typ.String(), val: l.Val}
	}
	return &erro{"Expected basic literal"}
}

func operandName(p *parser) Node {
	id := p.next() // get identifier
	return &opName{id: id.String()}
}
