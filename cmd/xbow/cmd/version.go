package cmd

import (
	"fmt"

	"github.com/rsclarke/xbow"
	"github.com/spf13/cobra"
)

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the xbow CLI and API version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("xbow version %s\n", version)
		fmt.Printf("api version %s\n", xbow.APIVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
