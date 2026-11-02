package cmd

import (
	"fmt"
	"os"

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

		l, err := lexer.NewLexer(dir)
		if err != nil {
			return err
		}

		toks, errors := l.Scan()

		// for _, tok := range toks {
		// 	fmt.Println(tok)
		// }

		if len(errors) > 0 {
			for _, lerr := range errors {
				fmt.Printf("kscript: %s: %v", color.RedString("error"), lerr)
			}

			fmt.Print("kscript: ")

			color.Set(color.FgRed)
			fmt.Printf("compilation failed with %d error(s)\n", len(errors))
			color.Unset()

			os.Exit(1)
		}

		p := parser.NewParser(toks, dir)
		ast, perrors := p.Parse()

		for _, err := range perrors {
			fmt.Printf("kscript: %s: %v\n", color.RedString("error"), err)
		}

		if len(perrors) != 0 {
			color.Set(color.FgRed)
			fmt.Printf("compilation failed with %d error(s)\n", len(errors))
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

		if len(rerrors) != 0 {
			color.Set(color.FgRed)
			fmt.Printf("compilation failed with %d error(s)\n", len(rerrors))
			color.Unset()

			os.Exit(1)
		}

		tc := typechecker.New(dir)
		terrs := tc.Check(ast)
		for _, err := range terrs {
			fmt.Printf("kscript: %s: %v\n", color.RedString("error"), err)
		}
		if len(terrs) != 0 {
			color.Set(color.FgRed)
			fmt.Printf("kscript: compilation failed with %d error(s)\n", len(terrs))
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
