package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish [file-path]",
	Short: "Publish a new snip",
	Long:  `Publish a new snip from a file or standard input.`,
	Run: func(cmd *cobra.Command, args []string) {
		var content []byte
		var err error

		if len(args) > 0 {
			content, err = ioutil.ReadFile(args[0])
		} else {
			content, err = ioutil.ReadAll(os.Stdin)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading content: %v\n", err)
			os.Exit(1)
		}

		if len(content) == 0 {
			fmt.Fprintln(os.Stderr, "Content cannot be empty")
			os.Exit(1)
		}

		password, _ := cmd.Flags().GetString("password")
		customSlug, _ := cmd.Flags().GetString("slug")
		expiresIn, _ := cmd.Flags().GetString("expires")
		destroyOnView, _ := cmd.Flags().GetBool("destroy-on-view")

		payload := map[string]interface{}{
			"content":       string(content),
			"isPrivate":     password != "",
			"password":      password,
			"customSlug":    customSlug,
			"expiresIn":     expiresIn,
			"destroyOnView": destroyOnView,
		}

		jsonPayload, _ := json.Marshal(payload)

		resp, err := http.Post("https://snipbox.vercel.app/api/clips", "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error publishing snip: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "Error: %s\n", string(body))
			os.Exit(1)
		}

		var result map[string]string
		json.Unmarshal(body, &result)

		fmt.Printf("Snip published: https://snipbox.vercel.app/clips/%s\n", result["slug"])
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringP("password", "p", "", "Password for private snip")
	publishCmd.Flags().StringP("slug", "s", "", "Custom slug for the snip")
	publishCmd.Flags().StringP("expires", "e", "1h", "Expiration time (e.g., 1h, 1d, 7d)")
	publishCmd.Flags().BoolP("destroy-on-view", "d", false, "Destroy snip after first view")
}
