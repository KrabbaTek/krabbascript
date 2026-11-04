package parser

import (
	"fmt"
	"strings"

	"github.com/KrabbaTek/krabbascript/internal/lexer"

	"slices"
)

func (p *Parser) currentToks(s *TokenStream) lexer.Token {
	if s.pos >= len(s.toks) {
		if len(s.toks) == 0 {
			return lexer.Token{
				Line:   1,
				Column: 1,
				Type:   lexer.TokenEof,
			}
		}

		return s.toks[len(s.toks)-1]
	}

	return s.toks[s.pos]
}

func (p *Parser) peekToks(s *TokenStream) lexer.Token {
	if s.pos+1 >= len(s.toks) {
		if len(s.toks) == 0 {
			return lexer.Token{
				Line:   1,
				Column: 1,
				Type:   lexer.TokenEof,
			}
		}

		return s.toks[len(s.toks)-1]
	}

	return s.toks[s.pos+1]
}

func (p *Parser) consumeToks(s *TokenStream) lexer.Token {
	if s.pos >= len(s.toks) {
		if len(s.toks) == 0 {
			return lexer.Token{
				Line:   1,
				Column: 1,
				Type:   lexer.TokenEof,
			}
		}

		return s.toks[len(s.toks)-1]
	}

	t := s.toks[s.pos]
	s.pos++

	return t
}

func (p *Parser) expectToks(s *TokenStream, types ...lexer.TokenType) *ParserError {
	t := p.currentToks(s)

	for _, tt := range types {
		if t.Type == tt {
			return nil
		}
	}

	expected := make([]string, len(types))
	for i, tt := range types {
		expected[i] = fmt.Sprintf("'%s'", tt)
	}

	// missing semicolon
	if slices.Contains(types, lexer.TokenSemi) {
        return &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected a semicolon but either found none or found an arbitrary value instead.",
				p.file, t.Line, t.Column),
			Fix: fmt.Sprintf(`
			You forgot to terminate line %d with a semicolon — add a semicolon at the end of the line.
			If you want to cut the crap then just use my beautiful ASI
			—— through the format command or --format flag with the build command ——
			because it's not the 70s anymore and semicolons are stupid.
			`, t.Line),
		}
    }

	// missing type for an argument in a variable declaration
	// for example:
	// func add(a, b)
	// instead of:
	// func add(a: I32, b: I32)

	// if I was an amazing and prudent person then I'd also add an && check to make sure
	// the given token is a comma — to make damn sure (reference!) that we're dealing with a missing type

	// but that the only edge cases are so fucking arbitrary and esoteric it's not even worth it —
	// you'd have to go on a tangent inside a function declaration which is just borderline schizophrenic.

	// if you want to be a jerkoff and write stupid stuff in a function declaration then don't get mad when we throw the wrong error at you.
	if slices.Contains(types, lexer.TokenColon) {
		return &ParserError{
			What: fmt.Errorf("%s:%d:%d: in a function declaration on line %d, you forgot to declare a type for one of your arguments.",
				p.file, t.Line, t.Column, t.Line),
		}
	}

	return &ParserError{
		What: fmt.Errorf("%s:%d:%d: expected %s, got '%s'",
			p.file, t.Line, t.Column, strings.Join(expected, ", "), t.Type),
	}
}

func (p *Parser) expectLitToks(s *TokenStream) *ParserError {
	return p.expectToks(s, lexer.TokenLiteral, lexer.TokenFloatLiteral, lexer.TokenStrLiteral, lexer.TokenNumLiteral)
}

func (p *Parser) expectLitToksExtra(s *TokenStream, types ...lexer.TokenType) *ParserError {
	base := []lexer.TokenType{
		lexer.TokenLiteral,
		lexer.TokenFloatLiteral,
		lexer.TokenStrLiteral,
		lexer.TokenNumLiteral,
	}

	toks := append(base, types...)

	return p.expectToks(s, toks...)
}

func (p *Parser) skipToks(s *TokenStream, types ...lexer.TokenType) *ParserError {
	err := p.expectToks(s, types...)
	if err != nil {
		return err
	}

	p.consumeToks(s)
	return nil
}

func (p *Parser) isCurrentBuiltInTypeToks(s *TokenStream) bool {
	switch p.currentToks(s).Type {
	case lexer.TokenBool, lexer.TokenStr,
		lexer.TokenI64, lexer.TokenI32, lexer.TokenI16, lexer.TokenI8,
		lexer.TokenU64, lexer.TokenU32, lexer.TokenU16, lexer.TokenU8,
		lexer.TokenAny, lexer.TokenArr:
		return true
	default:
		return false
	}
}

func (p *Parser) checkMisplacedListInit(s *TokenStream) bool {
	depth := 0
	for {
		t := p.currentToks(s)
		switch t.Type {
		case lexer.TokenEof:
			return false
		case lexer.TokenOpenBrack:
			depth++
		case lexer.TokenClosedBrack:
			depth--
			if depth == 0 {
				return p.peekToks(s).Type == lexer.TokenOpenBrack
			}
		}

		p.consumeToks(s)
	}
}
