package parser

import (
	"fmt"
	"kscript/internal/lexer"

	"github.com/fatih/color"
)

type TokenStream struct {
	toks []lexer.Token
	pos  int
}

type Parser struct {
	line, column int
	file         string

	toks TokenStream

	bindings map[lexer.TokenType]Binding
	errors   []ParserError
}

func newBindings() map[lexer.TokenType]Binding {
	return map[lexer.TokenType]Binding{
		// Addition and subtraction (+,-)
		lexer.TokenPlus:  {BindingPrecedence: PrecedenceTerm, BindingSide: LeftSide},
		lexer.TokenMinus: {BindingPrecedence: PrecedenceTerm, BindingSide: LeftSide},
		// Multiplication and division (*,/)
		lexer.TokenMul: {BindingPrecedence: PrecedenceFactor, BindingSide: LeftSide},
		lexer.TokenDiv: {BindingPrecedence: PrecedenceFactor, BindingSide: LeftSide},
		// Compasion (==, !=, <, >, >=, <=)
		lexer.TokenEqEq:        {BindingPrecedence: PrecedenceComparison, BindingSide: LeftSide},
		lexer.TokenNotEq:       {BindingPrecedence: PrecedenceComparison, BindingSide: LeftSide},
		lexer.TokenLessThan:    {BindingPrecedence: PrecedenceComparison, BindingSide: LeftSide},
		lexer.TokenGreaterThan: {BindingPrecedence: PrecedenceComparison, BindingSide: LeftSide},
		lexer.TokenLessEq:      {BindingPrecedence: PrecedenceComparison, BindingSide: LeftSide},
		lexer.TokenGreaterEq:   {BindingPrecedence: PrecedenceComparison, BindingSide: LeftSide},
		// Bitwise (&,|,^)
		lexer.TokenBand: {BindingPrecedence: PrecedenceBitwise, BindingSide: LeftSide},
		lexer.TokenBor:  {BindingPrecedence: PrecedenceBitwise, BindingSide: LeftSide},
		lexer.TokenBxor: {BindingPrecedence: PrecedenceBitwise, BindingSide: LeftSide},
		// Assignment (=,+=,-=,/=,*=,&=,|=,^=)
		lexer.TokenEq:      {BindingPrecedence: PrecedenceAssignment, BindingSide: RightSide},
		lexer.TokenPlusEq:  {BindingPrecedence: PrecedenceAssignment, BindingSide: RightSide},
		lexer.TokenMinusEq: {BindingPrecedence: PrecedenceAssignment, BindingSide: RightSide},
		lexer.TokenDivEq:   {BindingPrecedence: PrecedenceAssignment, BindingSide: RightSide},
		lexer.TokenMulEq:   {BindingPrecedence: PrecedenceAssignment, BindingSide: RightSide},
		lexer.TokenBandEq:  {BindingPrecedence: PrecedenceAssignment, BindingSide: RightSide},
		lexer.TokenBorEq:   {BindingPrecedence: PrecedenceAssignment, BindingSide: RightSide},
		lexer.TokenBxorEq:  {BindingPrecedence: PrecedenceAssignment, BindingSide: RightSide},
		// Postfix ((), [],.)
		lexer.TokenOpenParen:   {BindingPrecedence: PrecedencePostfix, BindingSide: LeftSide},
		lexer.TokenOpenSqBrack: {BindingPrecedence: PrecedencePostfix, BindingSide: LeftSide},
		lexer.TokenDot:         {BindingPrecedence: PrecedencePostfix, BindingSide: LeftSide},
	}
}

func NewParser(toks []lexer.Token, file string) Parser {
	return Parser{
		line:   1,
		column: 1,
		file:   file,

		toks: TokenStream{
			pos:  0,
			toks: toks,
		},
		bindings: newBindings(),
		errors:   make([]ParserError, 0, 0),
	}
}

func (p *Parser) parseListInit(s *TokenStream) (*Node, *ParserError) {
	start := p.currentToks(s)

	var buffer []lexer.Token
	depth := 0

	for {
		t := p.currentToks(s)

		if t.Type == lexer.TokenEof {
			return nil, &ParserError{
				What: fmt.Errorf(
					"%s:%d:%d: unclosed bracket",
					p.file, start.Line, start.Column,
				),
			}
		}

		switch t.Type {
		case lexer.TokenOpenBrack:
			depth++
		case lexer.TokenClosedBrack:
			depth--
		}

		buffer = append(buffer, t)
		p.consumeToks(s)

		if depth == 0 {
			break
		}
	}

	buffer = append(buffer, lexer.Token{
		Line:   start.Line,
		Column: start.Column,
		Type:   lexer.TokenEof,
	})

	stream := TokenStream{
		toks: buffer,
		pos:  0,
	}

	trimmed, err := p.trimListInit(&stream)
	if err != nil {
		return nil, err
	}

	left := &Node{
		line:      start.Line,
		column:    start.Column,
		Type:      NodeListInit,
		ExtraInfo: &NodeInfoBlock{},
	}

	for _, stream := range trimmed {
		node, err := p.parseExpressionWithMinBpToks(&stream, 0, true)
		if err != nil {
			return nil, err
		}

		p.appendToBlock(left, node)
	}

	return left, nil
}

