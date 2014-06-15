package lex

import (
	"strings"
)

type stateFn func(*lexer) stateFn

const eof = -1

const (
	Error TokenType = iota
	EOF
	Keyword
	OpOrDelim
	Identifier
	String
	Int
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
	comment               = "//"
	alphaLower            = "abcdefghijklmnopqrstuvwxyz"
	alphaUpper            = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alpha                 = alphaLower + alphaUpper
	letter                = alpha + "_"
	numSansZero           = "123456789"
	num                   = numSansZero + "0"
	alphaNum              = alpha + num
	letterNum             = letter + num
)

// sorted by length
var opDelims = [...]string{"&^=", "...", "<<=", ">>=", "+=", "&^", "&=", "--", "&&", "%=", "==", ">>", "!=", ":=", "-=", "++", "|=", "/=", "||", "<<", "<=", ">=", "*=", "<-", "^=", "+", ":", "&", ".", "(", "!", ")", "%", "-", ";", "|", ",", "<", "=", "[", "/", "]", "}", "*", "{", "^", ">"}

var keywords = [...]string{"break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type", "continue", "for", "import", "return", "var"}

func lexStart(l *lexer) stateFn {
	l.acceptRun(whitespaceSansNewline)
	l.ignore()
	if strings.HasPrefix(l.input[l.pos:], comment) {
		return lexComment
	}
	if l.accept(letter) {
		l.backup()
		return lexLetter
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
	// idea: group every 1 char operator together.
	if l.accept("+&=!()-|<[]*^<>{}/:,;%>.") {
		l.backup()
		return lexOpOrDelim
	}
	if l.next() == eof {
		l.backup()
		return lexEof
	}
	return nil
}

func semicolonRule(l *lexer) bool {
	if l.lastToken == nil {
		return false
	}
	validSemicolonInsert := false
	switch l.lastToken.Typ {
	case OpOrDelim:
		if l.lastToken.Val != "++" &&
			l.lastToken.Val != "--" &&
			l.lastToken.Val != ")" &&
			l.lastToken.Val != "]" &&
			l.lastToken.Val != "}" {
			break
		}
		validSemicolonInsert = true
	case Keyword:
		if l.lastToken.Val != "break" &&
			l.lastToken.Val != "continue" &&
			l.lastToken.Val != "fallthrough" &&
			l.lastToken.Val != "return" {
			break
		}
		validSemicolonInsert = true
	case Identifier:
		validSemicolonInsert = true
	case Int:
		validSemicolonInsert = true
	case String:
		validSemicolonInsert = true
	}
	if validSemicolonInsert {
		l.emitSemicolon()
	}
	return validSemicolonInsert
}

func lexNewline(l *lexer) stateFn {
	l.accept(newline)
	l.ignore()
	semicolonRule(l)
	l.lastToken = nil
	return lexStart
}

func lexLetter(l *lexer) stateFn {
	if l.accept(letter) {
		l.acceptRun(letterNum)
		if isKeyword(l.val()) {
			l.emit(Keyword)
		} else {
			l.emit(Identifier)
		}
		return lexStart
	}
	l.emitError("whoops")
	return nil
}

func lexDecimal(l *lexer) stateFn {
	if l.accept(num) {
		l.acceptRun(num)
		l.emit(Int)
		return lexStart
	}
	l.emitError("expected integer")
	return nil
}

func lexString(l *lexer) stateFn {
	if l.accept(quote) {
		// token.Val does not contain the quote char
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
	l.emit(String)
	l.next() // eat quote
	l.ignore()
	return lexStart
}

func lexOpOrDelim(l *lexer) stateFn {
	for _, od := range opDelims {
		if strings.HasPrefix(l.input[l.pos:], od) {
			l.pos += len(od)
			l.emit(OpOrDelim)
			return lexStart
		}
	}
	return nil
}

func lexEof(l *lexer) stateFn {
	if l.next() == eof {
		semicolonRule(l)
		l.emitEof()
		return nil
	}
	l.emitError("expected eof")
	return nil
}

func lexComment(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], comment) {
		l.acceptRunAllBut(newline)
		l.ignore()
		return lexNewline
	}
	l.emitError("error handling comment")
	return nil
}
