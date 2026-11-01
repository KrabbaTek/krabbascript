package typechecker

import (
	"fmt"
	"github.com/KrabbaTek/krabbascript/internal/parser"
)

type TypeChecker struct {
	file    string
	errors  []error
	scopes  []map[string]*Type
	returns []*Type
}

func New(file string) *TypeChecker {
	return &TypeChecker{file: file}
}
func (tc *TypeChecker) errorf(n *parser.Node, format string, args ...any) {
	line, col := 0, 0
	if n != nil {
		line, col = n.Line, n.Column
	}
	msg := fmt.Sprintf(format, args...)
	// errors must carry file:line:col
	tc.errors = append(tc.errors, fmt.Errorf("%s:%d:%d: %s", tc.file, line, col, msg))
}
func (tc *TypeChecker) push() {
	tc.scopes = append(tc.scopes, make(map[string]*Type))
}

func (tc *TypeChecker) pop() {
	tc.scopes = tc.scopes[:len(tc.scopes)-1]
}

func (tc *TypeChecker) declare(name string, t *Type) {
	if name == "" {
		return //anonymous nodes don't get bound
	}
	tc.scopes[len(tc.scopes)-1][name] = t
}

func (tc *TypeChecker) lookup(name string) (*Type, bool) {
	for i := len(tc.scopes) - 1; i >= 0; i-- {
		if t, ok := tc.scopes[i][name]; ok {
			return t, true
		}
	}
	return nil, false
}

// check walks the AST and collects every type error it finds.
func (tc *TypeChecker) Check(ast *parser.Node) []error {
	if ast == nil {
		return nil
	}
	block, ok := ast.ExtraInfo.(*parser.NodeInfoBlock)
	if !ok {
		return nil
	}
	tc.push()
	defer tc.pop()

	//register toplevel functions so they can be called
	// before their definition appears in the source. i have done
	// the same thing before so i will also work here(for now)
	for _, el := range block.Elements {
		if el == nil || el.Type != parser.NodeFunction {
			continue
		}
		if info, ok := el.ExtraInfo.(*parser.NodeInfoFunction); ok {
			tc.declare(el.Lexeme, funcType(info))
		}
	}

	for _, el := range block.Elements {
		tc.checkNode(el)
	}

	return tc.errors
}

func funcType(info *parser.NodeInfoFunction) *Type {
	ret := &Type{Kind: parser.NodeVoidType}
	if info.ReturnType != nil {
		ret = fromTypeNode(info.ReturnType)
	}
	ft := &Type{Kind: parser.NodeFunction, Elem: ret}
	for _, a := range info.Args {
		if a == nil {
			ft.Params = append(ft.Params, unk())
			continue
		}
		ft.Params = append(ft.Params, fromTypeNode(firstNode(a.Left, a.Right)))
	}
	return ft
}

func isExpression(t parser.NodeType) bool {
	switch t {
	case parser.NodeBinOp, parser.NodeUnary, parser.NodeNumLit, parser.NodeFloatLit,
		parser.NodeNullLit,
		parser.NodeLit, parser.NodeStrLit, parser.NodeArray,
		parser.NodePointer, parser.NodeIndex, parser.NodeFunctionCall,
		parser.NodeListInit, parser.NodeAnonymousFunction:
		return true
	}
	return false
}

// operand returns the single operand of a unary operator or return statement.
// right to left order
func operand(n *parser.Node) *parser.Node {
	return firstNode(n.Right, n.Left)
}

func (tc *TypeChecker) checkNode(n *parser.Node) {
	if n == nil {
		return
	}

	switch n.Type {
	case parser.NodeVariableDef, parser.NodeVariableDec, parser.NodeValueDef:
		tc.checkDecl(n)

	case parser.NodeFunction:
		if info, ok := n.ExtraInfo.(*parser.NodeInfoFunction); ok {
			ft := funcType(info)
			tc.declare(n.Lexeme, ft)
			tc.checkFunctionBody(info, ft)
		}

	case parser.NodeScope:
		tc.push()
		tc.checkChildren(n)
		tc.pop()

	case parser.NodeReturn:
		tc.checkReturn(n)

	default:
		if isExpression(n.Type) {
			tc.infer(n)
			return
		}
		// anything else in here.
		tc.checkChildren(n)
	}
}

func (tc *TypeChecker) checkChildren(n *parser.Node) {
	tc.checkNode(n.Left)
	tc.checkNode(n.Right)

	if block, ok := n.ExtraInfo.(*parser.NodeInfoBlock); ok {
		for _, el := range block.Elements {
			tc.checkNode(el)
		}
	}
}

func (tc *TypeChecker) checkFunctionBody(info *parser.NodeInfoFunction, ft *Type) {
	tc.push()

	for i, a := range info.Args {
		if a != nil {
			tc.declare(a.Lexeme, ft.Params[i])
		}
	}

	tc.returns = append(tc.returns, ft.Elem)
	tc.checkNode(info.Body)
	tc.returns = tc.returns[:len(tc.returns)-1]
	tc.pop()
}

