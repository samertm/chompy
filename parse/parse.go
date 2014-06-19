package parse

// - come up with data structure for top-level terminals

import (
	"github.com/samertm/chompy/lex"

	"log"
)

type parser struct {
	toks    chan lex.Token
	oldToks []*lex.Token
	// IDEA: stack of channels?
	// that's pretty good...
	// always emit to the channel on the top
	// push chan when going down a level

	// set when you want to backtrack
	// default to false
	recordTokens   bool
	recordedTokens []*lex.Token
	nodes          chan Node
	ast            Tree
}

func (p *parser) next() *lex.Token {
	if len(p.oldToks) != 0 {
		curr := p.oldToks[len(p.oldToks)-1]
		p.oldToks = p.oldToks[:len(p.oldToks)-1]
		if p.recordTokens {
			p.recordedTokens = append(p.recordedTokens, curr)
		}
		return curr
	}
	if t, ok := <-p.toks; ok {
		curr := &t
		if curr.Typ == lex.Error {
			log.Fatal("error lexing: ", curr)
			return nil
		}
		if p.recordTokens {
			p.recordedTokens = append(p.recordedTokens, curr)
		}
		return curr
	}
	log.Fatal("token stream closed")
	return nil
}

// It is illegal to push a token other than the one that was just
// recieved from next() when recording tokens
func (p *parser) push(t *lex.Token) {
	if t == nil {
		log.Fatal("bad push")
	}
	p.oldToks = append(p.oldToks, t)
	if p.recordTokens {
		// pop last token off slice
		p.recordedTokens = p.recordedTokens[:len(p.oldToks)-1]
	}
}

// turns off recording
func (p *parser) pushRecordedTokens() {
	if !p.recordTokens {
		return
	}
	for i := len(p.recordedTokens) - 1; i >= 0; i-- {
		p.oldToks = append(p.oldToks, p.recordedTokens[i])
	}
	p.recordedTokens = make([]*lex.Token, 0)
	p.recordTokens = false
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

func (p *parser) expect(tok lex.Token) *Erro {
	if p.accept(tok) {
		return nil
	}
	return &Erro{"expected " + tok.String() + " recieved " + p.peek().String()}
}
