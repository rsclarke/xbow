package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/rsclarke/xbow"
	"github.com/spf13/cobra"
)

var assetCmd = &cobra.Command{
	Use:     "asset",
	Aliases: []string{"assets"},
	Short:   "Manage assets",
}

func init() {
	rootCmd.AddCommand(assetCmd)
	assetCmd.AddCommand(assetGetCmd)
	assetCmd.AddCommand(assetCreateCmd)
	assetCmd.AddCommand(assetListCmd)
	assetCmd.AddCommand(assetUpdateCmd)
}

// get

var assetGetCmd = &cobra.Command{
	Use:   "get <asset-id>",
	Short: "Get an asset by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		asset, err := client.Assets.Get(context.Background(), args[0])
		if err != nil {
			return err
		}

		return printAsset(asset)
	},
}

// create

var (
	assetCreateOrgID string
	assetCreateName  string
	assetCreateSku   string
)

var assetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new asset",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		asset, err := client.Assets.Create(context.Background(), assetCreateOrgID, &xbow.CreateAssetRequest{
			Name: assetCreateName,
			Sku:  assetCreateSku,
		})
		if err != nil {
			return err
		}

		return printAsset(asset)
	},
}

func init() {
	assetCreateCmd.Flags().StringVar(&assetCreateOrgID, "org-id", "", "Organization ID (required)")
	assetCreateCmd.Flags().StringVar(&assetCreateName, "name", "", "Asset name (required)")
	assetCreateCmd.Flags().StringVar(&assetCreateSku, "sku", "standard-sku", "Asset SKU")
	_ = assetCreateCmd.MarkFlagRequired("org-id")
	_ = assetCreateCmd.MarkFlagRequired("name")
}

// list

var (
	assetListOrgID string
	assetListLimit int
)

var assetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List assets for an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		var opts *xbow.ListOptions
		if assetListLimit > 0 {
			opts = &xbow.ListOptions{Limit: assetListLimit}
		}

		return printAssetList(client.Assets.AllByOrganization(context.Background(), assetListOrgID, opts))
	},
}

func init() {
	assetListCmd.Flags().StringVar(&assetListOrgID, "org-id", "", "Organization ID (required)")
	assetListCmd.Flags().IntVar(&assetListLimit, "limit", 0, "Maximum number of results per page")
	_ = assetListCmd.MarkFlagRequired("org-id")
}

// update

var (
	assetUpdateName        string
	assetUpdateStartURL    string
	assetUpdateMaxRPS      int
	assetUpdateSku         string
	assetUpdateHeaders     []string
	assetUpdateFromFile    string
	assetUpdateCredentials []string
	assetUpdateDNSRules    []string
	assetUpdateHTTPRules   []string
)

