package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/KrabbaTek/krabbascript/internal/importresolver"
	"github.com/KrabbaTek/krabbascript/internal/lexer"
	"github.com/KrabbaTek/krabbascript/internal/parser"
	"github.com/KrabbaTek/krabbascript/internal/resolver"
	"github.com/KrabbaTek/krabbascript/tools"

	"github.com/KrabbaTek/krabbascript/internal/formatter" // formatter

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/KrabbaTek/krabbascript/internal/typechecker" // type checker
)

var buildCmd = &cobra.Command{
	Use:   "build [directory/file]",
	Short: "Build a KrabbaScript project/file",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// refactor this into path since it could be either a file or directory
		dir := args[0]

		// user wants their code to be formatted before compilation
		format, _ := cmd.Flags().GetBool("format")

		if format {
			info, err := os.Stat(dir)

			if err != nil {
				return err
			}

			if info.IsDir() {
				// batch process — format each file in the directory
				formatter.FormatDirectory(dir)
			} else {
				// format a single file

				// This reads amazingly! Format my file named dir.
				err := formatter.FormatFile(dir)

				if err != nil {
					return err
				}
			}
		}

		f, err := tools.IsFile(dir)
		if err != nil {
			return err
		}

		if f == false {
			return fmt.Errorf("building directories is not supported")
		}

		home, err := importresolver.GetKscriptHome()
		if err != nil {
			return err
		}

		err = os.MkdirAll(path.Join(home, "libs"), 0755)
		if err != nil {
			return err
		}

		l, _ := lexer.NewLexer(dir)
		toks, errors := l.Scan()

		for _, err := range errors {
			fmt.Printf("kscript: %s: %v\n", color.RedString("error"), err)
		}

		if len(errors) == 1 {
			fmt.Print("kscript: ")

			color.Set(color.FgRed)
			fmt.Printf("compilation failed with %d error\n", len(errors))
			color.Unset()

			os.Exit(1)
		}

		if len(errors) > 0 {
			fmt.Print("kscript: ")

			color.Set(color.FgRed)
			fmt.Printf("compilation failed with %d errors\n", len(errors))
			color.Unset()

			os.Exit(1)
		}

		p := parser.NewParser(toks, dir)
		ast, perrors := p.Parse()

		for _, err := range perrors {
			fmt.Printf("kscript: %s: %v\n", color.RedString("error"), err)
		}

		if len(perrors) == 1 {
			fmt.Print("kscript: ")

			color.Set(color.FgRed)
			fmt.Printf("compilation failed with %d error\n", len(perrors))
			color.Unset()

			os.Exit(1)
		}
		if len(perrors) > 0 {
			color.Set(color.FgRed)
			fmt.Printf("compilation failed with %d errors\n", len(perrors))
			color.Unset()

			os.Exit(1)
		}

		ires, err := importresolver.NewImportResolver(make([]string, 0))
		if err != nil {
			return err
		}

		irerrors := ires.ResolveLibs(ast)
		for _, err := range irerrors {
			fmt.Printf("kscript: %s: %v\n", color.RedString("error"), err)
		}

		if len(irerrors) == 1 {
			color.Set(color.FgRed)
			fmt.Printf("compilation failed with %d error\n", len(irerrors))
			color.Unset()

			os.Exit(1)
		}
		if len(irerrors) > 0 {
			color.Set(color.FgRed)
			fmt.Printf("compilation failed with %d errors\n", len(irerrors))
			color.Unset()

			os.Exit(1)
		}

		res := resolver.NewResolver(dir)
		symTable, rerrors := res.Walk(ast)

		fmt.Println("symbol table:")
		for _, htable := range symTable {
			for v, k := range htable {
				k.Print(v)
			}
		}

		for _, err := range rerrors {
			fmt.Printf("kscript: %s: %v\n", color.RedString("error"), err)
		}

		if len(rerrors) == 1 {
			color.Set(color.FgRed)
			fmt.Printf("compilation failed with %d error\n", len(rerrors))
			color.Unset()

			os.Exit(1)
		}

		if len(rerrors) > 0 {
			color.Set(color.FgRed)
			fmt.Printf("compilation failed with %d errors\n", len(rerrors))
			color.Unset()

			os.Exit(1)
		}

		tc := typechecker.New(dir)
		terrs := tc.Check(ast)
		for _, err := range terrs {
			fmt.Printf("kscript: %s: %v\n", color.RedString("error"), err)
		}

		if len(terrs) == 1 {
			color.Set(color.FgRed)
			fmt.Printf("kscript: compilation failed with %d error\n", len(terrs))
			color.Unset()

			os.Exit(1)
		}
		if len(terrs) > 0 {
			color.Set(color.FgRed)
			fmt.Printf("kscript: compilation failed with %d errors\n", len(terrs))
			color.Unset()

			os.Exit(1)
		}

		return nil
	},
}

func init() {
	buildCmd.Flags().Bool("format", false, "Insert semicolons before building.")
	rootCmd.AddCommand(buildCmd)
}
