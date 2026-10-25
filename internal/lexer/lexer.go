package lexer

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

type Lexer struct {
	line, column, pos int
	file              string

	fileCtx  []byte
	keywords map[string]TokenType
	symbols  map[string]TokenType

	buffer    bytes.Buffer
	numErrors int
}

func (l *Lexer) current() (byte, bool) {
	if l.pos >= len(l.fileCtx) {
		return 0, false
	}

	return l.fileCtx[l.pos], true
}

func (l *Lexer) consume() byte {
	c := l.fileCtx[l.pos]
	l.pos++

	if c == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}

	return c
}

func (l *Lexer) unconsume() {
	if l.pos == 0 {
		return // Nothing to unconsume
	}

	l.pos--
	c := l.fileCtx[l.pos]

	if c == '\n' {
		l.line--
		// Recalculate l.column by finding the length of the previous line
		l.column = 1
		for i := l.pos - 1; i >= 0; i-- {
			if l.fileCtx[i] == '\n' {
				break
			}
			l.column++
		}
	} else {
		l.column--
	}
}

func newKeywords() map[string]TokenType {
	return map[string]TokenType{
		// Keywords
		"var":    TokenVar,
		"val":    TokenVal,
		"func":   TokenFunc,
		"mod":    TokenMod,
		"return": TokenReturn,
		"if":     TokenIf,
		"elsif":  TokenElsif,
		"else":   TokenElse,
		"while":  TokenWhile,
		"for":    TokenFor,
		"loop":   TokenLoop,
		"repeat": TokenRepeat,
		"until":  TokenUntil,
		"when":   TokenWhen,
		"struct": TokenStruct,
		"import": TokenImport,
		"do":     TokenDo,

		// Types
		"I64":  TokenI64,
		"I32":  TokenI32,
		"I16":  TokenI16,
		"I8":   TokenI8,
		"U64":  TokenU64,
		"U32":  TokenU32,
		"U16":  TokenU16,
		"U8":   TokenU8,
		"F64":  TokenF64,
		"F32":  TokenF32,
		"Any":  TokenAny,
		"Str":  TokenStr,
		"Bool": TokenBool,
		"Arr":  TokenArr,
	}
}

func newSymbols() map[string]TokenType {
	return map[string]TokenType{
		"+":  TokenPlus,
		"++": TokenPlusPlus,
		"+=": TokenPlusEq,
		"-":  TokenMinus,
		"--": TokenMinusMinus,
		"-=": TokenMinusEq,
		"*":  TokenMul,
		"*=": TokenMulEq,
		"/":  TokenDiv,
		"/=": TokenDivEq,
		"=":  TokenEq,
		"==": TokenEqEq,
		"!=": TokenNotEq,
		"!":  TokenNot,
		"&":  TokenBand,
		"|":  TokenBor,
		"^":  TokenBxor,
		"~":  TokenBnot,
		"&=": TokenBandEq,
		"|=": TokenBorEq,
		"^=": TokenBxorEq,
		"->": TokenArrow,
		",":  TokenComma,
		";":  TokenSemi,
		":":  TokenColon,
		"(":  TokenOpenParen,
		")":  TokenClosedParen,
		"{":  TokenOpenBrack,
		"}":  TokenClosedBrack,
		"[":  TokenOpenSqBrack,
		"]":  TokenClosedSqBrack,
		".":  TokenDot,
		"<":  TokenLessThan,
		">":  TokenGreaterThan,
		"<=": TokenLessEq,
		">=": TokenGreaterEq,
	}
}

func NewLexer(file string) (*Lexer, error) {
	ctx, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return &Lexer{
		line:    1,
		column:  1,
		file:    file,
		fileCtx: ctx,
		buffer:  bytes.Buffer{},

		keywords: newKeywords(),
		symbols:  newSymbols(),
	}, nil
}

func NewLexerFromStr(src string, fName string) Lexer {
	ctx := []byte(src)

	return Lexer{
		line:    1,
		column:  1,
		file:    fName,
		fileCtx: ctx,
		buffer:  bytes.Buffer{},

		keywords: newKeywords(),
		symbols:  newSymbols(),
	}
}

func (l *Lexer) GetErrors() int {
	return l.numErrors
}

