# xbow

> **⚠️ Unofficial & Alpha**: This is **not** an official XBOW library. It is a personal/community project and is currently in **alpha**. The API may change without notice.

An idiomatic Go client library and CLI for the [XBOW API](https://docs.xbow.com/api/), following [google/go-github](https://github.com/google/go-github) patterns.

## Installation

### CLI

```bash
go install github.com/rsclarke/xbow/cmd/xbow@latest
```

### Library

```bash
go get github.com/rsclarke/xbow
```

Requires Go 1.23+ (uses `iter.Seq2` for pagination).

## CLI Usage

The `xbow` CLI provides command-line access to the XBOW API.

### Authentication

Set your API key via environment variable or flag:

```bash
# Organization key (for most operations)
export XBOW_ORG_KEY="your-org-key"

# Or pass directly
xbow --org-key "your-org-key" assessment list --asset-id abc123
```

### Assets

```bash
# Create an asset
xbow asset create --org-id <org-id> --name "My App" --sku standard-sku

# Get an asset
xbow asset get <asset-id>

# List all assets for an organization
xbow asset list --org-id <org-id>

# Update simple fields (GET-then-PUT; unspecified fields are preserved)
xbow asset update <asset-id> --name "New Name" --start-url "https://example.com" --max-rps 10

# Update with repeatable structured flags
# NOTE: Repeatable flags (--header, --credential, --dns-rule, --http-rule)
# perform a full replacement of that field. If you specify any values for a
# flag, they replace ALL existing values — omitted entries are removed.
# To keep existing entries, re-specify them alongside any new ones.
xbow asset update <asset-id> \
  --header "X-Custom: value" \
  --credential "name=admin,type=basic,username=u,password=p" \
  --dns-rule "action=allow-attack,type=hostname,filter=example.com,include-subdomains=true" \
  --http-rule "action=deny,type=url,filter=https://evil.com"

# Full replacement from a JSON file (or - for stdin)
xbow asset update <asset-id> --from-file asset.json
```

### Assessments

```bash
# Create an assessment
xbow assessment create --asset-id <asset-id> --attack-credits 100 --objective "Find vulnerabilities"

# Get an assessment
xbow assessment get <assessment-id>

# List all assessments for an asset
xbow assessment list --asset-id <asset-id>

# Control assessment execution
xbow assessment pause <assessment-id>
xbow assessment resume <assessment-id>
xbow assessment cancel <assessment-id>
```

### Findings

```bash
# Get a finding
xbow finding get <finding-id>

# List all findings for an asset
xbow finding list --asset-id <asset-id>

# Verify that a finding has been fixed (triggers a targeted assessment)
xbow finding verify-fix <finding-id>
```

### Reports

```bash
# Download a report as PDF
xbow report get <report-id> --output-file report.pdf

# Or pipe to stdout
xbow report get <report-id> > report.pdf

# Get the markdown summary
xbow report summary <report-id>

# Save the summary to a file
xbow report summary <report-id> --output-file summary.md

# List all reports for an asset
xbow report list --asset-id <asset-id>
```

### Organizations

Organization management requires an integration key (`--integration-key` or `XBOW_INTEGRATION_KEY`).
The `get` subcommand also accepts an organization key.

```bash
# Get an organization
xbow organization get <org-id>

# List organizations for an integration
xbow organization list --integration-id <integration-id>

# Create an organization with members
xbow organization create \
  --integration-id <integration-id> \
  --name "My Org" \
  --external-id "ext-123" \
  --member "email=alice@example.com,name=Alice" \
  --member "email=bob@example.com,name=Bob"

# Update an organization
xbow organization update <org-id> --name "New Name" --external-id "ext-456"

# Create an API key for an organization
xbow organization create-key <org-id> --name "CI Key" --expires-in-days 90

# Revoke an API key
xbow organization revoke-key <key-id>
```

### Webhooks

```bash
# List webhooks for an organization
xbow webhook list --org-id <org-id>

# Get a webhook
xbow webhook get <webhook-id>

# Create a webhook
xbow webhook create --org-id <org-id> --target-url "https://example.com/hook" \
  --event "assessment.changed" --event "finding.changed"

# Update a webhook
xbow webhook update <webhook-id> --target-url "https://example.com/new-hook"

# Delete a webhook
xbow webhook delete <webhook-id>

# Ping a webhook (send a test event)
xbow webhook ping <webhook-id>

# List deliveries for a webhook (table shows summary; use --output json for full payloads)
xbow webhook deliveries <webhook-id>
```

### Meta

```bash
# Get the OpenAPI specification
xbow meta openapi

# Save the OpenAPI specification to a file
xbow meta openapi --output-file openapi.json

# Get webhook signing keys
xbow meta signing-keys
```

### Output Formats

```bash
# Table output (default)
xbow assessment get <id>

# JSON output
xbow assessment get <id> --output json
```

### Global Flags

| Flag | Environment Variable | Description |
|------|---------------------|-------------|
| `--org-key` | `XBOW_ORG_KEY` | Organization API key |
| `--integration-key` | `XBOW_INTEGRATION_KEY` | Integration API key |
| `--output`, `-o` | - | Output format: `table` (default), `json` |
| `--version` | - | Print CLI and API version |

## Library Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/rsclarke/xbow"
)

