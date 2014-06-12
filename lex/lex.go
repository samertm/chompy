package lex

/* tasks remaining (ordered by significance):
 * - handle raw strings
 * - handle non-line comments
 * - other stuff probably
 */

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

type state int

type tokenType int

type token struct {
	typ tokenType
	val string
}

// for debugging purposes
func (t token) String() string {
	return fmt.Sprintf("(%s %s)\n", typeName[t.typ], t.val)
}

// for debugging purposes
var typeName = map[tokenType]string{
	tokenError:      "Error",
	tokenEOF:        "EOF",
	tokenKeyword:    "Keyword",
	tokenOpOrDelim:  "OpOrDelim",
	tokenIdentifier: "Identifier",
	tokenString:     "String",
	tokenInt:        "Int",
}

type lexer struct {
	name  string // used for errors
	input string // string being scanned
	start int    // start position of token
	pos   int    // current position of input
	width int    // width of last rune read
	// The last token processed on the line for newline insertion.
	// Error tokens are not stored.
	lastToken *token
	tokens    chan token // channel of scanned tokens
}

func isKeyword(val string) bool {
	for _, k := range keywords {
		if val == k {
			return true
		}
	}
	return false
}

func Lex(name, input string) (*lexer, chan token) {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make(chan token),
	}
	go l.run()
	return l, l.tokens
}

func (l *lexer) run() {
	for state := lexStart; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// reads & returns the next rune, steps width forward
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, s := utf8.DecodeRuneInString(l.input[l.pos:])
	if r == utf8.RuneError && s == 1 {
		log.Fatal("input error")
	}
	l.width = s
	l.pos += l.width
	return r
}

// can only be called once after each next
func (l *lexer) backup() {
	l.pos -= l.width
}

// accepts single rune in accepted
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) != -1 {
		return true
	}
	l.backup()
	return false
}

// accepts all runes in valid
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) != -1 {
	}
	l.backup()
}

func (l *lexer) acceptAllBut(invalid string) bool {
	for strings.IndexRune(invalid, l.next()) == -1 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRunAllBut(invalid string) {
	for strings.IndexRune(invalid, l.next()) == -1 {
	}
	l.backup()
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) emit(t tokenType) {
	v := l.input[l.start:l.pos]
	l.start = l.pos
	l.lastToken = &token{typ: t, val: v}
	l.tokens <- *l.lastToken
}

func (l *lexer) emitErrorf(format string, a ...interface{}) {
	l.tokens <- token{typ: tokenError, val: fmt.Sprintf(format, a)}
}

func (l *lexer) emitError(a ...interface{}) {
	l.tokens <- token{typ: tokenError, val: fmt.Sprint(a)}
}

func (l *lexer) emitSemicolon() {
	l.start = l.pos
	l.lastToken = &token{typ: tokenOpOrDelim, val: ";"}
	l.tokens <- token{typ: tokenOpOrDelim, val: ";"}
}

// peeks at the lexer's current value, without emitting it or changing
// the position.
func (l *lexer) val() string {
	return l.input[l.start:l.pos]
}