func (p *Parser) parseExpressionWithMinBpToks(s *TokenStream, minBp int, allowListInits bool) (*Node, *ParserError) {

	err := p.expectLitToksExtra(
		s,
		lexer.TokenOpenParen,
		lexer.TokenOpenBrack,

		lexer.TokenBool,
		lexer.TokenStr,

		lexer.TokenI64,
		lexer.TokenI32,
		lexer.TokenI16,
		lexer.TokenI8,

		lexer.TokenU64,
		lexer.TokenU32,
		lexer.TokenU16,
		lexer.TokenU8,

		lexer.TokenAny,
		lexer.TokenArr,
	)
	var left *Node
	if err != nil {
		c := p.currentToks(s)
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected an expression, got '%s'", p.file, c.Line, c.Column, c),
		}
	}

	if p.currentToks(s).Type == lexer.TokenOpenParen {
		p.consumeToks(s)
		l, err := p.parseExpressionWithMinBpToks(s, 0, true)
		if err != nil {
			return nil, err
		}

		left = l

		err = p.skipToks(s, lexer.TokenClosedParen)
		if err != nil {
			return nil, err
		}

	} else if p.currentToks(s).Type == lexer.TokenOpenBrack {
		if !allowListInits {
			c := p.currentToks(s)
			return nil, &ParserError{
				What: fmt.Errorf("%s:%d:%d: making a list initializer here is not allowed", p.file, c.Line, c.Column),
				Fix:  "wrap the list init in parentheses",
			}
		}

		left, err = p.parseListInit(s)
		if err != nil {
			return nil, err
		}
	} else if p.isCurrentBuiltInTypeToks(s) ||
		(p.currentToks(s).Type == lexer.TokenLiteral &&
			p.peekToks(s).Type == lexer.TokenOpenBrack) {

		nType, err := p.convertType(s)
		if err != nil {
			return nil, &ParserError{What: err}
		}

		if c := p.currentToks(s); c.Type == lexer.TokenOpenBrack {
			if !allowListInits {
				return nil, &ParserError{
					What: fmt.Errorf("%s:%d:%d: initializer is not allowed here", p.file, c.Line, c.Column),
					Fix:  "wrap the initializer in parentheses",
				}
			}

			var lerr *ParserError
			left, lerr = p.parseListInit(s)
			if lerr != nil {
				return nil, lerr
			}
			left.Left = nType
		} else {
			left = nType
		}
	} else if !allowListInits &&
		p.currentToks(s).Type == lexer.TokenLiteral &&
		p.peekToks(s).Type == lexer.TokenOpenBrack &&
		p.checkMisplacedListInit(&TokenStream{
			toks: s.toks,
			pos:  s.pos,
		}) {

		c := p.currentToks(s)
		return nil, &ParserError{
			What: fmt.Errorf(
				"%s:%d:%d: making a list initializer here is not allowed\n\t%s: wrap the list init in parentheses",
				p.file, c.Line, c.Column, color.BlueString("fix"),
			),
		}

	} else {
		t := p.consumeToks(s)

		var nodeType NodeType

		switch t.Type {
		case lexer.TokenNumLiteral:
			nodeType = NodeNumLit

		case lexer.TokenFloatLiteral:
			nodeType = NodeFloatLit

		case lexer.TokenStrLiteral:
			nodeType = NodeStrLit

		case lexer.TokenLiteral:
			nodeType = NodeLit

		default:
			return nil, &ParserError{
				What: fmt.Errorf(
					"%s:%d:%d: unexpected literal %q",
					p.file, t.Line, t.Column, t.Value,
				),
			}
		}

		left = &Node{
			line:   t.Line,
			column: t.Column,
			Type:   nodeType,
			Lexeme: t.Value,
		}
	}

	for {
		op := p.currentToks(s)
		if op.Type == lexer.TokenEof {
			break
		}

		bp, ok := p.bindings[op.Type]
		if !ok {
			break
		}

		if int(bp.BindingPrecedence) <= minBp {
			break
		}

		p.consumeToks(s)
		nodeOp := &Node{
			line:   op.Line,
			column: op.Column,
			Type:   NodeBinOp,

			Lexeme: op.Type.String(),
		}

		var right *Node
		var err *ParserError

		if bp.BindingSide == RightSide {
			right, err = p.parseExpressionWithMinBpToks(s, int(bp.BindingPrecedence-1), allowListInits)
			if err != nil {
				return nil, err
			}
		} else {
			right, err = p.parseExpressionWithMinBpToks(s, int(bp.BindingPrecedence), allowListInits)
			if err != nil {
				return nil, err
			}
		}

		nodeOp.Left = left
		nodeOp.Right = right

		left = nodeOp
	}

	return left, nil
}