var assetUpdateCmd = &cobra.Command{
	Use:   "update <asset-id>",
	Short: "Update an asset",
	Long: `Update an asset by ID. Fetches the current asset, applies changes, and saves.

Simple fields:
  --name, --start-url, --max-rps, --sku

Repeatable structured fields:
  --header "Key: Value"
  --credential "name=n,type=basic,username=u,password=p"
  --dns-rule "action=allow-attack,type=hostname,filter=example.com"
  --http-rule "action=deny,type=url,filter=https://evil.com"

  Optional sub-fields for --credential: email-address, authenticator-uri
  Optional sub-fields for --dns-rule/--http-rule: id, include-subdomains

Full replacement from JSON file:
  --from-file asset.json   (use - for stdin)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ctx := context.Background()

		var req *xbow.UpdateAssetRequest

		if assetUpdateFromFile != "" {
			req, err = loadUpdateRequestFromFile(assetUpdateFromFile)
			if err != nil {
				return err
			}
		} else {
			current, err := client.Assets.Get(ctx, args[0])
			if err != nil {
				return fmt.Errorf("fetching current asset: %w", err)
			}
			req = updateRequestFromAsset(current)

			if cmd.Flags().Changed("name") {
				req.Name = assetUpdateName
			}
			if cmd.Flags().Changed("start-url") {
				req.StartURL = assetUpdateStartURL
			}
			if cmd.Flags().Changed("max-rps") {
				req.MaxRequestsPerSecond = assetUpdateMaxRPS
			}
			if cmd.Flags().Changed("sku") {
				req.Sku = &assetUpdateSku
			}
			if cmd.Flags().Changed("header") {
				headers, err := parseHeaders(assetUpdateHeaders)
				if err != nil {
					return err
				}
				req.Headers = headers
			}
			if cmd.Flags().Changed("credential") {
				creds, err := parseCredentials(assetUpdateCredentials)
				if err != nil {
					return err
				}
				req.Credentials = creds
			}
			if cmd.Flags().Changed("dns-rule") {
				rules, err := parseDNSRules(assetUpdateDNSRules)
				if err != nil {
					return err
				}
				req.DNSBoundaryRules = rules
			}
			if cmd.Flags().Changed("http-rule") {
				rules, err := parseHTTPRules(assetUpdateHTTPRules)
				if err != nil {
					return err
				}
				req.HTTPBoundaryRules = rules
			}
		}

		asset, err := client.Assets.Update(ctx, args[0], req)
		if err != nil {
			return err
		}

		return printAsset(asset)
	},
}

func init() {
	assetUpdateCmd.Flags().StringVar(&assetUpdateName, "name", "", "Asset name")
	assetUpdateCmd.Flags().StringVar(&assetUpdateStartURL, "start-url", "", "Start URL")
	assetUpdateCmd.Flags().IntVar(&assetUpdateMaxRPS, "max-rps", 0, "Max requests per second")
	assetUpdateCmd.Flags().StringVar(&assetUpdateSku, "sku", "", "Asset SKU")
	assetUpdateCmd.Flags().StringArrayVar(&assetUpdateHeaders, "header", nil, `Header in "Key: Value" format (repeatable)`)
	assetUpdateCmd.Flags().StringArrayVar(&assetUpdateCredentials, "credential", nil, `Credential as "name=n,type=basic,username=u,password=p" (repeatable)`)
	assetUpdateCmd.Flags().StringArrayVar(&assetUpdateDNSRules, "dns-rule", nil, `DNS boundary rule as "action=allow-attack,type=hostname,filter=example.com" (repeatable)`)
	assetUpdateCmd.Flags().StringArrayVar(&assetUpdateHTTPRules, "http-rule", nil, `HTTP boundary rule as "action=deny,type=url,filter=https://example.com" (repeatable)`)
	assetUpdateCmd.Flags().StringVar(&assetUpdateFromFile, "from-file", "", "Load full update request from JSON file (- for stdin)")
}

// updateRequestFromAsset builds an UpdateAssetRequest from the current asset state.
func updateRequestFromAsset(a *xbow.Asset) *xbow.UpdateAssetRequest {
	req := &xbow.UpdateAssetRequest{
		Name:              a.Name,
		Sku:               &a.Sku,
		Credentials:       a.Credentials,
		DNSBoundaryRules:  a.DNSBoundaryRules,
		Headers:           a.Headers,
		HTTPBoundaryRules: a.HTTPBoundaryRules,
	}
	if a.StartURL != nil {
		req.StartURL = *a.StartURL
	}
	if a.MaxRequestsPerSecond != nil {
		req.MaxRequestsPerSecond = *a.MaxRequestsPerSecond
	}
	if a.ApprovedTimeWindows != nil {
		req.ApprovedTimeWindows = a.ApprovedTimeWindows
	}
	return req
}

func loadUpdateRequestFromFile(path string) (*xbow.UpdateAssetRequest, error) {
	var data []byte
	var err error

	if path == "-" {
		data, err = os.ReadFile("/dev/stdin")
	} else {
		data, err = os.ReadFile(filepath.Clean(path))
	}
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var req xbow.UpdateAssetRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return &req, nil
}

// parseHeaders parses "Key: Value" strings into a header map.
// Multiple values for the same key are accumulated.
func parseHeaders(raw []string) (map[string][]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	headers := make(map[string][]string)
	for _, h := range raw {
		key, value, ok := strings.Cut(h, ":")
		if !ok {
			return nil, fmt.Errorf("invalid header format %q, expected \"Key: Value\"", h)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			return nil, fmt.Errorf("empty header key in %q", h)
		}
		headers[key] = append(headers[key], value)
	}
	return headers, nil
}

// parseKV parses a "key1=val1,key2=val2" string into a map.
// Values may contain commas if the key=value pair contains an equals sign
// that isn't part of a subsequent key. This simple parser splits on commas
// and then on the first equals sign.
func parseKV(s string) map[string]string {
	m := make(map[string]string)
	for _, part := range strings.Split(s, ",") {
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		m[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return m
}

func parseCredentials(raw []string) ([]xbow.Credential, error) {
	creds := make([]xbow.Credential, 0, len(raw))
	for _, s := range raw {
		kv := parseKV(s)

		name := kv["name"]
		if name == "" {
			return nil, fmt.Errorf("credential missing required field 'name' in %q", s)
		}
		typ := kv["type"]
		if typ == "" {
			return nil, fmt.Errorf("credential missing required field 'type' in %q", s)
		}
		username := kv["username"]
		if username == "" {
			return nil, fmt.Errorf("credential missing required field 'username' in %q", s)
		}
		password := kv["password"]
		if password == "" {
			return nil, fmt.Errorf("credential missing required field 'password' in %q", s)
		}

		cred := xbow.Credential{
			ID:       kv["id"],
			Name:     name,
			Type:     typ,
			Username: username,
			Password: password,
		}
		if v, ok := kv["email-address"]; ok {
			cred.EmailAddress = &v
		}
		if v, ok := kv["authenticator-uri"]; ok {
			cred.AuthenticatorURI = &v
		}
		creds = append(creds, cred)
	}
	return creds, nil
}

func parseDNSRules(raw []string) ([]xbow.DNSBoundaryRule, error) {
	rules := make([]xbow.DNSBoundaryRule, 0, len(raw))
	for _, s := range raw {
		kv := parseKV(s)

		action := kv["action"]
		if action == "" {
			return nil, fmt.Errorf("dns-rule missing required field 'action' in %q", s)
		}
		typ := kv["type"]
		if typ == "" {
			return nil, fmt.Errorf("dns-rule missing required field 'type' in %q", s)
		}
		filter := kv["filter"]
		if filter == "" {
			return nil, fmt.Errorf("dns-rule missing required field 'filter' in %q", s)
		}

		rule := xbow.DNSBoundaryRule{
			ID:     kv["id"],
			Action: xbow.DNSBoundaryRuleAction(action),
			Type:   typ,
			Filter: filter,
		}
		if v, ok := kv["include-subdomains"]; ok {
			b := v == "true"
			rule.IncludeSubdomains = &b
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func parseHTTPRules(raw []string) ([]xbow.HTTPBoundaryRule, error) {
	rules := make([]xbow.HTTPBoundaryRule, 0, len(raw))
	for _, s := range raw {
		kv := parseKV(s)

		action := kv["action"]
		if action == "" {
			return nil, fmt.Errorf("http-rule missing required field 'action' in %q", s)
		}
		typ := kv["type"]
		if typ == "" {
			return nil, fmt.Errorf("http-rule missing required field 'type' in %q", s)
		}
		filter := kv["filter"]
		if filter == "" {
			return nil, fmt.Errorf("http-rule missing required field 'filter' in %q", s)
		}

		rule := xbow.HTTPBoundaryRule{
			ID:     kv["id"],
			Action: xbow.HTTPBoundaryRuleAction(action),
			Type:   typ,
			Filter: filter,
		}
		if v, ok := kv["include-subdomains"]; ok {
			b := v == "true"
			rule.IncludeSubdomains = &b
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// output helpers

func printAsset(a *xbow.Asset) error {
	if outputFormat == "json" {
		return printJSON(a)
	}

	w := newTabWriter()
	printRow(w, "ID:", a.ID)
	printRow(w, "NAME:", a.Name)
	printRow(w, "ORGANIZATION ID:", a.OrganizationID)
	printRow(w, "LIFECYCLE:", a.Lifecycle)
	printRow(w, "SKU:", a.Sku)
	if a.StartURL != nil {
		printRow(w, "START URL:", *a.StartURL)
	}
	if a.MaxRequestsPerSecond != nil {
		printRow(w, "MAX RPS:", *a.MaxRequestsPerSecond)
	}

	if len(a.Credentials) > 0 {
		printRow(w, "CREDENTIALS:", fmt.Sprintf("%d configured", len(a.Credentials)))
	}
	if len(a.DNSBoundaryRules) > 0 {
		printRow(w, "DNS RULES:", fmt.Sprintf("%d configured", len(a.DNSBoundaryRules)))
	}
	if len(a.HTTPBoundaryRules) > 0 {
		printRow(w, "HTTP RULES:", fmt.Sprintf("%d configured", len(a.HTTPBoundaryRules)))
	}
	if len(a.Headers) > 0 {
		printRow(w, "HEADERS:", fmt.Sprintf("%d configured", len(a.Headers)))
	}

	if a.Checks != nil {
		printRow(w, "CHECK REACHABLE:", a.Checks.AssetReachable.State)
		printRow(w, "CHECK CREDENTIALS:", a.Checks.Credentials.State)
		printRow(w, "CHECK DNS RULES:", a.Checks.DNSBoundaryRules.State)
	}

	printRow(w, "CREATED:", a.CreatedAt.Format("2006-01-02 15:04:05"))
	printRow(w, "UPDATED:", a.UpdatedAt.Format("2006-01-02 15:04:05"))
	return w.Flush()
}

func printAssetList(iter iter.Seq2[xbow.AssetListItem, error]) error {
	if outputFormat == "json" {
		var items []xbow.AssetListItem
		for a, err := range iter {
			if err != nil {
				return err
			}
			items = append(items, a)
		}
		return printJSON(items)
	}

	w := newTabWriter()
	printRow(w, "ID", "NAME", "LIFECYCLE", "CREATED")
	for a, err := range iter {
		if err != nil {
			return err
		}
		printRow(w, a.ID, a.Name, a.Lifecycle, a.CreatedAt.Format("2006-01-02"))
	}
	return w.Flush()
}
