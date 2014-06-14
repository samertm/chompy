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

type TokenType int

type Token struct {
	Typ TokenType
	Val string
}

// for debugging purposes
func (t Token) String() string {
	return fmt.Sprintf("(%s %s)", typeName[t.Typ], t.Val)
}

// for debugging purposes
var typeName = map[TokenType]string{
	Error:      "Error",
	EOF:        "EOF",
	Keyword:    "Keyword",
	OpOrDelim:  "OpOrDelim",
	Identifier: "Identifier",
	String:     "String",
	Int:        "Int",
}

type lexer struct {
	name  string // used for errors
	input string // string being scanned
	start int    // start position of token
	pos   int    // current position of input
	width int    // width of last rune read
	// The last token processed on the line for newline insertion.
	// Error tokens are not stored.
	lastToken *Token
	Tokens    chan Token // channel of scanned Tokens
}

func isKeyword(val string) bool {
	for _, k := range keywords {
		if val == k {
			return true
		}
	}
	return false
}

func Lex(name, input string) (*lexer, chan Token) {
	l := &lexer{
		name:   name,
		input:  input,
		Tokens: make(chan Token),
	}
	go l.run()
	return l, l.Tokens
}

func (l *lexer) run() {
	for state := lexStart; state != nil; {
		state = state(l)
	}
	close(l.Tokens)
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

func (l *lexer) emit(t TokenType) {
	v := l.input[l.start:l.pos]
	l.start = l.pos
	l.lastToken = &Token{Typ: t, Val: v}
	l.Tokens <- *l.lastToken
}

func (l *lexer) emitErrorf(format string, a ...interface{}) {
	l.Tokens <- Token{Typ: Error, Val: fmt.Sprintf(format, a)}
}

func (l *lexer) emitError(a ...interface{}) {
	l.Tokens <- Token{Typ: Error, Val: fmt.Sprint(a)}
}

func (l *lexer) emitSemicolon() {
	l.start = l.pos
	l.lastToken = &Token{Typ: OpOrDelim, Val: ";"}
	l.Tokens <- Token{Typ: OpOrDelim, Val: ";"}
}

// peeks at the lexer's current value, without emitting it or changing
// the position.
func (l *lexer) val() string {
	return l.input[l.start:l.pos]
}
