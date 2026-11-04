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
	Args: cobra.ExactArgs(1), // require a file to be passed
	DisableFlagsInUseLine: true, // format doesn't need any flags and this makes the usage message more accurate
	Long: `
	Run the ASI (Automatic Semicolon Insertion) engine on the given file or directory.
	This should be done before compiling a KrabbaScript file.

	If a directory is passed, then every KrabbaScript file within that directory will be formatted.


	Examples:

	krabbascript format denmark.ks

	krabbascript format ./europe/


	You're using the format command, 
	but consider using the --format flag in the build command — 
	it's a bit more convenient if you plan on compiling your formatted file.

	You can read a bit more about the engine here:
	https://github.com/KrabbaTek/krabbascript/blob/dev/DESIGN.md#formatter
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// once again, this can be either a file or directory.
		path := args[0]

		// used to check if the path is a file or directory later on
		stats, err := os.Stat(path)

		if err != nil {
			return err
		}

		// route to the correct formatter function based on whether the path type
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
