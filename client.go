// Package xbow provides an idiomatic Go client for the XBOW API.
//
// This package wraps the generated doordash oapi-codegen client with a
// more ergonomic, google/go-github-style API.
//
// Example usage:
//
//	// Most endpoints use an organization key
//	client, err := xbow.NewClient(xbow.WithOrganizationKey("your-org-key"))
//
//	// Organization management endpoints require an integration key
//	client, err := xbow.NewClient(xbow.WithIntegrationKey("your-integration-key"))
//
//	// Use both keys for full access
//	client, err := xbow.NewClient(
//	    xbow.WithOrganizationKey("your-org-key"),
//	    xbow.WithIntegrationKey("your-integration-key"),
//	)
//
//	// Get an assessment
//	assessment, err := client.Assessments.Get(ctx, "assessment-id")
//
//	// List all assessments for an asset with automatic pagination
//	for assessment, err := range client.Assessments.AllByAsset(ctx, assetID, nil) {
//	    if err != nil {
//	        return err
//	    }
//	    fmt.Println(assessment.Name)
//	}
package xbow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/doordash-oss/oapi-codegen-dd/v3/pkg/runtime"
	"github.com/rsclarke/xbow/internal/api"
)

// API configuration constants.
const (
	DefaultBaseURL = "https://console.xbow.com"
	APIVersion     = "2026-02-01"
)

// Client manages communication with the XBOW API.
type Client struct {
	raw            *api.Client
	orgKey         string
	integrationKey string
	baseURL        string
	httpClient     *http.Client

	// Services
	Assessments   *AssessmentsService
	Assets        *AssetsService
	Findings      *FindingsService
	Meta          *MetaService
	Organizations *OrganizationsService
	Reports       *ReportsService
	Webhooks      *WebhooksService
}

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*clientConfig)

type clientConfig struct {
	baseURL        string
	httpClient     *http.Client
	apiClientOpts  []runtime.APIClientOption
	orgKey         string
	integrationKey string
}

// WithBaseURL sets a custom base URL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *clientConfig) {
		c.baseURL = baseURL
	}
}

// httpClientWrapper adapts an *http.Client to the runtime.HttpRequestDoer interface.
type httpClientWrapper struct {
	client *http.Client
}

func (w *httpClientWrapper) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return w.client.Do(req.WithContext(ctx))
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *clientConfig) {
		c.httpClient = httpClient
		c.apiClientOpts = append(c.apiClientOpts, runtime.WithHTTPClient(&httpClientWrapper{client: httpClient}))
	}
}

// WithAPIClientOption adds a runtime.APIClientOption to the underlying client.
func WithAPIClientOption(opt runtime.APIClientOption) ClientOption {
	return func(c *clientConfig) {
		c.apiClientOpts = append(c.apiClientOpts, opt)
	}
}

// WithOrganizationKey sets the organization API key for authenticating with most endpoints.
func WithOrganizationKey(key string) ClientOption {
	return func(c *clientConfig) {
		c.orgKey = key
	}
}

// WithIntegrationKey sets the integration API key for authenticating with organization management endpoints.
func WithIntegrationKey(key string) ClientOption {
	return func(c *clientConfig) {
		c.integrationKey = key
	}
}

// NewClient creates a new XBOW API client.
func NewClient(opts ...ClientOption) (*Client, error) {
	cfg := &clientConfig{
		baseURL:    DefaultBaseURL,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	raw, err := api.NewDefaultClient(cfg.baseURL, cfg.apiClientOpts...)
	if err != nil {
		return nil, err
	}

	c := &Client{
		raw:            raw,
		orgKey:         cfg.orgKey,
		integrationKey: cfg.integrationKey,
		baseURL:        cfg.baseURL,
		httpClient:     cfg.httpClient,
	}

	c.Assessments = &AssessmentsService{client: c}
	c.Assets = &AssetsService{client: c}
	c.Findings = &FindingsService{client: c}
	c.Meta = &MetaService{client: c}
	c.Organizations = &OrganizationsService{client: c}
	c.Reports = &ReportsService{client: c}
	c.Webhooks = &WebhooksService{client: c}

	return c, nil
}

// Raw returns the underlying generated client for advanced use cases.
func (c *Client) Raw() *api.Client {
	return c.raw
}

// authEditorFor returns a request editor that adds authentication headers for the given key.
func (c *Client) authEditorFor(key string) runtime.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+key)
		return nil
	}
}

// orgAuthEditor returns a request editor using the organization key.
// Returns an error if the organization key is not set.
func (c *Client) orgAuthEditor() (runtime.RequestEditorFn, error) {
	if c.orgKey == "" {
		return nil, fmt.Errorf("xbow: organization key is required; provide xbow.WithOrganizationKey(...)")
	}
	return c.authEditorFor(c.orgKey), nil
}

// integrationAuthEditor returns a request editor using the integration key.
// Returns an error if the integration key is not set.
func (c *Client) integrationAuthEditor() (runtime.RequestEditorFn, error) {
	if c.integrationKey == "" {
		return nil, fmt.Errorf("xbow: integration key is required; provide xbow.WithIntegrationKey(...)")
	}
	return c.authEditorFor(c.integrationKey), nil
}

// orgOrIntegrationAuthEditor returns a request editor preferring integration key, falling back to org key.
// Returns an error if neither key is set.
func (c *Client) orgOrIntegrationAuthEditor() (runtime.RequestEditorFn, error) {
	if c.integrationKey != "" {
		return c.authEditorFor(c.integrationKey), nil
	}
	if c.orgKey != "" {
		return c.authEditorFor(c.orgKey), nil
	}
	return nil, fmt.Errorf("xbow: organization key or integration key is required; provide xbow.WithOrganizationKey(...) or xbow.WithIntegrationKey(...)")
}
