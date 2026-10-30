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
	NodeModule       // module

	NodeWhenArm     // when arm
	NodeWhenElseArm // when else arm
	NodeScope       // scope
	NodeLater       // later
	NodeSic         // sic

	NodeReturn            // return
	NodeAnonymousFunction // anonymous function

	NodeFieldDec  // field dec
	NodeStructDec // struct dec
	NodeListInit  // list init

	NodeFunctionArg // function arg
	NodeFunction    // function

	// Compiler stuff
	NodeRoot         // root
	NodeBinOp        // bin op
	NodeUnary        // unary
	NodeNumLit       // num lit
	NodeFloatLit     // float lit
	NodeNullLit      // null lit
	NodeLit          // lit
	NodeStrLit       // str lit
	NodeArray        // array
	NodePointer      // pointer
	NodeIndex        // index
	NodeFunctionCall // function call

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
	Line, Column int

	ExtraInfo any
}

type NodeInfoBlock struct {
	Elements []*Node
}

type NodeInfoFunction struct {
	Args       []*Node
	ReturnType *Node
	Body       *Node
}

type NodeInfoFunctionCall struct {
	Args   []*Node
	Callee *Node
}

func (n *Node) PrintWithDepth(depth int) {
	if n == nil {
		return
	}

	offset := strings.Repeat("   ", depth)

	// Print node's type and lexeme
	if n.Lexeme == "" {
		// Don't print the lexeme if it's empty
		fmt.Printf("%s%s", offset, n.Type)
	} else {
		fmt.Printf("%s%s with lexeme: %s", offset, n.Type, n.Lexeme)
	}

	// Print any extra info the node carries
	if n.ExtraInfo != nil {
		switch info := n.ExtraInfo.(type) {
		case *NodeInfoBlock:
			fmt.Println(" -> block:")

			for i, el := range info.Elements {
				fmt.Printf("%s   [%d]:\n", offset, i)
				el.PrintWithDepth(depth + 2)
			}
		case *NodeInfoFunction:
			fmt.Println(" -> function:")

			fmt.Printf("%s   return:\n", offset)
			info.ReturnType.PrintWithDepth(depth + 2)

			fmt.Printf("%s   body:\n", offset)
			info.Body.PrintWithDepth(depth + 2)

			fmt.Printf("%s   args:\n", offset)
			for i, arg := range info.Args {
				fmt.Printf("%s      [%d]:\n", offset, i)
				arg.PrintWithDepth(depth + 3)
			}
		case *NodeInfoFunctionCall:
			fmt.Println(" -> call:")

			fmt.Printf("%s   args:\n", offset)
			for i, arg := range info.Args {
				fmt.Printf("%s      [%d]:\n", offset, i)
				arg.PrintWithDepth(depth + 3)
			}
			fmt.Printf("%s   callee:\n", offset)
			info.Callee.PrintWithDepth(depth + 2)
		default:
			fmt.Printf(" -> einfo: %v\n", info)
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
