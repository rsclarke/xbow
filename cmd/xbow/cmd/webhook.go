package cmd

import (
	"context"
	"fmt"
	"iter"
	"strings"

	"github.com/rsclarke/xbow"
	"github.com/spf13/cobra"
)

var webhookCmd = &cobra.Command{
	Use:     "webhook",
	Aliases: []string{"webhooks"},
	Short:   "Manage webhooks",
}

func init() {
	rootCmd.AddCommand(webhookCmd)
	webhookCmd.AddCommand(webhookGetCmd)
	webhookCmd.AddCommand(webhookCreateCmd)
	webhookCmd.AddCommand(webhookUpdateCmd)
	webhookCmd.AddCommand(webhookDeleteCmd)
	webhookCmd.AddCommand(webhookPingCmd)
	webhookCmd.AddCommand(webhookListCmd)
	webhookCmd.AddCommand(webhookDeliveriesCmd)
}

// get

var webhookGetCmd = &cobra.Command{
	Use:   "get <webhook-id>",
	Short: "Get a webhook by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		webhook, err := client.Webhooks.Get(context.Background(), args[0])
		if err != nil {
			return err
		}

		return printWebhook(webhook)
	},
}

// create

var (
	webhookCreateOrgID      string
	webhookCreateTargetURL  string
	webhookCreateAPIVersion string
	webhookCreateEvents     []string
)

var webhookCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new webhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		events := make([]xbow.WebhookEventType, 0, len(webhookCreateEvents))
		for _, e := range webhookCreateEvents {
			events = append(events, xbow.WebhookEventType(e))
		}

		webhook, err := client.Webhooks.Create(context.Background(), webhookCreateOrgID, &xbow.CreateWebhookRequest{
			APIVersion: xbow.WebhookAPIVersion(webhookCreateAPIVersion),
			TargetURL:  webhookCreateTargetURL,
			Events:     events,
		})
		if err != nil {
			return err
		}

		return printWebhook(webhook)
	},
}

func init() {
	webhookCreateCmd.Flags().StringVar(&webhookCreateOrgID, "org-id", "", "Organization ID (required)")
	webhookCreateCmd.Flags().StringVar(&webhookCreateTargetURL, "target-url", "", "Webhook target URL (required)")
	webhookCreateCmd.Flags().StringVar(&webhookCreateAPIVersion, "api-version", "2026-02-01", "Webhook API version")
	webhookCreateCmd.Flags().StringArrayVar(&webhookCreateEvents, "event", nil, `Event type to subscribe to (repeatable, e.g. "assessment.changed")`)
	_ = webhookCreateCmd.MarkFlagRequired("org-id")
	_ = webhookCreateCmd.MarkFlagRequired("target-url")
	_ = webhookCreateCmd.MarkFlagRequired("event")
}

// update

var (
	webhookUpdateTargetURL  string
	webhookUpdateAPIVersion string
	webhookUpdateEvents     []string
)

var webhookUpdateCmd = &cobra.Command{
	Use:   "update <webhook-id>",
	Short: "Update a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		req := &xbow.UpdateWebhookRequest{}
		if cmd.Flags().Changed("target-url") {
			req.TargetURL = &webhookUpdateTargetURL
		}
		if cmd.Flags().Changed("api-version") {
			v := xbow.WebhookAPIVersion(webhookUpdateAPIVersion)
			req.APIVersion = &v
		}
		if cmd.Flags().Changed("event") {
			events := make([]xbow.WebhookEventType, 0, len(webhookUpdateEvents))
			for _, e := range webhookUpdateEvents {
				events = append(events, xbow.WebhookEventType(e))
			}
			req.Events = events
		}

		webhook, err := client.Webhooks.Update(context.Background(), args[0], req)
		if err != nil {
			return err
		}

		return printWebhook(webhook)
	},
}

func init() {
	webhookUpdateCmd.Flags().StringVar(&webhookUpdateTargetURL, "target-url", "", "Webhook target URL")
	webhookUpdateCmd.Flags().StringVar(&webhookUpdateAPIVersion, "api-version", "", "Webhook API version")
	webhookUpdateCmd.Flags().StringArrayVar(&webhookUpdateEvents, "event", nil, `Event type to subscribe to (repeatable)`)
}

// delete

