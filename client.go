// Package xbow provides an idiomatic Go client for the XBOW API.
//
// This package wraps the generated doordash oapi-codegen client with a
// more ergonomic, google/go-github-style API.
//
// Example usage:
//
//	client, err := xbow.NewClient("your-api-key")
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
	raw        *api.Client
	apiKey     string
	baseURL    string
	httpClient *http.Client

	// Services
	Assessments   *AssessmentsService
	Assets        *AssetsService
	Findings      *FindingsService
	Meta          *MetaService
	Organizations *OrganizationsService
}

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*clientConfig)

type clientConfig struct {
	baseURL       string
	httpClient    *http.Client
	apiClientOpts []runtime.APIClientOption
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

// NewClient creates a new XBOW API client with the given API key.
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
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
		raw:        raw,
		apiKey:     apiKey,
		baseURL:    cfg.baseURL,
		httpClient: cfg.httpClient,
	}

	c.Assessments = &AssessmentsService{client: c}
	c.Assets = &AssetsService{client: c}
	c.Findings = &FindingsService{client: c}
	c.Meta = &MetaService{client: c}
	c.Organizations = &OrganizationsService{client: c}

	return c, nil
}

// Raw returns the underlying generated client for advanced use cases.
func (c *Client) Raw() *api.Client {
	return c.raw
}

// authEditor returns a request editor that adds authentication headers.
func (c *Client) authEditor() runtime.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		return nil
	}
}
