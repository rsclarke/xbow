# xbow

> **⚠️ Unofficial & Alpha**: This is **not** an official XBOW library. It is a personal/community project and is currently in **alpha**. The API may change without notice.

An idiomatic Go client library for the [XBOW API](https://docs.xbow.com/api/), following [google/go-github](https://github.com/google/go-github) patterns.

## Installation

```bash
go get github.com/rsclarke/xbow
```

Requires Go 1.23+ (uses `iter.Seq2` for pagination).

## Usage

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

    // Create a new assessment
    assessment, err = client.Assessments.Create(ctx, "asset-id", &xbow.CreateAssessmentRequest{
        AttackCredits: 100,
    })
    if err != nil {
        log.Fatal(err)
    }

    // List all assessments for an asset with automatic pagination
    for assessment, err := range client.Assessments.AllByAsset(ctx, "asset-id", nil) {
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("- %s: %s\n", assessment.Name, assessment.State)
    }

    // Get an asset
    asset, err := client.Assets.Get(ctx, "asset-id")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Asset: %s (%s)\n", asset.Name, asset.Lifecycle)

    // Create a new asset
    asset, err = client.Assets.Create(ctx, "organization-id", &xbow.CreateAssetRequest{
        Name: "My Web App",
        Sku:  "standard-sku",
    })
    if err != nil {
        log.Fatal(err)
    }

    // List all assets in an organization
    for asset, err := range client.Assets.AllByOrganization(ctx, "organization-id", nil) {
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("- %s: %s\n", asset.Name, asset.Lifecycle)
    }
}
```

## API Coverage

| Service           | Operation                              | Implemented |
|-------------------|----------------------------------------|:-----------:|
| **Assessments**   | Get                                    | ✅          |
|                   | Create                                 | ✅          |
|                   | ListByAsset / AllByAsset               | ✅          |
|                   | Cancel                                 | ✅          |
|                   | Pause                                  | ✅          |
|                   | Resume                                 | ✅          |
| **Assets**        | Get                                    | ✅          |
|                   | Update                                 | ✅          |
|                   | ListByOrganization / AllByOrganization | ✅          |
|                   | Create                                 | ✅          |
| **Findings**      | Get                                    | ✅          |
|                   | VerifyFix                              | ✅          |
|                   | ListByAsset / AllByAsset               | ✅          |
| **Meta**          | GetOpenAPISpec                         | ✅          |
|                   | GetWebhookSigningKeys                  | ✅          |
| **Organizations** | Get                                    | ✅          |
|                   | Update                                 | ✅          |
|                   | ListByIntegration / AllByIntegration   | ✅          |
|                   | Create                                 | ✅          |
|                   | CreateKey                              | ✅          |
|                   | RevokeKey                              | ✅          |
| **Reports**       | Get                                    | ✅          |
|                   | GetSummary                             | ✅          |
|                   | ListByAsset / AllByAsset               | ✅          |
| **Webhooks**      | Get                                    | ✅          |
|                   | Update                                 | ✅          |
|                   | Delete                                 | ✅          |
|                   | Ping                                   | ✅          |
|                   | ListByOrganization / AllByOrganization | ✅          |
|                   | Create                                 | ✅          |
|                   | ListDeliveries / AllDeliveries         | ✅          |

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
