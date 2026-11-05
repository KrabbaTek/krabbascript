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
	module := modImport.Right

	// Check if the module is already imported
	entry, ok := r.lookUp(module.Lexeme)
	if ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: reimport of '%s' at %s:%d:%d", r.file, modImport.Line, modImport.Column, module.Lexeme,
				r.file, entry.Ref.Line, entry.Ref.Column),
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

		// Resolve the args
		for _, arg := range info.Args {
			err = r.resolveExpression(arg)
			if err != nil {
				return err
			}
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

func (r *Resolver) resolveFunction(function *parser.Node) *ResolverError {
	funct, ok := function.ExtraInfo.(*parser.NodeInfoFunction)
	if !ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: expected the ExtraInfo of function to be a function", r.file,
				function.Line, function.Column),
			Fix: "this is a programming error",
		}
	}

	// Look up if the function already exists
	entry, ok := r.lookUp(function.Lexeme)
	if ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: redeclaration of %s at %s:%d:%d", r.file, function.Line, function.Column, function.Lexeme,
				r.file, entry.Ref.Line, entry.Ref.Column),
		}
	}

	// Create a new scope for the function
	r.symTable = append(r.symTable, make(map[string]*SymTableEntry))

	// Get the function's scope
	scope := r.getCurrentScope()

	// Loop over the args and add bring them to the scope
	for _, arg := range funct.Args {
		scope[arg.Left.Lexeme] = &SymTableEntry{
			Kind: SymTableFunctionArg,
			Type: arg.Right,
			Ref:  arg,
		}
	}

	r.walkBlock(funct.Body)
	r.popScope()

	// Get the current scope
	scope = r.getCurrentScope()
	scope[function.Lexeme] = &SymTableEntry{
		Kind: SymTableFunction,
		Type: funct.ReturnType, // The type here is the return type
		Ref:  function,
	}

	return nil
}

func (r *Resolver) resolveMethod(method *parser.Node) *ResolverError {
	met, ok := method.ExtraInfo.(*parser.NodeInfoMethod)
	if !ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: expected the ExtraInfo of method to be a method", r.file,
				method.Line, method.Column),
			Fix: "this is a programming error",
		}
	}


	key := typeString(met.OfType) + "." + method.Lexeme

	entry, ok := r.lookUp(key)
	if ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: redeclaration of method %s on %s at %s:%d:%d", r.file, method.Line, method.Column,
				method.Lexeme, typeString(met.OfType), r.file, entry.Ref.Line, entry.Ref.Column),
		}
	}

	// Create a new scope for the method.
	r.symTable = append(r.symTable, make(map[string]*SymTableEntry))

	// Get the method's scope
	scope := r.getCurrentScope()

	// Loop over the args and bring them into the scope
	for _, arg := range met.Args {
		scope[arg.Left.Lexeme] = &SymTableEntry{
			Kind: SymTableMethodArg,
			Type: arg.Right,
			Ref:  arg,
		}
	}

	r.walkBlock(met.Body)
	r.popScope()

	// Register the method under its qualified key
	scope = r.getCurrentScope()
	scope[key] = &SymTableEntry{
		Kind: SymTableMethod,
		Type: met.ReturnType, // The type here is the return type
		Ref:  method,
	}

	return nil
}

func (r *Resolver) resolveStruct(ksStruct *parser.Node) *ResolverError {
	// We don't really care about the struct's fields right now.
	// Check if the struct is already defined

	entry, ok := r.lookUp(ksStruct.Lexeme)
	if ok {
		return &ResolverError{
			What: fmt.Errorf("%s:%d:%d: redeclaration of %s at %s:%d:%d", r.file, ksStruct.Line, ksStruct.Column,
				ksStruct.Lexeme, r.file, entry.Ref.Line, entry.Ref.Column),
		}
	}

	// Get the current scope
	scope := r.getCurrentScope()
	scope[ksStruct.Lexeme] = &SymTableEntry{
		Kind: SymTableStruct,
		Ref:  ksStruct,
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

	fmt.Println("AST:")
	ast.Print()

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

		case parser.NodeFunction:
			err := r.resolveFunction(el)
			if err != nil {
				r.errors = append(r.errors, err)
			}

		case parser.NodeMethod:
			err := r.resolveMethod(el)
			if err != nil {
				r.errors = append(r.errors, err)
			}

		case parser.NodeStructDec:
			err := r.resolveStruct(el)
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
}
