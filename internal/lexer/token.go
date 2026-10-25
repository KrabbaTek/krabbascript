package lexer

import "fmt"

type TokenType int

const (
	// Keywords
	TokenVar TokenType = iota
	TokenVal
	TokenFunc
	TokenMod
	TokenReturn

	TokenIf
	TokenElsif
	TokenElse

	TokenWhile
	TokenFor
	TokenLoop
	TokenRepeat
	TokenUntil
	TokenWhen
	TokenStruct
	TokenImport

	// Types
	TokenBool
	TokenStr

	TokenI64
	TokenI32
	TokenI16
	TokenI8

	TokenU64
	TokenU32
	TokenU16
	TokenU8

	TokenF64
	TokenF32

	TokenAny
	TokenArr
	TokenDo

	// Symbols
	TokenPlus
	TokenPlusPlus
	TokenPlusEq

	TokenMinus
	TokenMinusMinus
	TokenMinusEq

	TokenMul // AKA TokenStar
	TokenMulEq

	TokenDiv
	TokenDivEq

	TokenEq
	TokenEqEq

	TokenNotEq
	TokenNot

	TokenBand
	TokenBor
	TokenBxor
	TokenBnot

	TokenBandEq
	TokenBorEq
	TokenBxorEq

	TokenArrow

	TokenComma
	TokenSemi
	TokenColon

	TokenOpenParen
	TokenClosedParen

	TokenOpenBrack
	TokenClosedBrack

	TokenOpenSqBrack
	TokenClosedSqBrack

	TokenDot
	TokenGreaterThan
	TokenLessThan

	TokenLessEq
	TokenGreaterEq

	// Compiler stuff
	TokenLiteral
	TokenStrLiteral
	TokenNumLiteral
	TokenFloatLiteral

	TokenNone
	TokenEof
	TokenEol
)

func (tt TokenType) String() string {
	tokens := [...]string{
		// Keywords
		TokenVar:    "var",
		TokenVal:    "val",
		TokenFunc:   "func",
		TokenMod:    "mod",
		TokenReturn: "return",
		TokenIf:     "if",
		TokenElsif:  "elsif",
		TokenElse:   "else",
		TokenWhile:  "while",
		TokenFor:    "for",
		TokenLoop:   "loop",
		TokenRepeat: "repeat",
		TokenUntil:  "until",
		TokenWhen:   "when",
		TokenStruct: "struct",
		TokenImport: "import",
		TokenDo:     "do",

		// Types
		TokenBool: "Bool",
		TokenStr:  "Str",
		TokenI64:  "I64",
		TokenI32:  "I32",
		TokenI16:  "I16",
		TokenI8:   "I8",
		TokenU64:  "U64",
		TokenU32:  "U32",
		TokenU16:  "U16",
		TokenU8:   "U8",
		TokenF64:  "F64",
		TokenF32:  "F32",
		TokenAny:  "Any",
		TokenArr:  "Arr",

		// Symbols
		TokenPlus:          "+",
		TokenPlusPlus:      "++",
		TokenPlusEq:        "+=",
		TokenMinus:         "-",
		TokenMinusMinus:    "--",
		TokenMinusEq:       "-=",
		TokenMul:           "*",
		TokenMulEq:         "*=",
		TokenDiv:           "/",
		TokenDivEq:         "/=",
		TokenEq:            "=",
		TokenEqEq:          "==",
		TokenNotEq:         "!=",
		TokenNot:           "!",
		TokenBand:          "&",
		TokenBor:           "|",
		TokenBxor:          "^",
		TokenBnot:          "~",
		TokenBandEq:        "&=",
		TokenBorEq:         "|=",
		TokenBxorEq:        "^=",
		TokenArrow:         "->",
		TokenComma:         ",",
		TokenSemi:          ";",
		TokenColon:         ":",
		TokenOpenParen:     "(",
		TokenClosedParen:   ")",
		TokenOpenBrack:     "{",
		TokenClosedBrack:   "}",
		TokenOpenSqBrack:   "[",
		TokenClosedSqBrack: "]",
		TokenDot:           ".",
		TokenLessThan:      "<",
		TokenGreaterThan:   ">",
		TokenLessEq:        "<=",
		TokenGreaterEq:     ">=",

		// Compiler stuff
		TokenLiteral:      "literal",
		TokenStrLiteral:   "str literal",
		TokenNumLiteral:   "num literal",
		TokenFloatLiteral: "float literal",
		TokenNone:         "none",
		TokenEof:          "eof",
		TokenEol:          "newline",
	}

	if int(tt) >= 0 && int(tt) < len(tokens) {
		return tokens[tt]
	}
	return fmt.Sprintf("unknown token(%d)", tt)
}

type Token struct {
	Type         TokenType
	Line, Column int

	Value string
}

func (t Token) String() string {
	if t.Value != "" {
		return fmt.Sprintf("%s %q", t.Type, t.Value)
	} else {
		return fmt.Sprintf("%s", t.Type)
	}
}