func (p *Parser) convertType(s *TokenStream) (*Node, error) {
	typ := p.currentToks(s)

	n := &Node{
		line:   typ.Line,
		column: typ.Column,
	}

	// Assign a type
	switch typ.Type {
	case lexer.TokenI64:
		n.Type = NodeI64Type
		p.consumeToks(s)
	case lexer.TokenI32:
		n.Type = NodeI32Type
		p.consumeToks(s)
	case lexer.TokenI16:
		n.Type = NodeI16Type
		p.consumeToks(s)
	case lexer.TokenI8:
		n.Type = NodeI8Type
		p.consumeToks(s)

	case lexer.TokenU64:
		n.Type = NodeU64Type
		p.consumeToks(s)
	case lexer.TokenU32:
		n.Type = NodeU32Type
		p.consumeToks(s)
	case lexer.TokenU16:
		n.Type = NodeU16Type
		p.consumeToks(s)
	case lexer.TokenU8:
		n.Type = NodeU8Type
		p.consumeToks(s)
	case lexer.TokenF64:
		n.Type = NodeF64Type
		p.consumeToks(s)
	case lexer.TokenF32:
		n.Type = NodeF32Type
		p.consumeToks(s)

	case lexer.TokenStr:
		n.Type = NodeStrType
		p.consumeToks(s)
	case lexer.TokenAny:
		n.Type = NodeAnyType
		p.consumeToks(s)
	case lexer.TokenBool:
		n.Type = NodeBoolType
		p.consumeToks(s)

	case lexer.TokenLiteral:
		n.Type = NodeLit
		n.Lexeme = typ.Value
		p.consumeToks(s)
	case lexer.TokenFloatLiteral:
		n.Type = NodeFloatLit
		n.Lexeme = typ.Value
		p.consumeToks(s)
	case lexer.TokenNumLiteral:
		n.Type = NodeNumLit
		n.Lexeme = typ.Value
		p.consumeToks(s)
	case lexer.TokenStrLiteral:
		n.Type = NodeStrLit
		n.Lexeme = typ.Value
		p.consumeToks(s)
	case lexer.TokenArr:
		n.Type = NodeArray
		p.consumeToks(s)

		err := p.skipToks(s, lexer.TokenLessThan)
		if err != nil {
			return nil, err.What
		}

		elementType, cerr := p.convertType(s)
		if cerr != nil {
			return nil, cerr
		}

		err = p.skipToks(s, lexer.TokenGreaterThan)
		if err != nil {
			return nil, err.What
		}

		n.Left = elementType
	default:
		return nil, fmt.Errorf("%s:%d:%d: expected a type, got %s", p.file, typ.Line, typ.Column, typ.Type)
	}

	return n, nil
}

func (p *Parser) parseVarToks(s *TokenStream) (*Node, *ParserError) {
	var node *Node

	start := p.consumeToks(s) // Will use this for the line, col fields
	err := p.expectToks(s, lexer.TokenLiteral)
	if err != nil {
		c := p.currentToks(s)
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected a name for variable (example: var %s), got '%s'", p.file, c.Line, c.Column, color.GreenString("butterfly"), c),
		}
	}

	name := p.consumeToks(s) // Save this for later

	// Check if it's = or :
	// Since it can be:
	// var krabba = 12;
	// var krabba: i32 = 12;
	err = p.expectToks(s, lexer.TokenEq, lexer.TokenColon)
	if err != nil {
		c := p.currentToks(s)
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected an assignment or ':' followed by a type, got '%s'", p.file, c.Line, c.Column, c),
		}
	}

	colEq := p.consumeToks(s)
	if colEq.Type == lexer.TokenColon {
		nType, err := p.convertType(s)
		if err != nil {
			return nil, &ParserError{
				What: err,
			}
		}

		eerr := p.expectToks(s, lexer.TokenEq, lexer.TokenSemi)
		if eerr != nil {
			c := p.currentToks(s)
			return nil, &ParserError{
				What: fmt.Errorf("%s:%d:%d: expected an assignment or ';', got '%s'", p.file, c.Line, c.Column, c),
			}
		}

		// Something like var krabba: i32;
		if p.currentToks(s).Type == lexer.TokenSemi {
			p.consumeToks(s) // Skip the semicolon

			node = &Node{
				line:   start.Line,
				column: start.Column,
				Type:   NodeVariableDec,
				Lexeme: name.Value,
				Left:   nType,
			}
			goto exit

		} else if p.currentToks(s).Type == lexer.TokenEq {
			p.consumeToks(s) // Skip the =

			expr, err := p.parseExpressionWithMinBpToks(s, 0, true)
			if err != nil {
				return nil, err
			}

			err = p.skipToks(s, lexer.TokenSemi)
			if err != nil {
				return nil, err
			}

			node = &Node{
				line:   start.Line,
				column: start.Column,
				Type:   NodeVariableDef,
				Lexeme: name.Value,
				Left:   nType,
				Right:  expr,
			}
			goto exit
		}
	} else { // Equals
		expr, err := p.parseExpressionWithMinBpToks(s, 0, true)
		if err != nil {
			return nil, err
		}

		err = p.skipToks(s, lexer.TokenSemi)
		if err != nil {
			return nil, err
		}

		node = &Node{
			line:   start.Line,
			column: start.Column,
			Type:   NodeVariableDef,
			Lexeme: name.Value,
			Right:  expr,
		}
		goto exit
	}

