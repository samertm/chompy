package parse

// - come up with data structure for top-level terminals

import (
	"github.com/samertm/chompy/lex"

	"errors"
	"fmt"
	"log"
)

var _ = log.Fatal // debugging
var _ = errors.New

type parser struct {
	// THOUGHT: might remove curr
	curr    *lex.Token
	toks    chan lex.Token
	oldToks []*lex.Token
	// IDEA: stack of channels?
	// that's pretty good...
	// always emit to the channel on the top
	// push chan when going down a level
	nodes chan Node
	ast   tree
}

func (p *parser) next() *lex.Token {
	if len(p.oldToks) != 0 {
		curr := p.oldToks[0]
		p.oldToks = p.oldToks[1:]
		fmt.Println("oldToks:", curr)
		return curr
	}
	if t, ok := <-p.toks; ok {
		curr := &t
		if curr.Typ == lex.Error {
			log.Fatal("error lexing: ", curr)
			return nil
		} else if curr.Typ == lex.EOF {
			log.Fatal("hit eof")
			return nil
		}
		fmt.Println("chan:   ", curr)
		return curr
	}
	log.Fatal("token stream closed")
	return nil
}

func (p *parser) push(t *lex.Token) {
	if t == nil {
		log.Fatal("bad push")
	}
	p.oldToks = append(p.oldToks, t)
}

func (p *parser) peek() *lex.Token {
	t := p.next()
	p.push(t)
	return t
}

// accept and expect are both based on peek...
func (p *parser) accept(toks ...lex.Token) bool {
	nextTok := p.peek()
	for _, t := range toks {
		if lex.TokenEquiv(*nextTok, t) {
			return true
		}
	}
	return false
}

func (p *parser) expect(tok lex.Token) *erro {
	if p.accept(tok) {
		return nil
	}
	return &erro{"expected " + tok.String()}
}

type Node interface {
	Eval()
}

type grammarFn func(*parser) Node

type tree struct {
	kids []Node
}

func (t *tree) Eval() {
	for _, k := range t.kids {
		k.Eval()
	}
}

type pkg struct {
	name string
}

func (p *pkg) Eval() {
	fmt.Println("in package ", p.name)
}

type impts struct {
	imports []string
}

func (i *impts) Eval() {
	fmt.Println("imports: ", i.imports)
}

type impt struct {
	pkgName  string
	imptName string
}

func (i *impt) Eval() {
	fmt.Println("import: pkgName: " + i.pkgName + " imptName: " + i.imptName)
}

type erro struct {
	desc string
}

func (e *erro) Eval() {
	fmt.Println("error: ", e.desc)
}

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
	p.next()
	if p.accept(topImportDecl) {
		impts := importDecl(p)
		tr.kids = append(tr.kids, impts)
		if err := p.expect(tokSemicolon); err != nil {
			tr.kids = append(tr.kids, err)
		}
		p.next()
	}
	close(p.nodes)
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
	var n Node
	if p.accept(tokOpenParen) {
		p.next() // eat "("
		// TEMP: only grabs one importSpec
		if p.accept(topImportSpec...) {
			n = importSpec(p)
		}
		if err := p.expect(tokSemicolon); err != nil {
			return err
		}
		p.next() // eat ";"
		if err := p.expect(tokCloseParen); err != nil {
			return err
		}
		p.next() // eat ")"
		return n
	}
	// a single importSpec
	if !p.accept(topImportSpec...) {
		return &erro{"expected importSpec"}
	}
	return importSpec(p)
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
