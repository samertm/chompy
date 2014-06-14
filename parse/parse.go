package parse

// what to do tomorrow?
// -something, probably

import (
	"github.com/samertm/chompy/lex"

	"errors"
	"fmt"
	"log"
)

var _ = log.Fatal // debugging

type Node interface {
	Eval()
}

type tree struct {
	kids []*Node
}

type parser struct {
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

type pkg struct {
	name string
}

func (p pkg) Eval() {
	fmt.Println("in package ", p.name)
}

type impt struct {
	imports []string
}

func (i impt) Eval() {
	fmt.Println("imports: ", i.imports)
}

type erro struct {
	desc string
}

func (e erro) Eval() {
	fmt.Println("error: ", e.desc)
}

func (p *parser) next() error {
	if len(p.oldToks) != 0 {
		p.curr, p.oldToks = p.oldToks[0], p.oldToks[1:]
		fmt.Println("next: ", p.curr)
		return nil
	}
	if t, ok := <-p.toks; ok {
		p.curr = &t
		fmt.Println("next:", p.curr)
		return nil
	}
	log.Fatal("reached eof")
	return errors.New("unexpected EOF")
}

func (p *parser) backup() {
	if p.curr == nil {
		log.Fatal("bad push")
	}
	p.oldToks = append(p.oldToks, p.curr)
	p.curr = nil
}

func (p *parser) accept(typ lex.TokenType, val string) bool {
	if p.curr.Typ == typ {
		if val == "" {
			return true
		}
		return p.curr.Val == val
	}
	return false
}

func (p *parser) expect(typ lex.TokenType, val string) {
	if !p.accept(typ, val) {
		p.nodes <- erro{fmt.Sprint("died on ", p.curr, ", expected ", lex.Token{typ, val})}
	}
}

type parseFn func(*parser) error

func Start(toks chan lex.Token) chan Node {
	p := &parser{
		toks:    toks,
		oldToks: make([]*lex.Token, 0),
		nodes:   make(chan Node),
	}
	go sourceFile(p)
	return p.nodes
}

// should the states return their list?... probably but not rn
// need a way to set something as optional from a top level
func sourceFile(p *parser) {
	packageClause(p)
	p.next()
	p.expect(lex.OpOrDelim, ";")
	importDecl(p)
	p.next()
	p.expect(lex.OpOrDelim, ";")
	close(p.nodes)
}

func packageClause(p *parser) {
	p.next()
	p.expect(lex.Keyword, "package")
	packageName(p)
}

func packageName(p *parser) {
	p.next()
	p.expect(lex.Identifier, "")
}

// import is optional
func importDecl(p *parser) {
	if err := p.next(); err != nil {
		return
	}
	if ok := p.accept(lex.Keyword, "import"); !ok {
		p.backup()
		return
	}
	// semicolon-delimited list of importSpecs
	if ok := p.accept(lex.OpOrDelim, "("); ok {
		importSpec(p)
		p.expect(lex.OpOrDelim, ";")
		return
	}
	// only one importSpec
	importSpec(p)
	return
}

// ImportSpec       = [ "." | PackageName ] ImportPath .
// TODO finish spec
func importSpec(p *parser) {
	importPath(p)
	return
}

func importPath(p *parser) {
	p.next()
	p.expect(lex.String, "")
	return
}
