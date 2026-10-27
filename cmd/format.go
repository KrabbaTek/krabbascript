package cmd

// exposes the Formatter library to the outside world via the CLI.
// using the 'format' command, you can run the formatter (ASI) on any given file.

// - Reece Turner (reeceturner1358@gmail.com)

import(
	"github.com/spf13/cobra"
	"fmt"
	"os"

	"github.com/KrabbaTek/krabbascript/internal/formatter" // formatter

	"strings"
)


var formatCmd = &cobra.Command{
	Use: "format <file>",
	Short: "run the formatter on a file",
	Args: cobra.MinimumNArgs(1), // require a file to be passed
	DisableFlagsInUseLine: true, // format doesn't need any flags and this makes the usage message more accurate
	Long: `
	As of 10/26/2026, this will only run the ASI (Automatic Semicolon Insertion) engine on the given file.
	This should be done before compiling a KrabbaScript file.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		// read file
		raw_lines, read_error := os.ReadFile(file)

		if read_error != nil {
			fmt.Println("Error reading file:", file)
			return
		}

		// get an array of strings from the weird data structure Golang decides to return
		lines := strings.Split(string(raw_lines), "\n")

		new_lines := formatter.InsertSemicolons(lines)

		// write new lines to file
		// Golang is super pretentious with its arrays of bytes so please excuse the second argument——
		// it's just packaging new_lines in a way that Golang can tolerate.
		write_error := os.WriteFile(file, []byte(strings.Join(new_lines, "\n")), 0644)

		if write_error != nil {
			fmt.Println("An error occured while writing your new code to the output file.")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
}