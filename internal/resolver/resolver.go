package resolver

import (
	"fmt"

	"github.com/KrabbaTek/krabbascript/internal/parser"
	"github.com/fatih/color"
)

type ResolverError struct {
	What error
	Fix  string
}

func (e ResolverError) String() string {
	if e.Fix != "" {
		return fmt.Sprintf("%s\n\t%s: %s", e.What, color.BlueString("fix"), e.Fix)
	} else {
		return fmt.Sprintf("%s", e.What)
	}
}

// The absolute resolver
type Resolver struct {
	file     string
	errors   []*ResolverError
	symTable []map[string]*SymTableEntry
}

func NewResolver(file string) Resolver {
	return Resolver{
		file: file,
	}
}

func (r *Resolver) getCurrentScope() map[string]*SymTableEntry {
	return r.symTable[len(r.symTable)-1]
}

func (r *Resolver) lookUp(lexeme string) (*SymTableEntry, bool) {
	// Look up if the lexeme exists
	for i := len(r.symTable) - 1; i >= 0; i-- {
		// Look up the lexeme
		entry, ok := r.symTable[i][lexeme]
		if ok {
			return entry, true
		}
	}

	return nil, false
}

func (r *Resolver) resolveVariableDef(variable *parser.Node) *ResolverError {
	currentScope := r.getCurrentScope()

	// Check if the varible already exists
	entry, ok := r.lookUp(variable.Lexeme)
	if ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: redeclaration of %s at %s:%d:%d", r.file, variable.Line, variable.Column, variable.Lexeme,
				r.file, entry.Ref.Line, entry.Ref.Line),
		}
	}

	currentScope[variable.Lexeme] = &SymTableEntry{
		Kind: SymTableVariableDef,
		Type: variable.Left,
		Ref:  variable,
	}

	return nil
}

func (r *Resolver) resolveVariableDec(variable *parser.Node) *ResolverError {
	currentScope := r.getCurrentScope()

	// Check if the varible already exists
	entry, ok := r.lookUp(variable.Lexeme)
	if ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: redeclaration of %s at %s:%d:%d", r.file, variable.Line, variable.Column, variable.Lexeme,
				r.file, entry.Ref.Line, entry.Ref.Line),
		}
	}

	currentScope[variable.Lexeme] = &SymTableEntry{
		Kind: SymTableVariableDec,
		Type: variable.Left,
		Ref:  variable,
	}

	return nil
}

func (r *Resolver) resolveValueDef(value *parser.Node) *ResolverError {
	currentScope := r.getCurrentScope()

	// Check if the varible already exists
	entry, ok := r.lookUp(value.Lexeme)
	if ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: redefinition of %s at %s:%d:%d", r.file, value.Line, value.Column, value.Lexeme,
				r.file, entry.Ref.Line, entry.Ref.Line),
		}
	}

	currentScope[value.Lexeme] = &SymTableEntry{
		Kind: SymTableValueDef,
		Type: value.Left,
		Ref:  value,
	}

	return nil
}

func (r *Resolver) Walk(ast *parser.Node) ([]map[string]*SymTableEntry, []*ResolverError) {
	// Validate that the root of the AST is actually a root type
	if ast.Type != parser.NodeRoot {
		r.errors = append(r.errors, &ResolverError{
			What: fmt.Errorf("the root of the AST must have a NodeRoot type"),
			Fix:  "this is a programming error",
		})

		return r.symTable, r.errors
	}

	// Create a global scope which will always exist
	r.symTable = append(r.symTable, make(map[string]*SymTableEntry))

	v, ok := ast.ExtraInfo.(*parser.NodeInfoBlock)
	if !ok {
		r.errors = append(r.errors, &ResolverError{
			What: fmt.Errorf("the root of the AST must have it's ExtraInfo as a *NodeInfoBlock"),
			Fix:  "this is a programming error",
		})
	}

	for _, el := range v.Elements {
		switch el.Type {
		case parser.NodeVariableDef:
			err := r.resolveVariableDef(el)
			if err != nil {
				r.errors = append(r.errors, err)
			}
		case parser.NodeVariableDec:
			err := r.resolveVariableDec(el)
			if err != nil {
				r.errors = append(r.errors, err)
			}
		case parser.NodeValueDef:
			err := r.resolveValueDef(el)
			if err != nil {
				r.errors = append(r.errors, err)
			}
		default:
			r.errors = append(r.errors, &ResolverError{
				What: fmt.Errorf("%s:%d:%d: not supported by the resolver",
					r.file, el.Line, el.Column),
			})
		}
	}

	ast.Print()
	return r.symTable, r.errors
}
