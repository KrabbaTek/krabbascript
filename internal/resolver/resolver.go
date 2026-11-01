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

func (r *Resolver) popScope() {
	r.symTable = r.symTable[:len(r.symTable)-1]
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

	err := r.resolveExpression(variable.Right)
	if err != nil {
		return err
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

	err := r.resolveExpression(variable.Right)
	if err != nil {
		return err
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

	err := r.resolveExpression(value.Right)
	if err != nil {
		return err
	}

	currentScope[value.Lexeme] = &SymTableEntry{
		Kind: SymTableValueDef,
		Type: value.Left,
		Ref:  value,
	}

	return nil
}

func (r *Resolver) resolveImport(modImport *parser.Node) *ResolverError {
	currentScope := r.getCurrentScope()

	// Get the NodeInfoBlock field
	info, ok := modImport.Right.ExtraInfo.(*parser.NodeInfoBlock)
	if !ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: expected the ExtraInfo of expression to be a block", r.file,
				modImport.Line, modImport.Column),
			Fix: "this is a programming error",
		}
	}

	module := info.Elements[0]
	// Check if the module is already imported

	_, ok = r.lookUp(module.Lexeme)
	if ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: reimport of '%s'", r.file, modImport.Line, modImport.Column, module.Lexeme),
		}
	}

	currentScope[module.Lexeme] = &SymTableEntry{
		Kind: SymTableImport,
		Type: nil,
		Ref:  modImport,
	}

	return nil
}

func (r *Resolver) resolveIfStatement(ifStat *parser.Node) *ResolverError {
	// Create a new scope for the if statement
	r.symTable = append(r.symTable, make(map[string]*SymTableEntry))
	defer r.popScope() // We should pop the scope after so you can't access variables inside a scope

	// Resolve if's conditional which lives on the left
	err := r.resolveExpression(ifStat.Left)
	if err != nil {
		return err
	}

	r.walkBlock(ifStat)

	return nil
}

func (r *Resolver) resolveElseStatement(elseStat *parser.Node) *ResolverError {
	// Create a new scope for the elsif statement
	r.symTable = append(r.symTable, make(map[string]*SymTableEntry))
	defer r.popScope() // We should pop the scope after so you can't access variables inside a scope

	r.walkBlock(elseStat)

	return nil
}

func (r *Resolver) resolveExpression(expr *parser.Node) *ResolverError {
	if expr == nil {
		return nil
	}

	switch expr.Type {
	case parser.NodeBinOp:
		switch expr.Lexeme {
		case ".":
			// Only check if left exists
			err := r.resolveExpression(expr.Left)
			if err != nil {
				return err
			}
		default:
			// Check both left and right
			err := r.resolveExpression(expr.Left)
			if err != nil {
				return err
			}

			// Only check if left exists
			err = r.resolveExpression(expr.Right)
			if err != nil {
				return err
			}
		}
	case parser.NodeUnary:
		// Only check if right exists
		err := r.resolveExpression(expr.Right)
		if err != nil {
			return err
		}
	case parser.NodeLit:
		// Look up if it exists
		_, ok := r.lookUp(expr.Lexeme)
		if !ok {
			return &ResolverError{
				What: fmt.Errorf("%s:%d:%d: undeclared literal '%s'", r.file,
					expr.Line, expr.Column, expr.Lexeme),
			}
		}
	case parser.NodeFunctionCall:
		// Get the NodeInfoFunctionCall field
		info, ok := expr.ExtraInfo.(*parser.NodeInfoFunctionCall)
		if !ok {
			return &ResolverError{
				What: fmt.Errorf("%s:%d:%d: expected the ExtraInfo of expression to be a function call", r.file,
					expr.Line, expr.Column),
				Fix: "this is a programming error",
			}
		}

		// Look up if the function exists
		err := r.resolveExpression(info.Callee)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) resolveWhen(when *parser.Node) *ResolverError {
	// Resolve when's right arm containing a literal
	err := r.resolveExpression(when.Right)
	if err != nil {
		return err
	}

	// Try getting when's NodeInfoBlock
	block, ok := when.ExtraInfo.(*parser.NodeInfoBlock)
	if !ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: expected the ExtraInfo of when to be a block", r.file,
				when.Line, when.Column),
			Fix: "this is a programming error",
		}
	}

	for _, arm := range block.Elements {
		if arm.Type == parser.NodeWhenArm {
			// Resolve it's conditional
			err = r.resolveExpression(arm.Left.Left)
			if err != nil {
				return err
			}

			r.symTable = append(r.symTable, make(map[string]*SymTableEntry))
			defer r.popScope()

			r.walkBlock(arm.Right)
		} else {
			// It's an else arm, no need to resolve the conditional
			r.symTable = append(r.symTable, make(map[string]*SymTableEntry))
			defer r.popScope()

			r.walkBlock(arm.Right)
		}
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
	r.walkBlock(ast)

	return r.symTable, r.errors
}

func (r *Resolver) walkBlock(block *parser.Node) {
	v, ok := block.ExtraInfo.(*parser.NodeInfoBlock)
	if !ok {
		r.errors = append(r.errors, &ResolverError{
			What: fmt.Errorf("%s:%d:%d: expected the ExtraInfo of block to be a block", r.file,
				block.Line, block.Column),
			Fix: "this is a programming error",
		})
		return
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

		case parser.NodeImportStatement:
			err := r.resolveImport(el)
			if err != nil {
				r.errors = append(r.errors, err)
			}

		case parser.NodeWhenStatement:
			err := r.resolveWhen(el)
			if err != nil {
				r.errors = append(r.errors, err)
			}

		case parser.NodeIfStatement:
			err := r.resolveIfStatement(el)
			if err != nil {
				r.errors = append(r.errors, err)
			}
		case parser.NodeElsifStatement:
			// Reuse the if statement because Elsif and If statements are similar in structure
			err := r.resolveIfStatement(el)
			if err != nil {
				r.errors = append(r.errors, err)
			}
		case parser.NodeElseStatement:
			err := r.resolveElseStatement(el)
			if err != nil {
				r.errors = append(r.errors, err)
			}
		case parser.NodeScope:
			r.walkBlock(el)
		default:
			// r.errors = append(r.errors, &ResolverError{
			// 	What: fmt.Errorf("%s:%d:%d: not supported by the resolver",
			// 		r.file, el.Line, el.Column),
			// })

			err := r.resolveExpression(el)
			if err != nil {
				r.errors = append(r.errors, err)
			}
		}
	}

	block.Print()
}
