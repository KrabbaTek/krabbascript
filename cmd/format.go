package cmd

// exposes the Formatter library to the outside world via the CLI.
// using the 'format' command, you can run the formatter (ASI) on any given file.

// In addition to the format command, you can also pass the format flag to the build command,
// which will format your file/directory before building.

// - Reece Turner (reeceturner1358@gmail.com)

import(
	"github.com/spf13/cobra"

	"github.com/KrabbaTek/krabbascript/internal/formatter" // formatter

	"os"
)


var formatCmd = &cobra.Command{
	Use: "format [file/directory]",
	Short: "Run the formatter on a file or directory.",
	Args: cobra.MinimumNArgs(1), // require a file to be passed
	DisableFlagsInUseLine: true, // format doesn't need any flags and this makes the usage message more accurate
	Long: `
	Run the ASI (Automatic Semicolon Insertion) engine on the given file or directory.
	This should be done before compiling a KrabbaScript file.

	If a directory is passed, then every KrabbaScript file within that directory will be formatted.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		stats, err := os.Stat(path)

		if err != nil {
			return err
		}

		if stats.IsDir() {
			err = formatter.FormatDirectory(path)
		} else {
			err = formatter.FormatFile(path)
		}

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
}
