package lex

import (
	"strings"
)

type stateFn func(*lexer) stateFn

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

	return nil
}

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

func lexLetter(l *lexer) stateFn {
	if l.accept(letter) {
		l.acceptRun(letterNum)
		if isKeyword(l.val()) {
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

func lexOpOrDelim(l *lexer) stateFn {
	for _, od := range opDelims {
		if strings.HasPrefix(l.input[l.pos:], od) {
			l.pos += len(od)
			l.emit(tokenDelimiter)
			return lexStart
		}
	}
	return nil
}
