package lexer

import (
	"testing"
)

type expectedToken struct {
	expectedType  TokenType
	expectedValue string
}

func testTokenSequence(t *testing.T, source string, expected []expectedToken) {
	t.Helper()

	l := NewLexerFromStr(source, "test")
	tokens := l.Scan()

	if len(tokens) != len(expected) {
		t.Fatalf("token count mismatch: got %d tokens, want %d", len(tokens), len(expected))
	}

	for i, want := range expected {
		got := tokens[i]

		if got.Type != want.expectedType {
			t.Errorf("[%d] token type mismatch: got %v (%s), want %v (%s)",
				i, got.Type, got.Type, want.expectedType, want.expectedType)
		}

		if want.expectedValue != "" && got.Value != want.expectedValue {
			t.Errorf("[%d] token value mismatch: got %q, want %q",
				i, got.Value, want.expectedValue)
		}
	}
}

func TestEofToken(t *testing.T) {
	testTokenSequence(t, "", []expectedToken{
		{TokenEof, ""},
	})
}

func TestEolToken(t *testing.T) {
	testTokenSequence(t, "\n", []expectedToken{
		{TokenEol, ""},
		{TokenEof, ""},
	})
}

func TestKeywords(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect TokenType
	}{
		{"var", "var", TokenVar},
		{"val", "val", TokenVal},
		{"func", "func", TokenFunc},
		{"mod", "mod", TokenMod},
		{"ret", "ret", TokenRet},
		{"if", "if", TokenIf},
		{"elsif", "elsif", TokenElsif},
		{"else", "else", TokenElse},
		{"while", "while", TokenWhile},
		{"for", "for", TokenFor},
		{"loop", "loop", TokenLoop},
		{"repeat", "repeat", TokenRepeat},
		{"until", "until", TokenUntil},
		{"when", "when", TokenWhen},
		{"struct", "struct", TokenStruct},
		{"import", "import", TokenImport},
		{"do", "do", TokenDo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTokenSequence(t, tt.input, []expectedToken{
				{tt.expect, ""},
				{TokenEof, ""},
			})
		})
	}
}

func TestBuiltInTypes(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect TokenType
	}{
		{"Bool", "Bool", TokenBool},
		{"Str", "Str", TokenStr},
		{"I64", "I64", TokenI64},
		{"I32", "I32", TokenI32},
		{"I16", "I16", TokenI16},
		{"I8", "I8", TokenI8},
		{"U64", "U64", TokenU64},
		{"U32", "U32", TokenU32},
		{"U16", "U16", TokenU16},
		{"U8", "U8", TokenU8},
		{"F64", "F64", TokenF64},
		{"F32", "F32", TokenF32},
		{"Any", "Any", TokenAny},
		{"Arr", "Arr", TokenArr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTokenSequence(t, tt.input, []expectedToken{
				{tt.expect, ""},
				{TokenEof, ""},
			})
		})
	}
}

func TestSymbolsAndOperators(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect TokenType
	}{
		{"plus", "+", TokenPlus},
		{"increment", "++", TokenPlusPlus},
		{"plus_eq", "+=", TokenPlusEq},
		{"minus", "-", TokenMinus},
		{"decrement", "--", TokenMinusMinus},
		{"minus_eq", "-=", TokenMinusEq},
		{"multiply", "*", TokenMul},
		{"mul_eq", "*=", TokenMulEq},
		{"divide", "/", TokenDiv},
		{"div_eq", "/=", TokenDivEq},
		{"assign", "=", TokenEq},
		{"equal", "==", TokenEqEq},
		{"not_equal", "!=", TokenNotEq},
		{"not", "!", TokenNot},
		{"bitwise_and", "&", TokenBand},
		{"bitwise_or", "|", TokenBor},
		{"bitwise_xor", "^", TokenBxor},
		{"bitwise_not", "~", TokenBnot},
		{"band_eq", "&=", TokenBandEq},
		{"bor_eq", "|=", TokenBorEq},
		{"bxor_eq", "^=", TokenBxorEq},
		{"arrow", "->", TokenArrow},
		{"comma", ",", TokenComma},
		{"semi", ";", TokenSemi},
		{"colon", ":", TokenColon},
		{"open_paren", "(", TokenOpenParen},
		{"closed_paren", ")", TokenClosedParen},
		{"open_brace", "{", TokenOpenBrack},
		{"closed_brace", "}", TokenClosedBrack},
		{"open_bracket", "[", TokenOpenSqBrack},
		{"closed_bracket", "]", TokenClosedSqBrack},
		{"dot", ".", TokenDot},
		{"less_than", "<", TokenLessThan},
		{"greater_than", ">", TokenGreaterThan},
		{"less_eq", "<=", TokenLessEq},
		{"greater_eq", ">=", TokenGreaterEq},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTokenSequence(t, tt.input, []expectedToken{
				{tt.expect, ""},
				{TokenEof, ""},
			})
		})
	}
}

func TestVariableDeclaration(t *testing.T) {
	input := `var krabba: I32 = 0;`

	expected := []expectedToken{
		{TokenVar, ""},
		{TokenLiteral, "krabba"},
		{TokenColon, ""},
		{TokenI32, ""},
		{TokenEq, ""},
		{TokenNumLiteral, "0"},
		{TokenSemi, ""},
		{TokenEof, ""},
	}

	testTokenSequence(t, input, expected)
}

func TestFunctionDeclaration(t *testing.T) {
	input := `func add(a: I32, b: I32) -> I32 { ret a + b; }`

	expected := []expectedToken{
		{TokenFunc, ""},
		{TokenLiteral, "add"},
		{TokenOpenParen, ""},
		{TokenLiteral, "a"},
		{TokenColon, ""},
		{TokenI32, ""},
		{TokenComma, ""},
		{TokenLiteral, "b"},
		{TokenColon, ""},
		{TokenI32, ""},
		{TokenClosedParen, ""},
		{TokenArrow, ""},
		{TokenI32, ""},
		{TokenOpenBrack, ""},
		{TokenRet, ""},
		{TokenLiteral, "a"},
		{TokenPlus, ""},
		{TokenLiteral, "b"},
		{TokenSemi, ""},
		{TokenClosedBrack, ""},
		{TokenEof, ""},
	}

	testTokenSequence(t, input, expected)
}

func TestStringLiteral(t *testing.T) {
	input := `"Krabba!"`

	expected := []expectedToken{
		{TokenStrLiteral, "Krabba!"},
		{TokenEof, ""},
	}

	testTokenSequence(t, input, expected)
}
