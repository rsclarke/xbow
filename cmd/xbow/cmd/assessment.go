package cmd

import (
	"context"
	"fmt"
	"iter"

	"github.com/rsclarke/xbow"
	"github.com/spf13/cobra"
)

var assessmentCmd = &cobra.Command{
	Use:     "assessment",
	Aliases: []string{"assessments"},
	Short:   "Manage assessments",
}

func init() {
	rootCmd.AddCommand(assessmentCmd)
	assessmentCmd.AddCommand(assessmentGetCmd)
	assessmentCmd.AddCommand(assessmentCreateCmd)
	assessmentCmd.AddCommand(assessmentListCmd)
	assessmentCmd.AddCommand(assessmentCancelCmd)
	assessmentCmd.AddCommand(assessmentPauseCmd)
	assessmentCmd.AddCommand(assessmentResumeCmd)
}

var assessmentGetCmd = &cobra.Command{
	Use:   "get <assessment-id>",
	Short: "Get an assessment by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		assessment, err := client.Assessments.Get(context.Background(), args[0])
		if err != nil {
			return err
		}

		return printAssessment(assessment)
	},
}

var (
	createAssetID       string
	createAttackCredits int64
	createObjective     string
)

var assessmentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new assessment",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		req := &xbow.CreateAssessmentRequest{
			AttackCredits: createAttackCredits,
		}
		if createObjective != "" {
			req.Objective = &createObjective
		}

		assessment, err := client.Assessments.Create(context.Background(), createAssetID, req)
		if err != nil {
			return err
		}

		return printAssessment(assessment)
	},
}

func init() {
	assessmentCreateCmd.Flags().StringVar(&createAssetID, "asset-id", "", "Asset ID to create assessment for (required)")
	assessmentCreateCmd.Flags().Int64Var(&createAttackCredits, "attack-credits", 0, "Number of attack credits to use (required)")
	assessmentCreateCmd.Flags().StringVar(&createObjective, "objective", "", "Assessment objective")
	_ = assessmentCreateCmd.MarkFlagRequired("asset-id")
	_ = assessmentCreateCmd.MarkFlagRequired("attack-credits")
}

var (
	listAssetID string
	listLimit   int
)

var assessmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List assessments for an asset",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		var opts *xbow.ListOptions
		if listLimit > 0 {
			opts = &xbow.ListOptions{Limit: listLimit}
		}

		return printAssessmentList(client.Assessments.AllByAsset(context.Background(), listAssetID, opts))
	},
}

func init() {
	assessmentListCmd.Flags().StringVar(&listAssetID, "asset-id", "", "Asset ID to list assessments for (required)")
	assessmentListCmd.Flags().IntVar(&listLimit, "limit", 0, "Maximum number of results per page")
	_ = assessmentListCmd.MarkFlagRequired("asset-id")
}

var assessmentCancelCmd = &cobra.Command{
	Use:   "cancel <assessment-id>",
	Short: "Cancel a running assessment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		assessment, err := client.Assessments.Cancel(context.Background(), args[0])
		if err != nil {
			return err
		}

		return printAssessment(assessment)
	},
}

var assessmentPauseCmd = &cobra.Command{
	Use:   "pause <assessment-id>",
	Short: "Pause a running assessment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		assessment, err := client.Assessments.Pause(context.Background(), args[0])
		if err != nil {
			return err
		}

		return printAssessment(assessment)
	},
}

var assessmentResumeCmd = &cobra.Command{
	Use:   "resume <assessment-id>",
	Short: "Resume a paused assessment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		assessment, err := client.Assessments.Resume(context.Background(), args[0])
		if err != nil {
			return err
		}

		return printAssessment(assessment)
	},
}

func printAssessment(a *xbow.Assessment) error {
	if outputFormat == "json" {
		return printJSON(a)
	}

	w := newTabWriter()
	printRow(w, "ID:", a.ID)
	printRow(w, "NAME:", a.Name)
	printRow(w, "ASSET ID:", a.AssetID)
	printRow(w, "STATE:", a.State)
	printRow(w, "PROGRESS:", fmt.Sprintf("%.1f%%", a.Progress*100))
	printRow(w, "ATTACK CREDITS:", a.AttackCredits)
	printRow(w, "CREATED:", a.CreatedAt.Format("2006-01-02 15:04:05"))
	printRow(w, "UPDATED:", a.UpdatedAt.Format("2006-01-02 15:04:05"))
	return w.Flush()
}

func printAssessmentList(iter iter.Seq2[xbow.AssessmentListItem, error]) error {
	if outputFormat == "json" {
		var items []xbow.AssessmentListItem
		for a, err := range iter {
			if err != nil {
				return err
			}
			items = append(items, a)
		}
		return printJSON(items)
	}

	w := newTabWriter()
	printRow(w, "ID", "NAME", "STATE", "PROGRESS", "CREATED")
	for a, err := range iter {
		if err != nil {
			return err
		}
		printRow(w, a.ID, a.Name, a.State, fmt.Sprintf("%.1f%%", a.Progress*100), a.CreatedAt.Format("2006-01-02"))
	}
	return w.Flush()
}
