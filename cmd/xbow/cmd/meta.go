package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "API metadata and utilities",
}

func init() {
	rootCmd.AddCommand(metaCmd)
	metaCmd.AddCommand(metaOpenapiCmd)
	metaCmd.AddCommand(metaSigningKeysCmd)
}

// openapi

var metaOpenapiOutputFile string

var metaOpenapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "Get the OpenAPI specification",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		data, err := client.Meta.GetOpenAPISpec(context.Background())
		if err != nil {
			return err
		}

		if metaOpenapiOutputFile != "" {
			if err := os.WriteFile(filepath.Clean(metaOpenapiOutputFile), data, 0o644); err != nil {
				return fmt.Errorf("writing file: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Written to %s\n", metaOpenapiOutputFile)
			return nil
		}

		_, err = os.Stdout.Write(data)
		return err
	},
}

func init() {
	metaOpenapiCmd.Flags().StringVarP(&metaOpenapiOutputFile, "output-file", "f", "", "Write output to file")
}

// signing-keys

var metaSigningKeysCmd = &cobra.Command{
	Use:   "signing-keys",
	Short: "Get webhook signing keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		keys, err := client.Meta.GetWebhookSigningKeys(context.Background())
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(keys)
		}

		w := newTabWriter()
		printRow(w, "PUBLIC KEY")
		for _, k := range keys {
			printRow(w, k.PublicKey)
		}
		return w.Flush()
	},
}