var webhookDeleteCmd = &cobra.Command{
	Use:   "delete <webhook-id>",
	Short: "Delete a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		if err := client.Webhooks.Delete(context.Background(), args[0]); err != nil {
			return err
		}

		fmt.Println("Webhook deleted.")
		return nil
	},
}

// ping

var webhookPingCmd = &cobra.Command{
	Use:   "ping <webhook-id>",
	Short: "Send a ping event to a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		if err := client.Webhooks.Ping(context.Background(), args[0]); err != nil {
			return err
		}

		fmt.Println("Ping sent.")
		return nil
	},
}

// list

var (
	webhookListOrgID string
	webhookListLimit int
)

var webhookListCmd = &cobra.Command{
	Use:   "list",
	Short: "List webhooks for an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		var opts *xbow.ListOptions
		if webhookListLimit > 0 {
			opts = &xbow.ListOptions{Limit: webhookListLimit}
		}

		return printWebhookList(client.Webhooks.AllByOrganization(context.Background(), webhookListOrgID, opts))
	},
}

func init() {
	webhookListCmd.Flags().StringVar(&webhookListOrgID, "org-id", "", "Organization ID (required)")
	webhookListCmd.Flags().IntVar(&webhookListLimit, "limit", 0, "Maximum number of results per page")
	_ = webhookListCmd.MarkFlagRequired("org-id")
}

// deliveries

var webhookDeliveriesLimit int

var webhookDeliveriesCmd = &cobra.Command{
	Use:   "deliveries <webhook-id>",
	Short: "List deliveries for a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		var opts *xbow.ListOptions
		if webhookDeliveriesLimit > 0 {
			opts = &xbow.ListOptions{Limit: webhookDeliveriesLimit}
		}

		return printDeliveryList(client.Webhooks.AllDeliveries(context.Background(), args[0], opts))
	},
}

func init() {
	webhookDeliveriesCmd.Flags().IntVar(&webhookDeliveriesLimit, "limit", 0, "Maximum number of results per page")
}

// output helpers

func printWebhook(wh *xbow.Webhook) error {
	if outputFormat == "json" {
		return printJSON(wh)
	}

	w := newTabWriter()
	printRow(w, "ID:", wh.ID)
	printRow(w, "TARGET URL:", wh.TargetURL)
	printRow(w, "API VERSION:", wh.APIVersion)
	printRow(w, "EVENTS:", strings.Join(webhookEventStrings(wh.Events), ", "))
	printRow(w, "CREATED:", wh.CreatedAt.Format("2006-01-02 15:04:05"))
	printRow(w, "UPDATED:", wh.UpdatedAt.Format("2006-01-02 15:04:05"))
	return w.Flush()
}

func printWebhookList(iter iter.Seq2[xbow.WebhookListItem, error]) error {
	if outputFormat == "json" {
		var items []xbow.WebhookListItem
		for wh, err := range iter {
			if err != nil {
				return err
			}
			items = append(items, wh)
		}
		return printJSON(items)
	}

	w := newTabWriter()
	printRow(w, "ID", "TARGET URL", "API VERSION", "EVENTS", "CREATED")
	for wh, err := range iter {
		if err != nil {
			return err
		}
		printRow(w, wh.ID, wh.TargetURL, wh.APIVersion, strings.Join(webhookEventStrings(wh.Events), ","), wh.CreatedAt.Format("2006-01-02"))
	}
	return w.Flush()
}

func printDeliveryList(iter iter.Seq2[xbow.WebhookDelivery, error]) error {
	if outputFormat == "json" {
		var items []xbow.WebhookDelivery
		for d, err := range iter {
			if err != nil {
				return err
			}
			items = append(items, d)
		}
		return printJSON(items)
	}

	w := newTabWriter()
	printRow(w, "SENT AT", "SUCCESS", "STATUS")
	for d, err := range iter {
		if err != nil {
			return err
		}
		printRow(w, d.SentAt.Format("2006-01-02 15:04:05"), fmt.Sprintf("%v", d.Success), d.Response.Status)
	}
	return w.Flush()
}

func webhookEventStrings(events []xbow.WebhookEventType) []string {
	s := make([]string, 0, len(events))
	for _, e := range events {
		s = append(s, string(e))
	}
	return s
}
