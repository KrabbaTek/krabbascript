package typechecker

import "github.com/KrabbaTek/krabbascript/internal/parser"

type Type struct {
	Kind   parser.NodeType
	Elem   *Type
	Params []*Type
}

func unk() *Type { return &Type{Kind: parser.NodeUnknownType} }

var kindNames = map[parser.NodeType]string{
	parser.NodeUnknownType: "unknown",
	parser.NodeVoidType:    "void",
	parser.NodeNumLit:      "I32",
	parser.NodeFloatLit:    "F64",
	parser.NodeNullLit:     "null",
	parser.NodeI8Type:      "I8",
	parser.NodeI16Type:     "I16",
	parser.NodeI32Type:     "I32",
	parser.NodeI64Type:     "I64",
	parser.NodeU8Type:      "U8",
	parser.NodeU16Type:     "U16",
	parser.NodeU32Type:     "U32",
	parser.NodeU64Type:     "U64",
	parser.NodeF32Type:     "F32",
	parser.NodeF64Type:     "F64",
	parser.NodeStrType:     "Str",
	parser.NodeBoolType:    "Bool",
	parser.NodeAnyType:     "Any",
}

func (t *Type) String() string {
	if t == nil {
		return "unknown"
	}

	switch t.Kind {
	case parser.NodePointer:
		return "*" + t.Elem.String()
	case parser.NodeArray:
		return "list of " + t.Elem.String()
	case parser.NodeFunction:
		return "function"
	}

	return kindNames[t.Kind]
}

func isInteger(k parser.NodeType) bool {
	switch k {
	case parser.NodeNumLit,
		parser.NodeI8Type, parser.NodeI16Type, parser.NodeI32Type, parser.NodeI64Type,
		parser.NodeU8Type, parser.NodeU16Type, parser.NodeU32Type, parser.NodeU64Type:
		return true
	}
	return false
}

func isFloat(k parser.NodeType) bool {
	return k == parser.NodeFloatLit || k == parser.NodeF32Type || k == parser.NodeF64Type
}

func isNumeric(k parser.NodeType) bool { return isInteger(k) || isFloat(k) }

func isBuiltin(k parser.NodeType) bool {
	switch k {
	case parser.NodeStrType, parser.NodeBoolType, parser.NodeAnyType:
		return true
	}
	return isNumeric(k) && k != parser.NodeNumLit && k != parser.NodeFloatLit
}

func firstNode(a, b *parser.Node) *parser.Node {
	if a != nil {
		return a
	}
	return b
}

func fromTypeNode(n *parser.Node) *Type {
	if n == nil {
		return unk()
	}

	if isBuiltin(n.Type) {
		return &Type{Kind: n.Type}
	}

	switch n.Type {
	case parser.NodePointer:
		return &Type{Kind: parser.NodePointer, Elem: fromTypeNode(firstNode(n.Left, n.Right))}
	case parser.NodeArray:
		return &Type{Kind: parser.NodeArray, Elem: fromTypeNode(firstNode(n.Left, n.Right))}
	}

	return unk()
}

func defaultType(t *Type) *Type {
	if t == nil {
		return t
	}

	switch t.Kind {
	case parser.NodeNumLit:
		return &Type{Kind: parser.NodeI32Type}
	case parser.NodeFloatLit:
		return &Type{Kind: parser.NodeF64Type}
	case parser.NodeArray:
		return &Type{Kind: parser.NodeArray, Elem: defaultType(t.Elem)}
	}

	return t
}

func equal(a, b *Type) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Kind == parser.NodeUnknownType || b.Kind == parser.NodeUnknownType {
		return true
	}
	if a.Kind != b.Kind {
		return false
	}

	switch a.Kind {
	case parser.NodePointer, parser.NodeArray:
		return equal(a.Elem, b.Elem)
	}

	return true
}

func assignable(dst, src *Type) bool {
	if dst == nil || src == nil || dst.Kind == parser.NodeUnknownType || src.Kind == parser.NodeUnknownType {
		return true
	}
	if dst.Kind == parser.NodeAnyType {
		return src.Kind != parser.NodeVoidType
	}

	switch {
	case src.Kind == parser.NodeNumLit:
		return isNumeric(dst.Kind)
	case src.Kind == parser.NodeFloatLit:
		return isFloat(dst.Kind)
	case src.Kind == parser.NodeNullLit:
		return dst.Kind == parser.NodePointer
	case dst.Kind == parser.NodeArray && src.Kind == parser.NodeArray:
		return assignable(dst.Elem, src.Elem)
	}

	return equal(dst, src)
}

func unify(a, b *Type) (*Type, bool) {
	if a.Kind == parser.NodeUnknownType {
		return b, true
	}
	if b.Kind == parser.NodeUnknownType {
		return a, true
	}

	switch {
	case a.Kind == parser.NodeNumLit && isNumeric(b.Kind):
		return b, true
	case b.Kind == parser.NodeNumLit && isNumeric(a.Kind):
		return a, true
	case a.Kind == parser.NodeFloatLit && isFloat(b.Kind):
		return b, true
	case b.Kind == parser.NodeFloatLit && isFloat(a.Kind):
		return a, true
	}

	if equal(a, b) {
		return a, true
	}

	return nil, false
}
