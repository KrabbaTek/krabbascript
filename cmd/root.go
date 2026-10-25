package cmd

import (
	"os"

	"github.com/KrabbaTek/krabbascript/internal/version"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kscript",
	Short:   "KrabbaScript compiler",
	Version: version.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
