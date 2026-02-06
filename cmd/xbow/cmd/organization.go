package cmd

import (
	"context"
	"fmt"
	"iter"
	"strings"

	"github.com/rsclarke/xbow"
	"github.com/spf13/cobra"
)

var organizationCmd = &cobra.Command{
	Use:     "organization",
	Aliases: []string{"organizations"},
	Short:   "Manage organizations",
}

func init() {
	rootCmd.AddCommand(organizationCmd)
	organizationCmd.AddCommand(orgGetCmd)
	organizationCmd.AddCommand(orgCreateCmd)
	organizationCmd.AddCommand(orgUpdateCmd)
	organizationCmd.AddCommand(orgListCmd)
	organizationCmd.AddCommand(orgCreateKeyCmd)
	organizationCmd.AddCommand(orgRevokeKeyCmd)
}

// get

var orgGetCmd = &cobra.Command{
	Use:   "get <org-id>",
	Short: "Get an organization by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		org, err := client.Organizations.Get(context.Background(), args[0])
		if err != nil {
			return err
		}

		return printOrganization(org)
	},
}

// create

var (
	orgCreateIntegrationID string
	orgCreateName          string
	orgCreateExternalID    string
	orgCreateMembers       []string
)

var orgCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new organization",
	Long: `Create a new organization in an integration.

Members are specified as repeatable flags:
  --member "email=alice@example.com,name=Alice"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		members, err := parseMembers(orgCreateMembers)
		if err != nil {
			return err
		}

		req := &xbow.CreateOrganizationRequest{
			Name:    orgCreateName,
			Members: members,
		}
		if cmd.Flags().Changed("external-id") {
			req.ExternalID = &orgCreateExternalID
		}

		org, err := client.Organizations.Create(context.Background(), orgCreateIntegrationID, req)
		if err != nil {
			return err
		}

		return printOrganization(org)
	},
}

func init() {
	orgCreateCmd.Flags().StringVar(&orgCreateIntegrationID, "integration-id", "", "Integration ID (required)")
	orgCreateCmd.Flags().StringVar(&orgCreateName, "name", "", "Organization name (required)")
	orgCreateCmd.Flags().StringVar(&orgCreateExternalID, "external-id", "", "External ID")
	orgCreateCmd.Flags().StringArrayVar(&orgCreateMembers, "member", nil, `Member as "email=alice@example.com,name=Alice" (repeatable, at least one required)`)
	_ = orgCreateCmd.MarkFlagRequired("integration-id")
	_ = orgCreateCmd.MarkFlagRequired("name")
	_ = orgCreateCmd.MarkFlagRequired("member")
}

// update

var (
	orgUpdateName       string
	orgUpdateExternalID string
)

var orgUpdateCmd = &cobra.Command{
	Use:   "update <org-id>",
	Short: "Update an organization",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		req := &xbow.UpdateOrganizationRequest{
			Name: orgUpdateName,
		}
		if cmd.Flags().Changed("external-id") {
			req.ExternalID = &orgUpdateExternalID
		}

		org, err := client.Organizations.Update(context.Background(), args[0], req)
		if err != nil {
			return err
		}

		return printOrganization(org)
	},
}

func init() {
	orgUpdateCmd.Flags().StringVar(&orgUpdateName, "name", "", "Organization name (required)")
	orgUpdateCmd.Flags().StringVar(&orgUpdateExternalID, "external-id", "", "External ID (use empty string to clear)")
	_ = orgUpdateCmd.MarkFlagRequired("name")
}

// list

var (
	orgListIntegrationID string
	orgListLimit         int
)

var orgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List organizations for an integration",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		var opts *xbow.ListOptions
		if orgListLimit > 0 {
			opts = &xbow.ListOptions{Limit: orgListLimit}
		}

		return printOrganizationList(client.Organizations.AllByIntegration(context.Background(), orgListIntegrationID, opts))
	},
}

func init() {
	orgListCmd.Flags().StringVar(&orgListIntegrationID, "integration-id", "", "Integration ID (required)")
	orgListCmd.Flags().IntVar(&orgListLimit, "limit", 0, "Maximum number of results per page")
	_ = orgListCmd.MarkFlagRequired("integration-id")
}

// create-key

var (
	orgCreateKeyName        string
	orgCreateKeyExpiresDays int
)

var orgCreateKeyCmd = &cobra.Command{
	Use:   "create-key <org-id>",
	Short: "Create an API key for an organization",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		req := &xbow.CreateKeyRequest{
			Name: orgCreateKeyName,
		}
		if cmd.Flags().Changed("expires-in-days") {
			req.ExpiresInDays = &orgCreateKeyExpiresDays
		}

		key, err := client.Organizations.CreateKey(context.Background(), args[0], req)
		if err != nil {
			return err
		}

		return printAPIKey(key)
	},
}

func init() {
	orgCreateKeyCmd.Flags().StringVar(&orgCreateKeyName, "name", "", "Key name (required)")
	orgCreateKeyCmd.Flags().IntVar(&orgCreateKeyExpiresDays, "expires-in-days", 0, "Number of days until key expires")
	_ = orgCreateKeyCmd.MarkFlagRequired("name")
}

// revoke-key

var orgRevokeKeyCmd = &cobra.Command{
	Use:   "revoke-key <key-id>",
	Short: "Revoke an organization API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		if err := client.Organizations.RevokeKey(context.Background(), args[0]); err != nil {
			return err
		}

		fmt.Println("Key revoked successfully.")
		return nil
	},
}

// parsing helpers

func parseMembers(raw []string) ([]xbow.OrganizationMember, error) {
	members := make([]xbow.OrganizationMember, 0, len(raw))
	for _, s := range raw {
		kv := parseKV(s)

		email := kv["email"]
		if email == "" {
			return nil, fmt.Errorf("member missing required field 'email' in %q", s)
		}
		name := kv["name"]
		if name == "" {
			return nil, fmt.Errorf("member missing required field 'name' in %q", s)
		}

		members = append(members, xbow.OrganizationMember{
			Email: email,
			Name:  name,
		})
	}
	return members, nil
}

// output helpers

func printOrganization(o *xbow.Organization) error {
	if outputFormat == "json" {
		return printJSON(o)
	}

	w := newTabWriter()
	printRow(w, "ID:", o.ID)
	printRow(w, "NAME:", o.Name)
	if o.ExternalID != nil {
		printRow(w, "EXTERNAL ID:", *o.ExternalID)
	}
	printRow(w, "STATE:", strings.ToUpper(string(o.State)))
	printRow(w, "CREATED:", o.CreatedAt.Format("2006-01-02 15:04:05"))
	printRow(w, "UPDATED:", o.UpdatedAt.Format("2006-01-02 15:04:05"))
	return w.Flush()
}

func printOrganizationList(iter iter.Seq2[xbow.OrganizationListItem, error]) error {
	if outputFormat == "json" {
		var items []xbow.OrganizationListItem
		for o, err := range iter {
			if err != nil {
				return err
			}
			items = append(items, o)
		}
		return printJSON(items)
	}

	w := newTabWriter()
	printRow(w, "ID", "NAME", "STATE", "CREATED")
	for o, err := range iter {
		if err != nil {
			return err
		}
		printRow(w, o.ID, o.Name, strings.ToUpper(string(o.State)), o.CreatedAt.Format("2006-01-02"))
	}
	return w.Flush()
}

func printAPIKey(k *xbow.OrganizationAPIKey) error {
	if outputFormat == "json" {
		return printJSON(k)
	}

	w := newTabWriter()
	printRow(w, "ID:", k.ID)
	printRow(w, "NAME:", k.Name)
	printRow(w, "KEY:", k.Key)
	if k.ExpiresAt != nil {
		printRow(w, "EXPIRES AT:", k.ExpiresAt.Format("2006-01-02 15:04:05"))
	}
	printRow(w, "CREATED:", k.CreatedAt.Format("2006-01-02 15:04:05"))
	return w.Flush()
}
