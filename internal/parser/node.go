package parser

import (
	"fmt"
	"strings"
)

//go:generate go tool stringer -type=NodeType -linecomment
type NodeType int

const (
	NodeVariableDef NodeType = iota // variable def
	NodeVariableDec                 // variable dec

	NodeValueDef        // value def
	NodeIfStatement     // if statement
	NodeElsifStatement  // elsif statement
	NodeElseStatement   // else statement
	NodeImportStatement // import statement
	NodeWhenStatement   // when statement

	NodeModuleImport // module import
	NodeModulePath   // module path

	NodeWhenArm     // when arm
	NodeWhenElseArm // when else arm
	NodeScope       // scope

	NodeFieldDec  // field dec
	NodeStructDec // struct dec
	NodeListInit  // list init

	// Compiler stuff
	NodeRoot     // root
	NodeBinOp    // bin op
	NodeNumLit   // num lit
	NodeFloatLit // float lit
	NodeLit      // lit
	NodeStrLit   // str lit
	NodeArray    // array
	NodeIndex    // index

	// Types
	NodeI64Type // I64 type
	NodeI32Type // I32 type
	NodeI16Type // I16 type
	NodeI8Type  // I8 type

	NodeU64Type // U64 type
	NodeU32Type // U32 type
	NodeU16Type // U16 type
	NodeU8Type  // U8 type

	NodeF64Type // F64 type
	NodeF32Type // F32 type

	NodeStrType  // Str type
	NodeBoolType // Bool type
	NodeAnyType  // Any type
)

type Node struct {
	Left  *Node
	Right *Node
	Type  NodeType

	Lexeme       string
	line, column int

	ExtraInfo any
}

type NodeInfoBlock struct {
	Elements []*Node
}

func (n *Node) PrintWithDepth(depth int) {
	if n == nil {
		return
	}

	offset := strings.Repeat("   ", depth)

	// Print nodes type and lexeme
	if n.Lexeme == "" {
		// Don't print the lexeme if it's empty
		fmt.Printf("%s%s", offset, n.Type)
	} else {
		fmt.Printf("%s%s with lexeme: %s", offset, n.Type, n.Lexeme)
	}

	// Print any extra info the node carries
	if n.ExtraInfo != nil {
		switch i := n.ExtraInfo.(type) {
		case *NodeInfoBlock:
			fmt.Println(" -> block:")
			for i, el := range i.Elements {
				fmt.Printf("%s   [%d]:\n", offset, i)
				el.PrintWithDepth(depth + 2)
			}
		default:
			fmt.Printf(" -> einfo: %v\n", i)
		}
	} else {
		fmt.Println()
	}

	// Print out the left and right node
	if n.Left != nil {
		fmt.Printf("%s   left:\n", offset)
		n.Left.PrintWithDepth(depth + 2)
	}
	if n.Right != nil {
		fmt.Printf("%s   right:\n", offset)
		n.Right.PrintWithDepth(depth + 2)
	}
}

func (n *Node) Print() {
	n.PrintWithDepth(1)
}
