package xbow

import (
	"context"
	"net/http"

	"github.com/rsclarke/xbow/internal/api"
)

// MetaService handles meta-related API calls.
type MetaService struct {
	client *Client
}

// WebhookSigningKey represents a public key used to verify webhook signatures.
type WebhookSigningKey struct {
	// PublicKey is a Base64-encoded Ed25519 public key in SPKI format.
	PublicKey string `json:"publicKey"`
}

// GetOpenAPISpec retrieves the OpenAPI specification for the current API version.
// The response is returned as raw JSON bytes since the schema is dynamic.
func (s *MetaService) GetOpenAPISpec(ctx context.Context) ([]byte, error) {
	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	_, body, err := s.client.do(ctx, http.MethodGet, "/api/v1/meta/openapi.json", auth)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetWebhookSigningKeys retrieves the public keys used to sign webhook requests.
// Use these keys to verify webhook signatures. The array supports key rotation -
// during rotation, multiple keys may be active.
func (s *MetaService) GetWebhookSigningKeys(ctx context.Context) ([]WebhookSigningKey, error) {
	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	opts := &api.GetAPIV1MetaWebhooksSigningKeysRequestOptions{
		Header: &api.GetAPIV1MetaWebhooksSigningKeysHeaders{
			XXBOWAPIVersion: api.GetAPIV1MetaWebhooksSigningKeysHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.GetAPIV1MetaWebhooksSigningKeys(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return webhookSigningKeysFromResponse(resp), nil
}

// webhookSigningKeysFromResponse converts the generated response to domain types.
func webhookSigningKeysFromResponse(r *api.GetAPIV1MetaWebhooksSigningKeysResponse) []WebhookSigningKey {
	if r == nil {
		return nil
	}

	keys := make([]WebhookSigningKey, 0, len(*r))
	for _, item := range *r {
		keys = append(keys, WebhookSigningKey{
			PublicKey: item.PublicKey,
		})
	}
	return keys
}
