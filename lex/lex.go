package lex

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

func (t token) String() string {
	return fmt.Sprintf("(%d %s)", t.typ, t.val)
}

type stateFn func(*lexer) stateFn

type lexer struct {
	name   string     // used for errors
	input  string     // string being scanned
	start  int        // start position of token
	pos    int        // current position of input
	width  int        // width of last rune read
	tokens chan token // channel of scanned tokens
}

const (
	tokenError tokenType = iota
	tokenEOF
	tokenIdentifier
	tokenSemicolon
	tokenString
)

const (
	semicolon  = ";"
	newline    = "\n"
	space      = " "
	tab        = "\t"
	whitespace = " \n\t"
	quote = "\""
	backslash = "\\"
	alphaLower = "abcdefghijklmnopqrstuvwxyz"
	alphaUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alpha      = alphaLower + alphaUpper
	numSansZero = "123456789"
	num        = numSansZero + "0"
	alphaNum   = alpha + num
)

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
	l.tokens <- token{typ: t, val: v}
}

func (l *lexer) emitErrorf(format string, a ...interface{}) {
	l.tokens <- token{typ: tokenError, val: fmt.Sprintf(format, a)}
}

func (l *lexer) emitError(a ...interface{}) {
	l.tokens <- token{typ: tokenError, val: fmt.Sprint(a)}
}

func (l *lexer) emitSemicolon() {
	l.start = l.pos
	l.tokens <- token{typ: tokenSemicolon}
}

func lexStart(l *lexer) stateFn {
	l.acceptRun(whitespace)
	l.ignore()
	if l.accept(alpha) {
		l.backup()
		return lexAlpha
	}
	// if l.accept(numSansZero) {
	// 	l.backup()
	// 	return lexDecimal
	// }
	// TODO rune literal
	if l.accept(quote) {
		l.backup()
		return lexString
	}
	return nil
}

func lexAlpha(l *lexer) stateFn {
	l.acceptRun(whitespace)
	l.ignore()
	if l.accept(alphaNum) {
		l.acceptRun(alphaNum)
		l.emit(tokenIdentifier)
		return lexStart
	}
	l.emitError("package name not found")
	return nil
}

func lexString(l *lexer) stateFn {
	if l.accept(quote) {
		// token.val does not contain the quote char
		l.ignore()
		return lexStringIn
	}
	l.emitError("expected '\"'")
	return nil
}

func lexStringIn(l *lexer) stateFn {
	l.acceptRunAllBut(quote + backslash)
	if l.peek() == '\\' {
		return lexStringBackslash
	}
	if l.peek() == '"' {
		return lexStringOut
	}
	return nil
}

func lexStringBackslash(l *lexer) stateFn {
	l.next() // eat backslash
	l.next() // eat next rune
	return lexStringIn
}

func lexStringOut(l *lexer) stateFn {
	l.emit(tokenString)
	l.next() // eat quote
	l.ignore()
	return lexStart
}
