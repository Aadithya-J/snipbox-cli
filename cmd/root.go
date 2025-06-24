package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "snipbox",
	Version: "1.0.2",
	Short:   "A CLI for Snipbox",
}

func init() {
	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
