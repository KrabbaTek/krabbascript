package importresolver

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/KrabbaTek/krabbascript/internal/parser"
)

const kscriptDir string = ".kscript"

type ImportResolver struct {
	libPaths []string
	errors   []error
}

func GetKscriptHome() (string, error) {
	// Get the home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	kscriptHome := path.Join(home, kscriptDir)
	return kscriptHome, nil
}

// Takes in the absolute path to the libraries folders, the current file being parsed
// and creates an Import Resolver from that
func NewImportResolver(libPaths []string) (ImportResolver, error) {
	// Make a copy of libPaths and append the global libs folder
	paths := libPaths

	// Get KrabbaScript's home
	home, err := GetKscriptHome()
	if err != nil {
		return ImportResolver{}, nil
	}

	// Join the home path and add libs
	paths = append(paths, path.Join(home, "libs"))

	return ImportResolver{
		libPaths: paths,
	}, nil
}

func (ir *ImportResolver) resolveImport(modImport string) {

}

func (ir *ImportResolver) resolveImports(imports []*parser.Node) {
	// Walk all of the library folders
	for _, libDir := range ir.libPaths {
		filepath.WalkDir(libDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Skip folders like .git
			if d.IsDir() && strings.HasPrefix(d.Name(), ".") && path != libDir {
				return fs.SkipDir
			}

			// Print every .ks file path for now
			if !d.IsDir() && filepath.Ext(path) == ".ks" {
				fmt.Println(path)
			}

			return nil
		})
	}
}

func (ir *ImportResolver) ResolveLibs(ast *parser.Node) []error {
	// Walk though the AST and find all of the import statements
	// But first, try get a NodeInfoBlock from the AST's root
	block, ok := ast.ExtraInfo.(*parser.NodeInfoBlock)
	if !ok {
		ir.errors = append(ir.errors, fmt.Errorf("expected the ExtraInfo of AST root to be a block"))
		return ir.errors
	}

	// Now we can walk the AST
	var imports []*parser.Node

	for _, node := range block.Elements {
		if node.Type == parser.NodeImportStatement {
			imports = append(imports, node)
		}
	}

	ir.resolveImports(imports)
	return ir.errors
}
