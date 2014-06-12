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

const eof = -1

const (
	tokenError tokenType = iota
	tokenEOF
	tokenKeyword
	tokenOperator
	tokenDelimiter
	tokenIdentifier
	tokenSemicolon
	tokenString
	tokenInt
)

const (
	semicolon             = ";"
	newline               = "\n"
	space                 = " "
	tab                   = "\t"
	whitespaceSansNewline = space + tab
	whitespace            = whitespaceSansNewline + newline
	quote                 = "\""
	backslash             = "\\"
	alphaLower            = "abcdefghijklmnopqrstuvwxyz"
	alphaUpper            = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alpha                 = alphaLower + alphaUpper
	numSansZero           = "123456789"
	num                   = numSansZero + "0"
	alphaNum              = alpha + num
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
	l.tokens <- token{typ: tokenSemicolon}
}

// peeks at the lexer's current value, without emitting it or changing
// the position.
func (l *lexer) val() string {
	return l.input[l.start:l.pos]
}

func lexStart(l *lexer) stateFn {
	l.acceptRun(whitespaceSansNewline)
	l.ignore()
	if l.accept(alpha) {
		l.backup()
		return lexAlpha
	}
	if l.accept(numSansZero) {
		l.backup()
		return lexDecimal
	}
	// TODO rune literal
	if l.accept(quote) {
		l.backup()
		return lexString
	}
	if l.accept(newline) {
		l.backup()
		return lexNewline
	}
	return nil
}

// TODO get this to not emit semicolons on blank lines
func lexNewline(l *lexer) stateFn {
	l.accept(newline)
	l.ignore()
	if l.lastToken == nil {
		return lexStart
	}
	validSemicolonInsert := false
	switch l.lastToken.typ {
	case tokenOperator:
		if l.lastToken.val != "++" &&
			l.lastToken.val != "--" {
			break
		}
		validSemicolonInsert = true
	case tokenDelimiter:
		if l.lastToken.val != ")" &&
			l.lastToken.val != "]" &&
			l.lastToken.val != "}" {
			break
		}
		validSemicolonInsert = true
	case tokenKeyword:
		if l.lastToken.val != "break" &&
			l.lastToken.val != "continue" &&
			l.lastToken.val != "fallthrough" &&
			l.lastToken.val != "return" {
			break
		}
		validSemicolonInsert = true
	case tokenIdentifier:
		validSemicolonInsert = true
	case tokenInt:
		validSemicolonInsert = true
	case tokenSemicolon:
		validSemicolonInsert = true
	case tokenString:
		validSemicolonInsert = true
	}
	if validSemicolonInsert {
		l.emitSemicolon()
	}
	l.lastToken = nil
	return lexStart
}

func lexAlpha(l *lexer) stateFn {
	if l.accept(alphaNum) {
		l.acceptRun(alphaNum)
		isKeyword := false
		// this turned out jankier than I thought it would..
		switch l.val() {
		case "break":
			isKeyword = true
		case "default":
			isKeyword = true
		case "func":
			isKeyword = true
		case "interface":
			isKeyword = true
		case "select":
			isKeyword = true
		case "case":
			isKeyword = true
		case "defer":
			isKeyword = true
		case "go":
			isKeyword = true
		case "map":
			isKeyword = true
		case "struct":
			isKeyword = true
		case "chan":
			isKeyword = true
		case "else":
			isKeyword = true
		case "goto":
			isKeyword = true
		case "package":
			isKeyword = true
		case "switch":
			isKeyword = true
		case "const":
			isKeyword = true
		case "fallthrough":
			isKeyword = true
		case "if":
			isKeyword = true
		case "range":
			isKeyword = true
		case "type":
			isKeyword = true
		case "continue":
			isKeyword = true
		case "for":
			isKeyword = true
		case "import":
			isKeyword = true
		case "return":
			isKeyword = true
		case "var":
			isKeyword = true
		}
		if isKeyword {
			l.emit(tokenKeyword)
		} else {
			l.emit(tokenIdentifier)
		}
		return lexStart
	}
	l.emitError("whoops")
	return nil
}

func lexDecimal(l *lexer) stateFn {
	if l.accept(num) {
		l.acceptRun(num)
		l.emit(tokenInt)
		return lexStart
	}
	l.emitError("expected integer")
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

// TODO turn \[A-Z] into char code
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
