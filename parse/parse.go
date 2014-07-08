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

	// stream of recorded tokens
	recordedTokens []*lex.Token
	// index in recorded tokens each tracker is at
	trackers []int
	nodes    chan Node
	ast      Tree
}

func (p *parser) next() *lex.Token {
	if len(p.oldToks) != 0 {
		curr := p.oldToks[len(p.oldToks)-1]
		p.oldToks = p.oldToks[:len(p.oldToks)-1]
		if len(p.trackers) > 0 {
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
		if len(p.trackers) > 0 {
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
	if len(p.trackers) > 0 {
		// pop last token off recordedTokens
		p.recordedTokens = p.recordedTokens[:len(p.recordedTokens)-1]
	}
}

// hook up a tracker for backtracking
func (p *parser) hookTracker() {
	//fmt.Println("HOOK TRACKER")
	// add tracker at current index in recordedTokens stream
	p.trackers = append(p.trackers, len(p.recordedTokens))
}

// unhook a tracker
func (p *parser) unhookTracker() {
	//fmt.Println("UNHOOK TRACKER")
	if len(p.trackers) == 0 {
		log.Fatal("Error: unhookTracker called with zero trackers")
	} else if len(p.trackers) == 1 {
		// erase trackers, recordedTokens
		p.trackers = make([]int, 0)
		p.recordedTokens = make([]*lex.Token, 0)
	} else {
		// remove tracker, keep recorded tokens intact
		p.trackers = p.trackers[:len(p.trackers)-1]
	}
}

// does not unhook tracker
func (p *parser) backtrack() {
	//fmt.Println(strconv.Itoa(len(p.trackers)), "many backtrackers")
	//fmt.Println("BACKTRACKING START")
	if len(p.trackers) == 0 {
		log.Fatal("Error: backtrack called with zero trackers")
	}
	// invariant: i is always the index of the token we want to push on the stream
	// starts at the last index in recordedTokens
	// ends at the index where the tracker began tracking
	start := len(p.recordedTokens) - 1
	end := p.trackers[len(p.trackers)-1]
	for i := start; i >= end; i-- {
		t := p.recordedTokens[i]
		//fmt.Println("BACK: ", t)
		p.oldToks = append(p.oldToks, t)
	}
	// remove backtracked tokens from recordedTokens
	p.recordedTokens = p.recordedTokens[:end]
	//fmt.Println("BACKTRACKING END")
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
