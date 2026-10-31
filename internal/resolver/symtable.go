package resolver

import (
	"fmt"

	"github.com/KrabbaTek/krabbascript/internal/parser"
)

type SymTableEntryType int

//go:generate go tool stringer -type=SymTableEntryType -linecomment
const (
	SymTableVariableDef SymTableEntryType = iota // variable def
	SymTableVariableDec                          // variable dec
	SymTableValueDef                             // value def

	SymTableImport // import
)

type SymTableEntry struct {
	Kind SymTableEntryType // The Kind field tells what the symbol is, a function, a variable, etc
	Type *parser.Node

	Ref *parser.Node // Reference to the node
}

// Since there's no lexeme in the symtable entry and is rathed stored as a key in the hashmap,
// this function takes in the lexeme.
func (s SymTableEntry) Print(lexeme string) {
	fmt.Printf("symbol %s with kind %s:\n", lexeme, s.Kind)

	fmt.Println("type:")
	s.Type.Print()

	fmt.Println("ref:")
	s.Ref.Print()
}
