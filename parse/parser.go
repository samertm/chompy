package parse

import (
	"errors"
	"log"

	"github.com/samertm/chompy/lex"
)

// parser type keeps track of the state of the parser for the state
// functions in grammar.go. None of its fields should be accessed
// directly.
type parser struct {
	toks    chan lex.Token
	oldToks []*lex.Token
	// stream of recorded tokens
	recordedTokens []*lex.Token
	// index in recorded tokens that each tracker is on
	trackers []int
	errs     pErrors
	// Maps a tracker index to the earliest errors index.
	trackerToErrors map[int]int
}

type pErrors []string

func (e pErrors) Error() string {
	if len(e) == 0 {
		return "No errors."
	}
	str := make([]byte, 0)
	for i, err := range e {
		if i != 0 {
			str = append(str, '\n')
		}
		str = append(str, err...)
	}
	return string(str)
}

func newParser(toks chan lex.Token) *parser {
	return &parser{
		toks:            toks,
		oldToks:         make([]*lex.Token, 0),
		errs:            make([]string, 0),
		trackerToErrors: make(map[int]int),
	}
}

// gets the next token in the stream
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

// pushes a token onto the stream
// it is illegal to push a token other than the one that was just
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
// does not backtrack
func (p *parser) unhookTracker() {
	//fmt.Println("UNHOOK TRACKER")
	if len(p.trackers) == 0 {
		log.Fatal("Error: unhookTracker called with zero trackers")
	} else if len(p.trackers) == 1 {
		// erase trackers, recordedTokens
		p.trackers = make([]int, 0)
		p.recordedTokens = make([]*lex.Token, 0)
		p.trackerToErrors = make(map[int]int)
	} else {
		// remove tracker, keep recorded tokens intact
		p.trackers = p.trackers[:len(p.trackers)-1]
		delete(p.trackerToErrors, len(p.trackers))
	}
}

// pushes all tokens of the tracker onto the stream
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
	// remove associated errors
	if erroridx, ok := p.trackerToErrors[len(p.trackers)-1]; ok {
		p.errs = p.errs[:erroridx]
	}
}

// peek looks at the next token without modifying the stream
func (p *parser) peek() *lex.Token {
	t := p.next()
	p.push(t)
	return t
}

// accept returns true if the next token in the stream matches
// all of the tokens passed in as args. accept does not modify
// the stream.
func (p *parser) accept(toks ...lex.Token) bool {
	nextTok := p.peek()
	for _, t := range toks {
		if lex.TokenEquiv(*nextTok, t) {
			return true
		}
	}
	return false
}

// expect returns an error if it cannot accept the token that is
// passed in. Use it as a stronger version of accept. expect does
// not modify the stream.
func (p *parser) expect(tok lex.Token) error {
	if p.accept(tok) {
		return nil
	}
	return errors.New("expected " + tok.String() + " recieved " + p.peek().String())
}

func (p *parser) addError(e string) {
	p.errs = append(p.errs, e)
	if len(p.trackers) != 0 {
		i := len(p.trackers) - 1
		if _, ok := p.trackerToErrors[i]; !ok {
			p.trackerToErrors[i] = len(p.errs) - 1
		}
	}
}

// returns true if the current tracker has not hit any errors
func (p *parser) valid() bool {
	_, ok := p.trackerToErrors[len(p.trackers)-1]
	return !ok
}
