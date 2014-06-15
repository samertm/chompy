package parse

// - come up with data structure for top-level terminals

import (
	"github.com/samertm/chompy/lex"

	"log"
)

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
		// fmt.Println("oldToks:", curr)
		return curr
	}
	if t, ok := <-p.toks; ok {
		curr := &t
		if curr.Typ == lex.Error {
			log.Fatal("error lexing: ", curr)
			return nil
		}
		//fmt.Println("chan:   ", curr)
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

