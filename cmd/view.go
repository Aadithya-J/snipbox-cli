package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var viewCmd = &cobra.Command{
	Use:   "view [slug]",
	Short: "View a snip",
	Long:  `View a snip by providing its slug.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := args[0]
		url := fmt.Sprintf("https://snipbox.vercel.app/api/clips/%s", slug)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error viewing snip: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Print("Enter password: ")
			password, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Fprintf(os.Stderr, "\nError reading password: %v\n", err)
				os.Exit(1)
			}
			fmt.Println()

			url = fmt.Sprintf("%s?password=%s", url, string(password))
			resp, err = http.Get(url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error viewing snip: %v\n", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			body, _ = ioutil.ReadAll(resp.Body)
		}

		if resp.StatusCode != http.StatusOK {
			var errorResult map[string]string
			json.Unmarshal(body, &errorResult)
			fmt.Fprintf(os.Stderr, "Error: %s\n", errorResult["error"])
			os.Exit(1)
		}

		var result map[string]interface{}
		json.Unmarshal(body, &result)

		fmt.Println(result["content"])
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
