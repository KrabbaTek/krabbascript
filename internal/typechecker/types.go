package typechecker

import "github.com/KrabbaTek/krabbascript/internal/parser" // import parser

// Kind is what the checker see's a type is. Names match the
//surface syntax so error messages don't need a table.
type Kind int
const (
	// variable names with K. We can edit it later.
	KUnknown  Kind = iota 
	KVoid             
	KIntLit     // untyped integer literal, collapses to I32 at binding 
	KFloatLit   // same but F64 
	KNull
	KI8
	KI16
	KI32
	KI64
	KU8
	KU16
	KU32
	KU64
	KF32
	KF64
	KStr
	KBool
	KAny
	KPointer
	KList
	KFunc
)

type Type struct {
	Kind   Kind
	Elem   *Type   // list element
	Params []*Type 
}
func unk() *Type { return &Type{Kind: KUnknown} }
var kindNames = map[Kind]string{
		KUnknown:  "unknown",
		KVoid:     "void",
		KIntLit:   "integer",
		KFloatLit: "float",
		KNull:     "null",
		KI8:       "I8",
		KI16:      "I16",
		KI32:      "I32",
		KI64:      "I64",
		KU8:       "U8",
		KU16:      "U16",
		KU32:      "U32",
		KU64:      "U64",
		KF32:      "F32",
		KF64:      "F64",
		KStr:      "Str",
		KBool:     "Bool",
		KAny:      "Any",
}
//type 
func (t *Type) String() string {
	if t == nil {
		return "unknown"
	}
	switch t.Kind {
	case KPointer:
		return "*" + t.Elem.String()
	case KList:
		return "list of " + t.Elem.String()
	case KFunc:
		return "function"
	}

	return kindNames[t.Kind]
}
// The three numeric predicates are the only place that hardcodes
func isInteger(k Kind) bool { return k == KIntLit || (k >= KI8 && k <= KU64) }
func isFloat(k Kind) bool   { return k == KFloatLit || k == KF32 || k == KF64 }
func isNumeric(k Kind) bool { return isInteger(k) || isFloat(k) }
// typekinds 
var typeKinds = map[parser.NodeType]Kind{
	parser.NodeI8Type: KI8, 
	parser.NodeI16Type: KI16, 
	parser.NodeI32Type: KI32, 
	parser.NodeI64Type: KI64,
	parser.NodeU8Type: KU8, 
	parser.NodeU16Type: KU16, 
	parser.NodeU32Type: KU32, 
	parser.NodeU64Type: KU64,
	parser.NodeF32Type: KF32, parser.NodeF64Type: KF64,
	parser.NodeStrType: KStr, parser.NodeBoolType: KBool, 
	parser.NodeAnyType: KAny,
}
func firstNode(a, b *parser.Node) *parser.Node {
	if a != nil {
		return a
	}
	return b
}
// this converts a type annotation from the AST.
func fromTypeNode(n *parser.Node) *Type {
	if n == nil {
		return unk()
	}


	
	if k, ok := typeKinds[n.Type]; ok {
		return &Type{Kind: k}
	}
	switch n.Type {
	case parser.NodePointer:
		return &Type{Kind: KPointer, Elem: fromTypeNode(firstNode(n.Left, n.Right))}
	case parser.NodeArray:
		return &Type{Kind: KList, Elem: fromTypeNode(firstNode(n.Left, n.Right))}
	}


	return unk()
}
// this checks untyped literals down at binding time
func defaultType(t *Type) *Type {
	if t == nil {
		return t
	}
	switch t.Kind {
	case KIntLit:
		return &Type{Kind: KI32}
	case KFloatLit:
		return &Type{Kind: KF64}
	case KList:
		return &Type{Kind: KList, Elem: defaultType(t.Elem)}
	}
	return t
}

// equal compares two types and Unknown matches anything. for
// now this should work fine. 
func equal(a, b *Type) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Kind == KUnknown || b.Kind == KUnknown {
		return true
	}
	if a.Kind != b.Kind {
		return false
	}
	switch a.Kind {
	case KPointer, KList:
		return equal(a.Elem, b.Elem)
	}


	
	return true
}
// assignable is equal with the literal/implicit-conversion rules
// keep it small...
func assignable(dst, src *Type) bool {
	if dst == nil || src == nil || dst.Kind == KUnknown || src.Kind == KUnknown {
		return true
	}
	if dst.Kind == KAny {
		return src.Kind != KVoid
	}
	switch {
	case src.Kind == KIntLit:
		return isNumeric(dst.Kind)
	case src.Kind == KFloatLit:
		return isFloat(dst.Kind)
	case src.Kind == KNull:
		return dst.Kind == KPointer
	case dst.Kind == KList && src.Kind == KList:
		return assignable(dst.Elem, src.Elem)
	}
	return equal(dst,src) // unit test equal
}

//unify
func unify(a, b *Type) (*Type, bool) {
	if a.Kind == KUnknown {
		return b, true
	}
	if b.Kind == KUnknown {
		return a, true
	}
	switch {  // this is showing errrs.
	case a.Kind == KIntLit && isNumeric(b.Kind):
		return b, true
	case b.Kind == KIntLit && isNumeric(a.Kind):
		return a, true
	case a.Kind == KFloatLit && isFloat(b.Kind):
		return b, true
	case b.Kind == KFloatLit && isFloat(a.Kind):
		return a, true
	}
	if equal(a, b) {
		return a, true
	}
	return nil, false
}