func (tc *TypeChecker) checkReturn(n *parser.Node) {
	expr := operand(n)
	if expr == nil {
		return
	}

	got := tc.infer(expr)
	if len(tc.returns) == 0 {
		return
	}

	want := tc.returns[len(tc.returns)-1]
	switch {
	case want.Kind == parser.NodeVoidType:
		tc.errorf(n, "void function cannot return a value")
	case !assignable(want, got):
		tc.errorf(expr, "expected `%s`, got `%s`", want, got)
	}
}

// checkDecl handles var / val declarations.
// If a type annotation is present it must be compatible with the initializer or
// otherwise the type is inferred from the initializer.
func (tc *TypeChecker) checkDecl(decl *parser.Node) {
	var got *Type
	if decl.Right != nil {
		got = tc.infer(decl.Right)
	}

	switch {
	case decl.Left != nil:
		declared := fromTypeNode(decl.Left)
		if got != nil && !assignable(declared, got) {
			tc.errorf(decl.Left, "expected `%s`, got `%s`", declared, got)
		}
		tc.declare(decl.Lexeme, declared)

	case got != nil:
		tc.declare(decl.Lexeme, tc.inferred(decl, got))

	default:
		tc.declare(decl.Lexeme, unk())
	}
}

func (tc *TypeChecker) inferred(decl *parser.Node, got *Type) *Type {
	switch got.Kind {
	case parser.NodeNullLit:
		tc.errorf(decl, "`%s` is untyped: `null` carries no type and none was declared", decl.Lexeme)
		return unk()
	case parser.NodeVoidType:
		tc.errorf(decl, "cannot use a value that returns nothing to initialize `%s`", decl.Lexeme)
		return unk()
	}
	return defaultType(got)
}

// infer determines the type of an expression and reports any errors found
// while doing so. It returns an unknown type when it cannot decide, so it
// never produces false-positive errors.
func (tc *TypeChecker) infer(n *parser.Node) *Type {
	if n == nil {
		return unk()
	}
	// this should also work
	switch n.Type {
	case parser.NodeNumLit:
		return &Type{Kind: parser.NodeNumLit}
	case parser.NodeFloatLit:
		return &Type{Kind: parser.NodeFloatLit}
	case parser.NodeStrLit:
		return &Type{Kind: parser.NodeStrType}
	case parser.NodeNullLit:
		return &Type{Kind: parser.NodeNullLit}
	case parser.NodeLit:
		return tc.inferIdentifier(n)
	case parser.NodeBinOp:
		return tc.inferBinOp(n)
	case parser.NodeUnary:
		return tc.inferUnary(n)
	case parser.NodeFunctionCall:
		return tc.inferCall(n)
	case parser.NodeIndex:
		return tc.inferIndex(n)
	case parser.NodeListInit, parser.NodeArray:
		return tc.inferList(n)
	case parser.NodeAnonymousFunction:
		if info, ok := n.ExtraInfo.(*parser.NodeInfoFunction); ok {
			ft := funcType(info)
			tc.checkFunctionBody(info, ft)
			return ft
		}
	}

	return unk() //unknown
}

func (tc *TypeChecker) inferIdentifier(n *parser.Node) *Type {
	if t, ok := tc.lookup(n.Lexeme); ok {
		return t
	}
	if n.Lexeme == "true" || n.Lexeme == "false" {
		return &Type{Kind: parser.NodeBoolType}
	}
	// Undeclared names shouldn't be a problem
	return unk()
}

func (tc *TypeChecker) inferBinOp(n *parser.Node) *Type {
	//the right-hand side is a name
	if n.Lexeme == "." {
		tc.infer(n.Left)
		return unk()
	}
	l := tc.infer(n.Left)
	r := tc.infer(n.Right)

	switch n.Lexeme {
	//arthematic operators
	case "+", "-", "*", "/", "%":
		if l.Kind == parser.NodeUnknownType || r.Kind == parser.NodeUnknownType {
			return unk()
		}
		if n.Lexeme == "+" && l.Kind == parser.NodeStrType && r.Kind == parser.NodeStrType {
			return &Type{Kind: parser.NodeStrType}
		}
		if !isNumeric(l.Kind) || !isNumeric(r.Kind) {
			tc.errorf(n, "operator `%s` can and not be applied to `%s` and `%s`", n.Lexeme, l, r)
			return unk()
		}
		t, ok := unify(l, r)
		if !ok {
			tc.errorf(n, "mismatched tipes `%s` and `%s` for operator `%s`", l, r, n.Lexeme)
			return unk()
		}
		return t

		// equal and un-equal
	case "==", "!=":
		if !assignable(l, r) && !assignable(r, l) {
			tc.errorf(n, "can't compare `%s` with `%s`", l, r)
		}
		return &Type{Kind: parser.NodeBoolType}
		//
	case "<", ">", "<=", ">=":
		if l.Kind == parser.NodeUnknownType || r.Kind == parser.NodeUnknownType {
			return &Type{Kind: parser.NodeBoolType}
		}
		numeric := isNumeric(l.Kind) && isNumeric(r.Kind)
		text := l.Kind == parser.NodeStrType && r.Kind == parser.NodeStrType
		if !numeric && !text {
			tc.errorf(n, "operator `%s` cannot be applied to `%s` and `%s`", n.Lexeme, l, r)
		} else if _, ok := unify(l, r); !ok {
			tc.errorf(n, "mismatcheed typess `%s` and `%s` for operator `%s`", l, r, n.Lexeme)
		}
		return &Type{Kind: parser.NodeBoolType}
		// and or operators
	case "&&", "||":
		for _, side := range []*Type{l, r} {
			if side.Kind != parser.NodeUnknownType && side.Kind != parser.NodeBoolType {
				tc.errorf(n, "operator `%s` needs `Bool` operands, got `%s`", n.Lexeme, side)
			}
		}
		return &Type{Kind: parser.NodeBoolType}
		// remaining
	case "=", "+=", "-=", "*=", "/=", "%=":
		if !assignable(l, r) {
			tc.errorf(n, "expected `%s`, got `%s`", l, r)
		}
		return l
	}

	return unk()
}