exit:
	return node, nil
}

func (p *Parser) parseValToks(s *TokenStream) (*Node, *ParserError) {
	var node *Node

	start := p.consumeToks(s) // Will use this for the line, col fields
	err := p.expectToks(s, lexer.TokenLiteral)
	if err != nil {
		c := p.currentToks(s)
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected a name for value (example: val %s), got '%s'", p.file, c.Line, c.Column, color.GreenString("butterfly"), c),
		}
	}

	name := p.consumeToks(s) // Save this for later

	// Check if it's = or :
	// Since it can be:
	// val krabba = 12;
	// val krabba: i32 = 12;
	err = p.expectToks(s, lexer.TokenEq, lexer.TokenColon)
	if err != nil {
		c := p.currentToks(s)
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected an assignment or ':' followed by a type, got '%s'", p.file, c.Line, c.Column, c),
		}
	}

	colEq := p.consumeToks(s)
	if colEq.Type == lexer.TokenColon {
		nType, err := p.convertType(s)
		if err != nil {
			return nil, &ParserError{
				What: err,
			}
		}

		eerr := p.expectToks(s, lexer.TokenEq)
		if eerr != nil {
			return nil, eerr
		}

		p.consumeToks(s) // Skip the =

		expr, eerr := p.parseExpressionWithMinBpToks(s, 0, true)
		if eerr != nil {
			return nil, eerr
		}

		eerr = p.skipToks(s, lexer.TokenSemi)
		if eerr != nil {
			return nil, eerr
		}

		node = &Node{
			line:   start.Line,
			column: start.Column,
			Type:   NodeValueDef,
			Lexeme: name.Value,
			Left:   nType,
			Right:  expr,
		}
		goto exit
	} else { // Equals
		expr, err := p.parseExpressionWithMinBpToks(s, 0, true)
		if err != nil {
			return nil, err
		}

		err = p.skipToks(s, lexer.TokenSemi)
		if err != nil {
			return nil, err
		}

		node = &Node{
			line:   start.Line,
			column: start.Column,
			Type:   NodeValueDef,
			Lexeme: name.Value,
			Right:  expr,
		}
		goto exit
	}

exit:
	return node, nil
}

func (p *Parser) parseScopeToks(s *TokenStream) (*Node, *ParserError) {
	err := p.skipToks(s, lexer.TokenOpenBrack)
	if err != nil {
		return nil, err
	}

	var toks []lexer.Token
	depth := 0
	for {
		t := p.currentToks(s)
		if t.Type == lexer.TokenEof {
			break
		}

		if t.Type == lexer.TokenOpenBrack {
			depth++
		} else if t.Type == lexer.TokenClosedBrack {
			if depth == 0 {
				break
			}
			depth--
		}

		toks = append(toks, t)
		p.consumeToks(s)
	}

	err = p.skipToks(s, lexer.TokenClosedBrack)
	if err != nil {
		return nil, err
	}

	// Nothing to parse...
	if len(toks) == 0 {
		return nil, nil
	}

	last := toks[len(toks)-1]
	toks = append(toks, lexer.Token{
		Line:   last.Line,
		Column: last.Column,
		Type:   lexer.TokenEof,
	})

	// The errors will get pushed to errors field
	n, _ := p.ParseToks(&TokenStream{toks: toks, pos: 0})
	n.Type = NodeScope

	return n, nil
}

func (p *Parser) parseStructFieldToks(s *TokenStream) (*Node, *ParserError) {
	err := p.expectToks(s, lexer.TokenLiteral)
	if err != nil {
		return nil, err
	}

	name := p.consumeToks(s)

	err = p.skipToks(s, lexer.TokenColon)
	if err != nil {
		return nil, err
	}

	typ, cerr := p.convertType(s)
	if cerr != nil {
		return nil, &ParserError{
			What: cerr,
		}
	}

	err = p.skipToks(s, lexer.TokenSemi)
	if err != nil {
		return nil, err
	}

	n := Node{
		line:   name.Line,
		column: name.Column,

		Type:   NodeFieldDec,
		Lexeme: name.Value,

		Left: typ,
	}

	return &n, nil
}

func (p *Parser) parseStructFieldsToks(s *TokenStream) (*Node, *ParserError) {
	node := &Node{
		ExtraInfo: &NodeInfoBlock{},
		Type:      NodeStructDec,
	}
	for {
		t := p.currentToks(s)
		if t.Type == lexer.TokenEof {
			return node, nil
		}

		n, err := p.parseStructFieldToks(s)
		if err != nil {
			return nil, err
		}

		p.appendToBlock(node, n)
	}
}

