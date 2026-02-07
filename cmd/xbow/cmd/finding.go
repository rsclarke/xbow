package cmd

import (
	"context"
	"iter"

	"github.com/rsclarke/xbow"
	"github.com/spf13/cobra"
)

var findingCmd = &cobra.Command{
	Use:     "finding",
	Aliases: []string{"findings"},
	Short:   "Manage findings",
}

func init() {
	rootCmd.AddCommand(findingCmd)
	findingCmd.AddCommand(findingGetCmd)
	findingCmd.AddCommand(findingListCmd)
	findingCmd.AddCommand(findingVerifyFixCmd)
}

// get

var findingGetCmd = &cobra.Command{
	Use:   "get <finding-id>",
	Short: "Get a finding by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		finding, err := client.Findings.Get(context.Background(), args[0])
		if err != nil {
			return err
		}

		return printFinding(finding)
	},
}

// list

var (
	findingListAssetID string
	findingListLimit   int
)

var findingListCmd = &cobra.Command{
	Use:   "list",
	Short: "List findings for an asset",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		var opts *xbow.ListOptions
		if findingListLimit > 0 {
			opts = &xbow.ListOptions{Limit: findingListLimit}
		}

		return printFindingList(client.Findings.AllByAsset(context.Background(), findingListAssetID, opts))
	},
}

func init() {
	findingListCmd.Flags().StringVar(&findingListAssetID, "asset-id", "", "Asset ID to list findings for (required)")
	findingListCmd.Flags().IntVar(&findingListLimit, "limit", 0, "Maximum number of results per page")
	_ = findingListCmd.MarkFlagRequired("asset-id")
}

// verify-fix

var findingVerifyFixCmd = &cobra.Command{
	Use:   "verify-fix <finding-id>",
	Short: "Verify that a finding has been fixed",
	Long:  `Triggers a targeted assessment to verify that a finding's vulnerability has been mitigated.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		assessment, err := client.Findings.VerifyFix(context.Background(), args[0])
		if err != nil {
			return err
		}

		return printAssessment(assessment)
	},
}

// output helpers

func printFinding(f *xbow.Finding) error {
	if outputFormat == "json" {
		return printJSON(f)
	}

	w := newTabWriter()
	printRow(w, "ID:", f.ID)
	printRow(w, "NAME:", f.Name)
	printRow(w, "SEVERITY:", f.Severity)
	printRow(w, "STATE:", f.State)
	printRow(w, "SUMMARY:", f.Summary)
	printRow(w, "IMPACT:", f.Impact)
	printRow(w, "MITIGATIONS:", f.Mitigations)
	printRow(w, "RECIPE:", f.Recipe)
	printRow(w, "EVIDENCE:", f.Evidence)
	printRow(w, "CREATED:", f.CreatedAt.Format("2006-01-02 15:04:05"))
	printRow(w, "UPDATED:", f.UpdatedAt.Format("2006-01-02 15:04:05"))
	return w.Flush()
}

func printFindingList(iter iter.Seq2[xbow.FindingListItem, error]) error {
	if outputFormat == "json" {
		var items []xbow.FindingListItem
		for f, err := range iter {
			if err != nil {
				return err
			}
			items = append(items, f)
		}
		return printJSON(items)
	}

	w := newTabWriter()
	printRow(w, "ID", "NAME", "SEVERITY", "STATE", "CREATED")
	for f, err := range iter {
		if err != nil {
			return err
		}
		printRow(w, f.ID, f.Name, f.Severity, f.State, f.CreatedAt.Format("2006-01-02"))
	}
	return w.Flush()
}