func main() {
    // Most endpoints use an organization key
    client, err := xbow.NewClient(xbow.WithOrganizationKey("your-org-key"))
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get an assessment
    assessment, err := client.Assessments.Get(ctx, "assessment-id")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Assessment: %s (%s)\n", assessment.Name, assessment.State)

    // List all findings for the assessment's asset
    for finding, err := range client.Findings.AllByAsset(ctx, assessment.AssetID, nil) {
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("- %s: %s (%s)\n", finding.Name, finding.Severity, finding.State)
    }
}
```

### Webhook Verification

Verify incoming webhook requests using Ed25519 signatures. Fetch the signing keys from the API, then create a verifier:

```go
keys, err := client.Meta.GetWebhookSigningKeys(ctx)
if err != nil {
    log.Fatal(err)
}

verifier, err := xbow.NewWebhookVerifier(keys)
if err != nil {
    log.Fatal(err)
}
```

Use it as HTTP middleware, which returns `401 Unauthorized` for invalid signatures:

```go
http.Handle("/webhook", verifier.Middleware(myHandler))
```

Or verify requests manually:

```go
if err := verifier.Verify(r); err != nil {
    // handle invalid signature
}
```

Options can be passed to `NewWebhookVerifier` to adjust clock skew tolerance (default 5 minutes) and maximum body size (default 5 MB):

```go
verifier, err := xbow.NewWebhookVerifier(keys,
    xbow.WithMaxClockSkew(10*time.Minute),
    xbow.WithMaxBodyBytes(10*1024*1024), // 10 MB
)
```

## Authentication

The XBOW API uses two types of API keys:

- **Organization key** - Used for most endpoints (assessments, assets, findings, etc.)
- **Integration key** - Required for organization management endpoints

```go
// Most users - organization key only
client, _ := xbow.NewClient(xbow.WithOrganizationKey("your-org-key"))

// Integration key only (for organization management)
client, _ := xbow.NewClient(xbow.WithIntegrationKey("your-integration-key"))

// Both keys for full access
client, _ := xbow.NewClient(
    xbow.WithOrganizationKey("your-org-key"),
    xbow.WithIntegrationKey("your-integration-key"),
)
```

## Configuration

```go
// Custom base URL
client, _ := xbow.NewClient(
    xbow.WithOrganizationKey("your-org-key"),
    xbow.WithBaseURL("https://custom.xbow.com"),
)

// Custom HTTP client
client, _ := xbow.NewClient(
    xbow.WithOrganizationKey("your-org-key"),
    xbow.WithHTTPClient(myHTTPClient),
)
```

## Rate Limiting

The API may return `429 Too Many Requests` responses. You can configure a rate limiter to automatically throttle requests:

```go
import "golang.org/x/time/rate"

// 10 requests per second with burst of 10
limiter := rate.NewLimiter(rate.Every(time.Second), 10)

client, _ := xbow.NewClient(
    xbow.WithOrganizationKey("your-org-key"),
    xbow.WithRateLimiter(limiter),
)
```

The `RateLimiter` interface requires only a `Wait(context.Context) error` method, so you can provide any custom implementation.

## Pagination

List methods return a single page. Use `All*` methods for automatic pagination:

```go
// Single page
page, err := client.Assessments.ListByAsset(ctx, assetID, &xbow.ListOptions{Limit: 10})

// All pages (iterator)
for assessment, err := range client.Assessments.AllByAsset(ctx, assetID, nil) {
    // ...
}
```

## Error Handling

Errors from the API are returned as `*xbow.Error` with structured error codes:

```go
assessment, err := client.Assessments.Get(ctx, "invalid-id")
if err != nil {
    var apiErr *xbow.Error
    if errors.As(err, &apiErr) {
        fmt.Printf("API error: %s - %s\n", apiErr.Code, apiErr.Message)
    }
}
```

## License

MIT
