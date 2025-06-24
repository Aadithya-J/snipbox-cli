package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snipbox",
	Short: "A CLI for Snipbox",
	Long:  `A simple and secure text sharing platform with a CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Snipbox CLI!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
