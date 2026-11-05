package resolver

import "github.com/KrabbaTek/krabbascript/internal/parser"

var builtinTypeNames = map[parser.NodeType]string{
	parser.NodeI8Type:   "I8",
	parser.NodeI16Type:  "I16",
	parser.NodeI32Type:  "I32",
	parser.NodeI64Type:  "I64",
	parser.NodeU8Type:   "U8",
	parser.NodeU16Type:  "U16",
	parser.NodeU32Type:  "U32",
	parser.NodeU64Type:  "U64",
	parser.NodeF32Type:  "F32",
	parser.NodeF64Type:  "F64",
	parser.NodeStrType:  "Str",
	parser.NodeBoolType: "Bool",
	parser.NodeAnyType:  "Any",
}

// typeString will turn a type node into a canonical string such as "I32"or
// "Point" and two type nodes are the same type exactly
// when their strings are equal.
func typeString(n *parser.Node) string {
	if n == nil {
		return "?"
	}
	switch n.Type {
	case parser.NodeLit:
		return n.Lexeme
	case parser.NodeArray:
		return "Arr<" + typeString(n.Left) + ">"
	case parser.NodePointer:
		return "*" + typeString(n.Left)
	}

	if name, ok := builtinTypeNames[n.Type]; ok {
		return name
	}

	return n.Type.String()
}
