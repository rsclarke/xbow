package cmd

import (
	"fmt"
	"os"

	"github.com/rsclarke/xbow"
	"github.com/spf13/cobra"
)

var (
	orgKey         string
	integrationKey string
	outputFormat   string
)

var rootCmd = &cobra.Command{
	Use:   "xbow",
	Short: "XBOW CLI - Interact with the XBOW API",
	Long:  `A command-line interface for interacting with the XBOW security assessment platform.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&orgKey, "org-key", "", "Organization API key (or set XBOW_ORG_KEY env var)")
	rootCmd.PersistentFlags().StringVar(&integrationKey, "integration-key", "", "Integration API key (or set XBOW_INTEGRATION_KEY env var)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json")
}

func newClient() (*xbow.Client, error) {
	opts := []xbow.ClientOption{}

	key := orgKey
	if key == "" {
		key = os.Getenv("XBOW_ORG_KEY")
	}
	if key != "" {
		opts = append(opts, xbow.WithOrganizationKey(key))
	}

	intKey := integrationKey
	if intKey == "" {
		intKey = os.Getenv("XBOW_INTEGRATION_KEY")
	}
	if intKey != "" {
		opts = append(opts, xbow.WithIntegrationKey(intKey))
	}

	if key == "" && intKey == "" {
		return nil, fmt.Errorf("API key required: use --org-key/--integration-key or set XBOW_ORG_KEY/XBOW_INTEGRATION_KEY")
	}

	return xbow.NewClient(opts...)
}
