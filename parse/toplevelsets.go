package parse

import (
	"github.com/samertm/chompy/lex"
)

var (
	tokString            = lex.Token{Typ: lex.String}
	tokIdentifier        = lex.Token{Typ: lex.Identifier}
	tokInt               = lex.Token{Typ: lex.Int}
	tokIf                = lex.Token{Typ: lex.Keyword, Val: "if"}
	tokElse              = lex.Token{Typ: lex.Keyword, Val: "else"}
	tokFor               = lex.Token{Typ: lex.Keyword, Val: "for"}
	tokGo                = lex.Token{Typ: lex.Keyword, Val: "go"}
	tokReturn            = lex.Token{Typ: lex.Keyword, Val: "return"}
	tokBreak             = lex.Token{Typ: lex.Keyword, Val: "break"}
	tokContinue          = lex.Token{Typ: lex.Keyword, Val: "continue"}
	tokGoto              = lex.Token{Typ: lex.Keyword, Val: "goto"}
	tokFallthrough       = lex.Token{Typ: lex.Keyword, Val: "fallthrough"}
	tokDefer             = lex.Token{Typ: lex.Keyword, Val: "defer"}
	tokRange             = lex.Token{Typ: lex.Keyword, Val: "range"}
	tokSemicolon         = lex.Token{Typ: lex.OpOrDelim, Val: ";"}
	tokDot               = lex.Token{Typ: lex.OpOrDelim, Val: "."}
	tokOpenParen         = lex.Token{Typ: lex.OpOrDelim, Val: "("}
	tokCloseParen        = lex.Token{Typ: lex.OpOrDelim, Val: ")"}
	tokOpenSquiggly      = lex.Token{Typ: lex.OpOrDelim, Val: "{"}
	tokCloseSquiggly     = lex.Token{Typ: lex.OpOrDelim, Val: "}"}
	tokComma             = lex.Token{Typ: lex.OpOrDelim, Val: ","}
	tokEqual             = lex.Token{Typ: lex.OpOrDelim, Val: "="}
	tokColonEqual        = lex.Token{Typ: lex.OpOrDelim, Val: ":="}
	tokDotDotDot         = lex.Token{Typ: lex.OpOrDelim, Val: "..."}
	tokLeftArrow         = lex.Token{Typ: lex.OpOrDelim, Val: "<-"}
	tokColon             = lex.Token{Typ: lex.OpOrDelim, Val: ":"}
	tokOpenSquareBrace   = lex.Token{Typ: lex.OpOrDelim, Val: "["}
	tokClosedSquareBrace = lex.Token{Typ: lex.OpOrDelim, Val: "]"}
	tokIncDec            = []lex.Token{
		lex.Token{Typ: lex.OpOrDelim, Val: "++"},
		lex.Token{Typ: lex.OpOrDelim, Val: "--"},
	}
	tokUnaryOp = []lex.Token{
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
	tokAssignOp = append(append([]lex.Token{
		tokEqual}, tokAddOp...), tokMulOp...)
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
	topTopLevelDecl = append(append([]lex.Token{},
		topDeclaration...), topFunctionDecl)
	topDeclaration = []lex.Token{
		topConstDecl,
		topTypeDecl,
		topVarDecl,
	}
	topConstDecl      = lex.Token{Typ: lex.Keyword, Val: "const"}
	topConstSpec      = topIdentifierList
	topIdentifierList = tokIdentifier
	topExpressionList = topExpression
	topExpression     = append([]lex.Token{}, topUnaryExpr...)
	topUnaryExpr      = append(append([]lex.Token{},
		topPrimaryExpr...), tokUnaryOp...)
	topPrimaryExpr = append(append([]lex.Token{topBuiltinCall}, topOperand...), topConversion)
	topOperand     = append([]lex.Token{topOperandName}, topLiteral...)
	topLiteral     = topBasicLit
	topBasicLit    = []lex.Token{
		tokInt,
		tokString,
	}
	topOperandName   = tokIdentifier
	topType          = []lex.Token{tokIdentifier, tokOpenParen}
	topTypeName      = tokIdentifier
	topTypeDecl      = lex.Token{Typ: lex.Keyword, Val: "type"}
	topTypeSpec      = tokIdentifier
	topVarDecl       = lex.Token{Typ: lex.Keyword, Val: "var"}
	topVarSpec       = topIdentifierList
	topFunctionDecl  = lex.Token{Typ: lex.Keyword, Val: "func"}
	topFunctionName  = tokIdentifier
	topFunction      = topSignature
	topFunctionBody  = topBlock
	topBlock         = tokOpenSquiggly
	topStatementList = topStatement
	topStatement     = append(append([]lex.Token{
		topLabeledStmt, topGoStmt, topReturnStmt, topBreakStmt,
		topContinueStmt, topGotoStmt, topFallthroughStmt, topBlock,
		topIfStmt, topForStmt, topDeferStmt,
	}, topDeclaration...), topSimpleStmt...)
	topSignature     = topParameters
	topResult        = append([]lex.Token{topParameters}, topType...)
	topParameters    = tokOpenParen
	topParameterList = topParameterDecl
	topParameterDecl = topIdentifierList
	// all simple statements start with an expression
	topSimpleStmt      = topExpression
	topLabeledStmt     = topLabel
	topLabel           = tokIdentifier
	topExpressionStmt  = topExpression
	topSendStmt        = topChannel
	topChannel         = topExpression
	topIncDecStmt      = topExpression
	topAssignment      = topExpression
	topIfStmt          = tokIf
	topForStmt         = tokFor
	topCondition       = topExpression
	topForClause       = append([]lex.Token{tokSemicolon}, topInitStmt...)
	topInitStmt        = topSimpleStmt
	topPostStmt        = topSimpleStmt
	topRangeClause     = append([]lex.Token{topIdentifierList}, topExpressionList...)
	topGoStmt          = tokGo
	topReturnStmt      = tokReturn
	topBreakStmt       = tokBreak
	topContinueStmt    = tokContinue
	topGotoStmt        = tokGoto
	topFallthroughStmt = tokFallthrough
	topDeferStmt       = tokDefer
	topShortVarDecl    = topIdentifierList
	topConversion      = topType
	topBuiltinCall     = tokIdentifier
	topBuiltinArgs     = topType
	topSelector        = tokDot
	topIndex           = tokOpenSquareBrace
	topSlice           = tokOpenSquareBrace
	topTypeAssertion   = tokDot
	topCall            = tokOpenParen
	topArgumentList    = topExpressionList
)
