package parse

import (
	"github.com/samertm/chompy/lex"
)

var (
	tokString     = lex.Token{Typ: lex.String}
	tokIdentifier = lex.Token{Typ: lex.Identifier}
	tokInt        = lex.Token{Typ: lex.Int}
	tokSemicolon  = lex.Token{Typ: lex.OpOrDelim, Val: ";"}
	tokDot        = lex.Token{Typ: lex.OpOrDelim, Val: "."}
	tokOpenParen  = lex.Token{Typ: lex.OpOrDelim, Val: "("}
	tokCloseParen = lex.Token{Typ: lex.OpOrDelim, Val: ")"}
	tokComma      = lex.Token{Typ: lex.OpOrDelim, Val: ","}
	tokEqual      = lex.Token{Typ: lex.OpOrDelim, Val: "="}
	tokUnaryOp    = []lex.Token{
		lex.Token{Typ: lex.OpOrDelim, Val: "+"},
		lex.Token{Typ: lex.OpOrDelim, Val: "-"},
		lex.Token{Typ: lex.OpOrDelim, Val: "!"},
		lex.Token{Typ: lex.OpOrDelim, Val: "^"},
		lex.Token{Typ: lex.OpOrDelim, Val: "*"},
		lex.Token{Typ: lex.OpOrDelim, Val: "&"},
		lex.Token{Typ: lex.OpOrDelim, Val: "<-"},
	}
	tokMulOp = []lex.Token{
		lex.Token{Typ: lex.OpOrDelim, Val: "*"},
		lex.Token{Typ: lex.OpOrDelim, Val: "/"},
		lex.Token{Typ: lex.OpOrDelim, Val: "%"},
		lex.Token{Typ: lex.OpOrDelim, Val: "<<"},
		lex.Token{Typ: lex.OpOrDelim, Val: ">>"},
		lex.Token{Typ: lex.OpOrDelim, Val: "&"},
		lex.Token{Typ: lex.OpOrDelim, Val: "&^"},
	}
	tokAddOp = []lex.Token{
		lex.Token{Typ: lex.OpOrDelim, Val: "+"},
		lex.Token{Typ: lex.OpOrDelim, Val: "-"},
		lex.Token{Typ: lex.OpOrDelim, Val: "|"},
		lex.Token{Typ: lex.OpOrDelim, Val: "^"},
	}
	tokRelOp = []lex.Token{
		lex.Token{Typ: lex.OpOrDelim, Val: "=="},
		lex.Token{Typ: lex.OpOrDelim, Val: "!="},
		lex.Token{Typ: lex.OpOrDelim, Val: "<"},
		lex.Token{Typ: lex.OpOrDelim, Val: "<="},
		lex.Token{Typ: lex.OpOrDelim, Val: ">"},
		lex.Token{Typ: lex.OpOrDelim, Val: ">="},
	}
	tokBinaryOp = append(append(append([]lex.Token{
		lex.Token{Typ: lex.OpOrDelim, Val: "||"},
		lex.Token{Typ: lex.OpOrDelim, Val: "&&"},
	}, tokMulOp...), tokAddOp...), tokRelOp...)
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
		topTypeDecl,
	}
	topConstDecl      = lex.Token{Typ: lex.Keyword, Val: "const"}
	topConstSpec      = topIdentifierList
	topIdentifierList = tokIdentifier
	topExpressionList = topExpression
	topExpression     = append([]lex.Token{}, topUnaryExpr...)
	topUnaryExpr      = append(append([]lex.Token{},
		topPrimaryExpr...), tokUnaryOp...)
	topPrimaryExpr = append([]lex.Token{}, topOperand...)
	topOperand     = append([]lex.Token{topOperandName}, topLiteral...)
	topLiteral     = topBasicLit
	topBasicLit    = []lex.Token{
		tokInt,
		tokString,
	}
	topOperandName = tokIdentifier
	topType        = []lex.Token{tokIdentifier, tokOpenParen}
	topTypeName    = tokIdentifier
	topTypeDecl    = lex.Token{Typ: lex.Keyword, Val: "type"}
	topTypeSpec    = tokIdentifier
)
