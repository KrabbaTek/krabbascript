package cmd

// exposes the Formatter library to the outside world via the CLI.
// using the 'format' command, you can run the formatter (ASI) on any given file.

// In addition to the format command, you can also pass the format flag to the build command,
// which will format your file/directory before building.

// - Reece Turner (reeceturner1358@gmail.com)

import (
	"github.com/spf13/cobra"

	"github.com/KrabbaTek/krabbascript/internal/formatter" // formatter
)

var formatCmd = &cobra.Command{
	Use:                   "format <file>",
	Short:                 "Run the formatter on a file",
	Args:                  cobra.ExactArgs(1), // require a file to be passed
	DisableFlagsInUseLine: true,               // format doesn't need any flags and this makes the usage message more accurate
	Long: `
	Run the ASI (Automatic Semicolon Insertion) engine on the given file.
	This should be done before compiling a KrabbaScript file.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		err := formatter.FormatFile(file)

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
}