func (p *Parser) parseStructScopeToks(s *TokenStream) (*Node, *ParserError) {
	err := p.skipToks(s, lexer.TokenOpenBrack)
	if err != nil {
		return nil, err
	}

	var toks []lexer.Token
	for {
		t := p.currentToks(s)
		if t.Type == lexer.TokenEof || t.Type == lexer.TokenClosedBrack {
			break
		}

		toks = append(toks, t)
		p.consumeToks(s)
	}

	err = p.skipToks(s, lexer.TokenClosedBrack)
	if err != nil {
		return nil, err
	}

	for _, t := range toks {
		fmt.Println(t)
	}

	if len(toks) != 0 {

		last := toks[len(toks)-1]
		toks = append(toks, lexer.Token{
			Line:   last.Line,
			Column: last.Column,
			Type:   lexer.TokenEof,
		})
	}

	stream := TokenStream{
		toks: toks,
		pos:  0,
	}
	return p.parseStructFieldsToks(&stream)
}

func (p *Parser) parseStructToks(s *TokenStream) (*Node, *ParserError) {
	start := p.consumeToks(s) // Save this for line, col fields

	err := p.expectToks(s, lexer.TokenLiteral)
	if err != nil {
		return nil, err
	}

	name := p.consumeToks(s) // Save the name too

	n, err := p.parseStructScopeToks(s)
	if err != nil {
		return nil, err
	}

	if n == nil {
		return nil, nil
	}

	n.line = start.Line
	n.column = start.Column

	n.Lexeme = name.Value

	return n, nil
}

/*
Some code for trimming { 1, 2 }, 3 for example

Get current
Depth = 0

While current != EOF
    While Depth > 0 || current != , && current != EOF
        If current == {
            Depth++
        Else if current == }
            If Depth == 0 Error out
            Depth--

        Add to buffer
        Update current

    If Depth != 0 Error out

    Add buffer to the array of token streams
    Clear buffer

    If current == EOF break
    Else skip ,

Return trimmed
*/

func (p *Parser) trimListInit(s *TokenStream) ([]TokenStream, *ParserError) {
	var trimmed []TokenStream
	var buffer []lexer.Token
	depth := 0

	// Delete the { if present
	if p.currentToks(s).Type == lexer.TokenOpenBrack {
		p.consumeToks(s)
	}

	for {
		t := p.currentToks(s)

		if t.Type == lexer.TokenEof || (t.Type == lexer.TokenClosedBrack && depth == 0) {
			break
		}

		if t.Type == lexer.TokenComma && depth == 0 {
			p.consumeToks(s)

			start := p.currentToks(s) // Save this for error handling
			if start.Type == lexer.TokenClosedBrack {
				return nil, &ParserError{
					What: fmt.Errorf("%s:%d:%d: expected an element, got {", p.file, start.Line, start.Column),
				}
			}
			if len(buffer) > 0 {
				buffer = append(buffer, lexer.Token{
					Line:   t.Line,
					Column: t.Column,
					Type:   lexer.TokenEof,
				})

				trimmed = append(trimmed, TokenStream{
					toks: buffer,
					pos:  0,
				})

				buffer = nil
			}

			continue
		}
		switch t.Type {
		case lexer.TokenOpenBrack:
			depth++
		case lexer.TokenClosedBrack:
			depth--
			if depth < 0 {
				return nil, &ParserError{
					What: fmt.Errorf("%s:%d:%d: unexpected closing bracket", p.file, t.Line, t.Column),
				}
			}
		}

		buffer = append(buffer, t)
		p.consumeToks(s)
	}

	if depth != 0 {
		t := p.currentToks(s)
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: unclosed bracket", p.file, t.Line, t.Column),
		}
	}

	if len(buffer) > 0 {
		buffer = append(buffer, lexer.Token{
			Line:   1,
			Column: 1,
			Type:   lexer.TokenEof,
		})

		trimmed = append(trimmed, TokenStream{
			toks: buffer,
			pos:  0,
		})
	}

	// Same thing for }
	if p.currentToks(s).Type == lexer.TokenClosedBrack {
		p.consumeToks(s)
	}

	return trimmed, nil
}

func (p *Parser) parseIfToks(s *TokenStream) (*Node, *ParserError) {
	start := p.consumeToks(s) // Will use this for the line, col fields

	expr, err := p.parseExpressionWithMinBpToks(s, 0, false)
	if err != nil {
		return nil, err
	}

	body, err := p.parseScopeToks(s)
	if err != nil {
		return nil, err
	}

	node := &Node{
		line:      start.Line,
		column:    start.Column,
		Type:      NodeIfStatement,
		ExtraInfo: &NodeInfoBlock{},
		Left:      expr,
	}

	if body != nil {
		p.appendToBlock(node, body)
	}

	return node, nil
}

func (p *Parser) parseElsifToks(s *TokenStream) (*Node, *ParserError) {
	start := p.consumeToks(s)

	expr, err := p.parseExpressionWithMinBpToks(s, 0, false)
	if err != nil {
		return nil, err
	}

	body, err := p.parseScopeToks(s)
	if err != nil {
		return nil, err
	}

	node := &Node{
		line:      start.Line,
		column:    start.Column,
		Type:      NodeElsifStatement,
		ExtraInfo: &NodeInfoBlock{},
		Left:      expr,
	}

	if body != nil {
		p.appendToBlock(node, body)
	}

	return node, nil
}

