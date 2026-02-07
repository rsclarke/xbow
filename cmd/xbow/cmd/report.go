package cmd

import (
	"context"
	"fmt"
	"iter"
	"os"

	"github.com/rsclarke/xbow"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:     "report",
	Aliases: []string{"reports"},
	Short:   "Manage reports",
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.AddCommand(reportGetCmd)
	reportCmd.AddCommand(reportSummaryCmd)
	reportCmd.AddCommand(reportListCmd)
}

// get

var reportGetOutputFile string

var reportGetCmd = &cobra.Command{
	Use:   "get <report-id>",
	Short: "Download a report as PDF",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		data, err := client.Reports.Get(context.Background(), args[0])
		if err != nil {
			return err
		}

		if reportGetOutputFile != "" {
			return os.WriteFile(reportGetOutputFile, data, 0o644) //nolint:gosec // PDF output file; 0644 is intentional
		}

		_, err = os.Stdout.Write(data)
		return err
	},
}

func init() {
	reportGetCmd.Flags().StringVarP(&reportGetOutputFile, "output-file", "f", "", "Path to write the PDF file")
}

// summary

var reportSummaryOutputFile string

var reportSummaryCmd = &cobra.Command{
	Use:   "summary <report-id>",
	Short: "Get the markdown summary of a report",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		summary, err := client.Reports.GetSummary(context.Background(), args[0])
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(summary)
		}

		if reportSummaryOutputFile != "" {
			return os.WriteFile(reportSummaryOutputFile, []byte(summary.Markdown), 0o644) //nolint:gosec // markdown output file; 0644 is intentional
		}

		fmt.Print(summary.Markdown)
		return nil
	},
}

func init() {
	reportSummaryCmd.Flags().StringVarP(&reportSummaryOutputFile, "output-file", "f", "", "Path to write the markdown summary")
}

// list

var (
	reportListAssetID string
	reportListLimit   int
)

var reportListCmd = &cobra.Command{
	Use:   "list",
	Short: "List reports for an asset",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		var opts *xbow.ListOptions
		if reportListLimit > 0 {
			opts = &xbow.ListOptions{Limit: reportListLimit}
		}

		return printReportList(client.Reports.AllByAsset(context.Background(), reportListAssetID, opts))
	},
}

func init() {
	reportListCmd.Flags().StringVar(&reportListAssetID, "asset-id", "", "Asset ID to list reports for (required)")
	reportListCmd.Flags().IntVar(&reportListLimit, "limit", 0, "Maximum number of results per page")
	_ = reportListCmd.MarkFlagRequired("asset-id")
}

// output helpers

func printReportList(iter iter.Seq2[xbow.ReportListItem, error]) error {
	if outputFormat == "json" {
		var items []xbow.ReportListItem
		for r, err := range iter {
			if err != nil {
				return err
			}
			items = append(items, r)
		}
		return printJSON(items)
	}

	w := newTabWriter()
	printRow(w, "ID", "VERSION", "CREATED")
	for r, err := range iter {
		if err != nil {
			return err
		}
		printRow(w, r.ID, r.Version, r.CreatedAt.Format("2006-01-02"))
	}
	return w.Flush()
}