func (tc *TypeChecker) inferUnary(n *parser.Node) *Type {
	t := tc.infer(operand(n))

	switch n.Lexeme {
	case "-":
		if t.Kind == parser.NodeUnknownType {
			return t
		}
		if !isNumeric(t.Kind) {
			tc.errorf(n, "operator `-` cannot be applied to `%s`", t)
			return unk()
		}
		return t

	case "!":
		switch t.Kind {
		case parser.NodeUnknownType:
			return unk()
		case parser.NodeBoolType:
			return t
		case parser.NodeFunction: // `fn!` is a zer0argument call
			if len(t.Params) != 0 {
				tc.errorf(n, "expected %d argument(s), got 0", len(t.Params))
			}
			return t.Elem
		}
		tc.errorf(n, "operator `!` needs a `Bool`, got `%s`", t)
		return unk()

	case "&":
		return &Type{Kind: parser.NodePointer, Elem: t}

	case "*":
		switch t.Kind {
		case parser.NodeUnknownType:
			return unk()
		case parser.NodePointer:
			return t.Elem
		}
		tc.errorf(n, "cannot dereference `%s`", t)
		return unk()
	}

	return unk()
}

func (tc *TypeChecker) inferCall(n *parser.Node) *Type {
	info, ok := n.ExtraInfo.(*parser.NodeInfoFunctionCall)
	if !ok {
		return unk()
	}

	callee := tc.infer(info.Callee)
	args := make([]*Type, len(info.Args))
	for i, a := range info.Args {
		args[i] = tc.infer(a)
	}

	if callee.Kind != parser.NodeFunction {
		// Builtin/module functions are still unknown here.
		return unk()
	}

	if len(args) != len(callee.Params) {
		tc.errorf(n, "expected %d argument(s), got %d", len(callee.Params), len(args))
		return callee.Elem
	}

	for i, got := range args {
		if !assignable(callee.Params[i], got) {
			tc.errorf(info.Args[i], "argument %d: expected `%s`, got `%s`", i+1, callee.Params[i], got)
		}
	}

	return callee.Elem
}

func (tc *TypeChecker) inferIndex(n *parser.Node) *Type {
	coll := tc.infer(n.Left)
	idx := tc.infer(n.Right)

	if idx.Kind != parser.NodeUnknownType && !isInteger(idx.Kind) {
		tc.errorf(n, "index must be an integer, got `%s`", idx)
	}

	if coll.Kind == parser.NodeArray {
		return coll.Elem
	}
	return unk()
}

// inferList type-checks list / array initializers.
// When an element type is declared, every element must be assignable to it.
// Otherwise all elements must unify with each other.
func (tc *TypeChecker) inferList(n *parser.Node) *Type {
	block, ok := n.ExtraInfo.(*parser.NodeInfoBlock)
	if !ok {
		return unk()
	}

	var want *Type
	if n.Left != nil {
		t := fromTypeNode(n.Left)
		switch t.Kind {
		case parser.NodeUnknownType:
			// Likely a struct initializer
			return unk()
		case parser.NodeArray:
			want = t.Elem
		default:
			want = t
		}
	}

	elem := unk()
	for _, el := range block.Elements {
		got := tc.infer(el)

		if want != nil {
			if !assignable(want, got) {
				tc.errorf(el, "list element: expected `%s`, got `%s`", want, got)
			}
			continue
		}

		u, ok := unify(elem, got)
		if !ok {
			tc.errorf(el, "list elements must have the same type: expected `%s`, got `%s`", elem, got)
			continue
		}
		elem = u
	}

	if want != nil {
		elem = want
	}
	return &Type{Kind: parser.NodeArray, Elem: elem}
}