func (p *Parser) parseElseToks(s *TokenStream) (*Node, *ParserError) {
	start := p.consumeToks(s)

	body, err := p.parseScopeToks(s)
	if err != nil {
		return nil, err
	}

	node := &Node{
		line:      start.Line,
		column:    start.Column,
		Type:      NodeElseStatement,
		ExtraInfo: &NodeInfoBlock{},
	}

	if body != nil {
		p.appendToBlock(node, body)
	}

	return node, nil
}

func (p *Parser) parseImportModulePathToks(s *TokenStream) (*Node, *ParserError) {
	// It's guaranteed that the stream we got is not EOF, but check just in case
	start := p.currentToks(s)

	if start.Type == lexer.TokenEof {
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected a module path (example: import %s;), got %s", p.file, start.Line, start.Column, color.GreenString("Std.SomeMod"), start),
		}
	}

	path := Node{
		line:   start.Line,
		column: start.Column,
		Type:   NodeModulePath,

		ExtraInfo: &NodeInfoBlock{}, // This will hold the module path
	}

	for {
		c := p.currentToks(s)

		if c.Type == lexer.TokenEof {
			// Reached end of stream
			break
		}

		err := p.expectToks(s, lexer.TokenLiteral)
		if err != nil {
			return nil, &ParserError{
				What: fmt.Errorf("%s:%d:%d: expected module name, got '%s'", p.file, c.Line, c.Column, c),
			}
		}

		next := p.peekToks(s)

		if next.Type != lexer.TokenDot && next.Type != lexer.TokenEof {
			return nil, &ParserError{
				What: fmt.Errorf("%s:%d:%d: expected '.' or ';' after module, got %s", p.file, next.Line, next.Column, next),
			}
		}

		// Add the literal to the NodeBlock
		lit := p.consumeToks(s)
		p.appendToBlock(&path, &Node{
			Lexeme: lit.Value,
			Type:   NodeModuleImport,
		})

		c = p.currentToks(s)

		// Skip the dot if exists
		if c.Type == lexer.TokenDot {
			p.consumeToks(s)

			// If there is a dot after EOF, throw a mean error
			c = p.currentToks(s)
			if c.Type == lexer.TokenEof {
				return nil, &ParserError{
					What: fmt.Errorf("%s:%d:%d: trailing comma. expected module name after '.', got '%s'", p.file, c.Line, c.Column, c),
				}
			}
		}
	}

	return &path, nil
}

func (p *Parser) parseImportToks(s *TokenStream) (*Node, *ParserError) {
	start := p.consumeToks(s)

	// Collect everything into a buffer which will hold the node path.

	// import Std.EpicMod;
	//        ^^^^^^^^^^^

	current := p.currentToks(s)
	if current.Type == lexer.TokenEof || current.Type == lexer.TokenSemi {
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected a module path (example: import %s;), got %s", p.file, current.Line, current.Column, color.GreenString("Std.SomeMod"), current),
		}
	}

	var toks []lexer.Token
	for {
		t := p.currentToks(s)
		if t.Type == lexer.TokenEof || t.Type == lexer.TokenSemi {
			break
		}

		toks = append(toks, t)
		p.consumeToks(s)
	}

	err := p.skipToks(s, lexer.TokenSemi)
	if err != nil {
		return nil, err
	}

	node := Node{
		line:   start.Line,
		column: start.Column,

		Type: NodeImportStatement,
	}

	// Add an EOF by getting the last element and using it as the line and column for the EOF token
	last := toks[len(toks)-1]
	toks = append(toks, lexer.Token{
		Line:   last.Line,
		Column: last.Column,
		Type:   lexer.TokenEof,
	})

	// Create a new stream for ImportModulePath to work with
	stream := TokenStream{
		toks: toks,
		pos:  0, // This is just boilerplate
	}

	// Parse the module path
	path, merror := p.parseImportModulePathToks(&stream)
	if merror != nil {
		return nil, merror
	}

	// Set the path to the right side
	node.Right = path
	return &node, nil
}