func (l *Lexer) Scan() []Token {
	var toks []Token
	for {
		c, ok := l.current()
		if !ok { // We hit EOF
			break
		}

		if unicode.IsLetter(rune(c)) {
			startL, startC := l.line, l.column // Save the position

			// If it's a letter then consume it into a buffer
			for {
				c, ok := l.current()

				// A literal may contain underscores and numbers, so check for those too
				if !ok || (!unicode.IsLetter(rune(c)) && c != '_' && !unicode.IsDigit(rune(c))) {
					break
				}

				l.buffer.WriteByte(c)
				l.consume()
			}

			// Try find the literal in keywords
			k, ok := l.keywords[l.buffer.String()]
			if !ok {
				// It's a literal
				tok := Token{
					Line:   startL,
					Column: startC,
					Type:   TokenLiteral,

					Value: l.buffer.String(),
				}

				toks = append(toks, tok)
				l.buffer.Reset()
			} else {
				// It's a keyword
				tok := Token{
					Line:   startL,
					Column: startC,
					Type:   k,
				}

				toks = append(toks, tok)
				l.buffer.Reset()
			}
		} else if unicode.IsSpace(rune(c)) { // Covers \n, \t, spaces, etc
			if c == '\n' {
				tok := Token{
					Line:   l.line,
					Column: l.column,
					Type:   TokenEol,
				}

				toks = append(toks, tok)
			}
			l.consume()
		} else if c == '#' {
			for {
				c, ok := l.current()

				if !ok || c == '\n' {
					break
				}

				l.consume()
			}
		} else if unicode.IsDigit(rune(c)) { // Handle numbers and floats
			startL, startC := l.line, l.column // Save the position

			for {
				c, ok := l.current()

				if !ok || (!unicode.IsDigit(rune(c)) && c != '.') {
					break
				}

				l.buffer.WriteByte(c)
				l.consume()
			}

			if strings.Contains(l.buffer.String(), ".") {
				// It's a float
				f, err := strconv.ParseFloat(l.buffer.String(), 64)

				if err != nil {
					l.numErrors++
					fmt.Printf("kscript: %s: %s:%d:%d: %v\n", color.RedString("error"), l.file, startL, startC, err)

					l.buffer.Reset()
				} else {
					tok := Token{
						Line:   startL,
						Column: startC,
						Type:   TokenFloatLiteral,
						Value:  fmt.Sprintf("%f", f),
					}

					toks = append(toks, tok)
					l.buffer.Reset()
				}
			} else {
				// It's a number
				i, err := strconv.ParseInt(l.buffer.String(), 10, 64)

				if err != nil {
					l.numErrors++
					fmt.Printf("kscript: error: %s:%d:%d: %v\n", l.file, startL, startC, err)

					l.buffer.Reset()
				} else {
					tok := Token{
						Line:   startL,
						Column: startC,
						Type:   TokenNumLiteral,
						Value:  fmt.Sprintf("%d", i),
					}

					toks = append(toks, tok)
					l.buffer.Reset()
				}
			}
		} else { // Handle symbols
			startL, startC := l.line, l.column

			c1, ok := l.current()
			if !ok {
				break // EOF
			}

			if l.pos+1 < len(l.fileCtx) {
				twoChar := string(l.fileCtx[l.pos : l.pos+2])
				if sym, ok := l.symbols[twoChar]; ok {
					// Consume the 2 chars
					l.consume()
					l.consume()

					tok := Token{
						Line:   startL,
						Column: startC,
						Type:   sym,
					}

					toks = append(toks, tok)
					continue // Jump back to the loop
				}
			}

			oneChar := string(c1)
			if sym, ok := l.symbols[oneChar]; ok {
				l.consume()

				tok := Token{
					Line:   startL,
					Column: startC,

					Type: sym,
				}

				toks = append(toks, tok)
			} else if c == '"' {
				// Handle strings
				startL, startC := l.line, l.column
				l.consume()

				for {
					c, ok := l.current()
					if !ok || c == '"' {
						break
					}

					l.buffer.WriteByte(c)
					l.consume()
				}

				c, ok := l.current()
				if !ok || c != '"' {
					l.numErrors++
					if c != 0 {
						fmt.Printf("kscript: %s: %s:%d:%d: expected \", got %s\n", color.RedString("error"), l.file, startL, startC, string(c))
					} else {
						fmt.Printf("kscript: %s: %s:%d:%d: expected \", got <eof>\n", color.RedString("error"), l.file, startL, startC)
					}

					continue
				}
				l.consume()

				tok := Token{
					Line:   startL,
					Column: startC,
					Type:   TokenStrLiteral,

					Value: l.buffer.String(),
				}

				toks = append(toks, tok)
				l.buffer.Reset()
			} else {
				l.numErrors++
				fmt.Printf("kscript: %s: %s:%d:%d: unknown symbol %q\n", color.RedString("error"), l.file, startL, startC, oneChar)

				l.consume()
			}
		}
	}

	toks = append(toks, Token{
		Line:   l.line,
		Column: l.column,
		Type:   TokenEof,
	})
	return toks
}
