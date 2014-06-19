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
	trackers [][]*lex.Token
	nodes    chan Node
	ast      Tree
}

func (p *parser) next() *lex.Token {
	if len(p.oldToks) != 0 {
		curr := p.oldToks[len(p.oldToks)-1]
		p.oldToks = p.oldToks[:len(p.oldToks)-1]
		if len(p.trackers) > 0 {
			p.trackers[len(p.trackers)-1] = append(p.trackers[len(p.trackers)-1], curr)
		}
		return curr
	}
	if t, ok := <-p.toks; ok {
		curr := &t
		if curr.Typ == lex.Error {
			log.Fatal("error lexing: ", curr)
			return nil
		}
		if len(p.trackers) > 0 {
			p.trackers[len(p.trackers)-1] = append(p.trackers[len(p.trackers)-1], curr)
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

// hook up a tracker for backtracking
func (p *parser) hookTracker() {
	p.trackers = append(p.trackers, make([]*lex.Token, 0))
}

// unhook a tracker
func (p *parser) unhookTracker() {
	if len(p.trackers) == 0 {
		log.Fatal("Error: unhookTracker called with zero trackers")
	}
	p.trackers = p.trackers[:len(p.trackers)-1]
}

// does not unhook tracker
func (p *parser) backtrack() {
	if len(p.trackers) == 0 {
		log.Fatal("Error: backtrack called with zero trackers")
	}
	for i := len(p.recordedTokens) - 1; i >= 0; i-- {
		p.oldToks = append(p.oldToks, p.trackers[len(p.trackers)-1][i])
	}
	p.trackers[len(p.trackers)-1] = make([]*lex.Token, 0)
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