func (p *Parser) parseWhenArm(s *TokenStream) (*Node, *ParserError) {
	// The start of an arm in KrabbaScript is a comparison.
	// Example:
	//		when breakfast {
	//			== "eggs" do cook_eggs();
	//			^^ This
	//		}
	// The comparator can be instead an else, example:
	//		when breakfast {
	//			== "eggs" do cook_eggs();
	//			else do Std.print("I don't know what to cook!");
	//		}

	start := p.currentToks(s)
	arm := Node{
		line:   start.Line,
		column: start.Column,
		Type:   NodeWhenArm,
	}

	bp, ok := p.bindings[start.Type]
	if !ok && start.Type != lexer.TokenElse {
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected a comparison (like: %s or an else), got %s", p.file, start.Line, start.Column,
				color.GreenString("==, !="), start),
		}
	}
	if bp.BindingPrecedence != PrecedenceComparison && start.Type != lexer.TokenElse {
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected a comparison (like: %s), got %s", p.file, start.Line, start.Column,
				color.GreenString("== \"krabba\""), start),
		}
	}

	if start.Type == lexer.TokenElse {
		p.consumeToks(s)

		body, err := p.parseWhenDo(s, start)
		if err != nil {
			return nil, err
		}

		return &Node{
			line:   start.Line,
			column: start.Column,
			Type:   NodeWhenElseArm,
			Right:  body,
		}, nil
	}

	p.consumeToks(s)

	arm.Left = &Node{
		line:   start.Line,
		column: start.Column,
		Type:   NodeBinOp,
		Lexeme: start.Type.String(),
	}

	err := p.expectLitToks(s)
	if err != nil {
		c := p.currentToks(s)
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected a literal to compare against (like: %s), got %s", p.file, c.Line, c.Column,
				color.GreenString("== \"krabba\""), c),
		}
	}

	n, cerr := p.convertType(s)
	if cerr != nil {
		return nil, &ParserError{What: cerr}
	}
	arm.Left.Left = n

	body, perr := p.parseWhenDo(s, start)
	if perr != nil {
		return nil, perr
	}
	arm.Right = body

	return &arm, nil
}

func (p *Parser) parseWhenDo(s *TokenStream, armStart lexer.Token) (*Node, *ParserError) {
	// The do acts as a separator instead of an anonymous function. Skip it
	err := p.skipToks(s, lexer.TokenDo)
	if err != nil {
		return nil, err
	}

	// The arm body can be both a scope and a 1 liner. Example:
	// when breakfast {
	//		== "eggs" do cookEggs();
	//		== "coffee 'n sandwich" do {
	//			brew_coffee();
	//			make_sandwich();
	//		}
	// }
	if p.currentToks(s).Type == lexer.TokenOpenBrack {
		return p.parseScopeToks(s)
	}

	var toks []lexer.Token
	for {
		t := p.currentToks(s)
		if t.Type == lexer.TokenEof || t.Type == lexer.TokenEol {
			break
		}
		toks = append(toks, t)
		p.consumeToks(s)
	}

	if p.currentToks(s).Type == lexer.TokenEol {
		p.consumeToks(s)
	}

	eofLine, eofCol := armStart.Line, armStart.Column
	if len(toks) > 0 {
		last := toks[len(toks)-1]
		eofLine, eofCol = last.Line, last.Column
	}
	toks = append(toks, lexer.Token{Line: eofLine, Column: eofCol, Type: lexer.TokenEof})

	body, _ := p.ParseToks(&TokenStream{toks: toks, pos: 0})
	body.Type = NodeScope
	return body, nil
}

func (p *Parser) parseWhenArms(s *TokenStream) ([]*Node, *ParserError) {
	var arms []*Node
	for {
		// Remove all EOLs
		for p.currentToks(s).Type == lexer.TokenEol {
			p.consumeToks(s)
		}

		t := p.currentToks(s)
		if t.Type == lexer.TokenEof {
			return arms, nil
		}

		arm, err := p.parseWhenArm(s)
		if err != nil {
			return nil, err
		}

		if arm.Type == NodeWhenElseArm {
			for p.currentToks(s).Type == lexer.TokenEol {
				p.consumeToks(s)
			}
			if next := p.currentToks(s); next.Type != lexer.TokenEof {
				return nil, &ParserError{
					What: fmt.Errorf("%s:%d:%d: unreachable arm after 'else'", p.file, next.Line, next.Column),
					Fix:  "the else arm must be the last arm of a when",
				}
			}
			return append(arms, arm), nil
		}
		arms = append(arms, arm)
	}
}

func (p *Parser) parseWhenToks(s *TokenStream) (*Node, *ParserError) {
	start := p.consumeToks(s)

	node := Node{
		line:      start.Line,
		column:    start.Column,
		Type:      NodeWhenStatement,
		ExtraInfo: &NodeInfoBlock{},
	}

	err := p.expectLitToks(s)
	if err != nil {
		c := p.currentToks(s)
		return nil, &ParserError{
			What: fmt.Errorf("%s:%d:%d: expected a literal after 'when', got %s", p.file, c.Line, c.Column, c),
			Fix:  fmt.Sprintf("put a literal after 'when' (example: when %s { ... })", color.GreenString("butterfly_type")),
		}
	}

	lit := p.consumeToks(s)
	node.Right = &Node{
		line:   lit.Line,
		column: lit.Column,
		Type:   NodeLit,
		Lexeme: lit.Value,
	}

	err = p.expectToks(s, lexer.TokenOpenBrack)
	if err != nil {
		return nil, err
	}

	scopeStart := p.consumeToks(s)

	depth := 0
	var toks []lexer.Token

	for {
		t := p.currentToks(s)
		if t.Type == lexer.TokenEof {
			return nil, &ParserError{
				What: fmt.Errorf("%s:%d:%d: unclosed bracket", p.file, scopeStart.Line, scopeStart.Column),
			}
		}

		if t.Type == lexer.TokenOpenBrack {
			depth++
		} else if t.Type == lexer.TokenClosedBrack {
			if depth == 0 {
				break
			}
			depth--
		}

		toks = append(toks, t)
		p.consumeToks(s)
	}

	p.consumeToks(s) // Cnsume the closing brack

	// Append an EOF by getting the last element to use it as the line and column
	if len(toks) == 0 {
		// No EOF, use the startScope var
		toks = append(toks, lexer.Token{
			Line:   scopeStart.Line,
			Column: scopeStart.Column,
			Type:   lexer.TokenEof,
		})
	} else {
		last := toks[len(toks)-1]

		toks = append(toks, lexer.Token{
			Line:   last.Line,
			Column: last.Column,
			Type:   lexer.TokenEof,
		})
	}

	// Create a new stream with an EOF for the parseWhenArms
	stream := TokenStream{
		toks: toks,
		pos:  0, // Boilerplate
	}

	arms, err := p.parseWhenArms(&stream)
	if err != nil {
		return nil, err
	}

	// Add all of the arms to the when statement
	for _, arm := range arms {
		p.appendToBlock(&node, arm)
	}

	return &node, nil
}

