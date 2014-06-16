package parse

import (
	"github.com/samertm/chompy/lex"
)

var (
	tokSemicolon  = lex.Token{Typ: lex.OpOrDelim, Val: ";"}
	tokIdentifier = lex.Token{Typ: lex.Identifier}
	tokDot        = lex.Token{Typ: lex.OpOrDelim, Val: "."}
	tokString     = lex.Token{Typ: lex.String}
	tokOpenParen  = lex.Token{Typ: lex.OpOrDelim, Val: "("}
	tokCloseParen = lex.Token{Typ: lex.OpOrDelim, Val: ")"}
	tokComma      = lex.Token{Typ: lex.OpOrDelim, Val: ","}
)

var (
	topPackageClause = lex.Token{Typ: lex.Keyword, Val: "package"}
	topPackageName   = tokIdentifier
	topImportDecl    = lex.Token{Typ: lex.Keyword, Val: "import"}
	topImportSpec    = []lex.Token{
		tokDot,
		topPackageName,
		topImportPath,
	}
	topImportPath   = tokString
	topTopLevelDecl = append([]lex.Token{}, topDeclaration...)
	topDeclaration  = []lex.Token{
		topConstDecl,
	}
	topConstDecl      = lex.Token{Typ: lex.Keyword, Val: "const"}
	topConstSpec      = topIdentifierList
	topIdentifierList = tokIdentifier
)