func (p *Parser) appendToBlock(block *Node, node *Node) {
	switch v := block.ExtraInfo.(type) {
	case *NodeInfoBlock:
		v.Elements = append(v.Elements, node)
	}
}

// Will advance the stream until a safe place to continue parsing is reached
func (p *Parser) sync(s *TokenStream) {
	depth := 0

	for {
		t := p.currentToks(s)

		if t.Type == lexer.TokenEof {
			return
		}

		switch t.Type {
		case lexer.TokenOpenBrack:
			depth++
			p.consumeToks(s)

		case lexer.TokenClosedBrack:
			if depth > 0 {
				depth--
				p.consumeToks(s)
				continue
			}

			// Don't consume the enclosing scope's closing brace
			return

		case lexer.TokenSemi:
			if depth == 0 {
				p.consumeToks(s)
				return
			}
			p.consumeToks(s)

		default:
			// If we reach a new statement at the top level,
			// let the main parser handle it
			if depth == 0 {
				switch t.Type {
				case lexer.TokenVar, lexer.TokenVal,
					lexer.TokenIf, lexer.TokenElsif,
					lexer.TokenElse, lexer.TokenStruct,
					lexer.TokenImport, lexer.TokenWhen:
					return
				}
			}

			p.consumeToks(s)
		}
	}
}

func (p *Parser) ParseToks(s *TokenStream) (*Node, []ParserError) {
	ast := &Node{
		Type:      NodeRoot,
		ExtraInfo: &NodeInfoBlock{},
	}

loop:
	for {
		t := p.currentToks(s)

		switch t.Type {
		case lexer.TokenEof:
			break loop

		case lexer.TokenVar:
			n, err := p.parseVarToks(s)
			if err != nil {
				p.errors = append(p.errors, *err)
				p.sync(s)
			} else {
				p.appendToBlock(ast, n)
			}
		case lexer.TokenVal:
			n, err := p.parseValToks(s)
			if err != nil {
				p.errors = append(p.errors, *err)
				p.sync(s)
			} else {
				p.appendToBlock(ast, n)
			}
		case lexer.TokenIf:
			n, err := p.parseIfToks(s)
			if err != nil {
				p.errors = append(p.errors, *err)
				p.sync(s)
			} else {
				p.appendToBlock(ast, n)
			}
		case lexer.TokenElsif:
			n, err := p.parseElsifToks(s)
			if err != nil {
				p.errors = append(p.errors, *err)
				p.sync(s)
			} else {
				p.appendToBlock(ast, n)
			}
		case lexer.TokenElse:
			n, err := p.parseElseToks(s)
			if err != nil {
				p.errors = append(p.errors, *err)
				p.sync(s)
			} else {
				p.appendToBlock(ast, n)
			}
		case lexer.TokenStruct:
			n, err := p.parseStructToks(s)
			if err != nil {
				p.errors = append(p.errors, *err)
				p.sync(s)
			} else {
				p.appendToBlock(ast, n)
			}
		case lexer.TokenImport:
			n, err := p.parseImportToks(s)
			if err != nil {
				p.errors = append(p.errors, *err)
				p.sync(s)
			} else {
				p.appendToBlock(ast, n)
			}
		case lexer.TokenWhen:
			n, err := p.parseWhenToks(s)
			if err != nil {
				p.errors = append(p.errors, *err)
				p.sync(s)
			} else {
				p.appendToBlock(ast, n)
			}
		case lexer.TokenSemi:
			p.consumeToks(s)
		case lexer.TokenEol:
			p.consumeToks(s)
		default:
			p.errors = append(p.errors, ParserError{
				What: fmt.Errorf("%s:%d:%d: unexpected token %s", p.file, t.Line, t.Column, t),
			})

			p.sync(s)
		}
	}

	return ast, p.errors
}

func (p *Parser) Parse() (*Node, []ParserError) {
	return p.ParseToks(&p.toks)
}
