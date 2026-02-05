package xbow

import (
	"testing"

	"github.com/rsclarke/xbow/internal/api"
)

func TestWebhookSigningKeysFromResponse(t *testing.T) {
	t.Run("converts multiple keys", func(t *testing.T) {
		resp := api.GetAPIV1MetaWebhooksSigningKeysResponse{
			{PublicKey: "MCowBQYDK2VwAyEA...key1"},
			{PublicKey: "MCowBQYDK2VwAyEA...key2"},
		}

		got := webhookSigningKeysFromResponse(&resp)

		if len(got) != 2 {
			t.Fatalf("got %d keys, want 2", len(got))
		}
		if got[0].PublicKey != "MCowBQYDK2VwAyEA...key1" {
			t.Errorf("got[0].PublicKey = %q, want 'MCowBQYDK2VwAyEA...key1'", got[0].PublicKey)
		}
		if got[1].PublicKey != "MCowBQYDK2VwAyEA...key2" {
			t.Errorf("got[1].PublicKey = %q, want 'MCowBQYDK2VwAyEA...key2'", got[1].PublicKey)
		}
	})

	t.Run("handles empty response", func(t *testing.T) {
		resp := api.GetAPIV1MetaWebhooksSigningKeysResponse{}

		got := webhookSigningKeysFromResponse(&resp)

		if len(got) != 0 {
			t.Errorf("got %d keys, want 0", len(got))
		}
	})

	t.Run("handles nil response", func(t *testing.T) {
		got := webhookSigningKeysFromResponse(nil)

		if got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})

	t.Run("handles single key for rotation", func(t *testing.T) {
		resp := api.GetAPIV1MetaWebhooksSigningKeysResponse{
			{PublicKey: "MCowBQYDK2VwAyEAsinglekey"},
		}

		got := webhookSigningKeysFromResponse(&resp)

		if len(got) != 1 {
			t.Fatalf("got %d keys, want 1", len(got))
		}
		if got[0].PublicKey != "MCowBQYDK2VwAyEAsinglekey" {
			t.Errorf("got[0].PublicKey = %q, want 'MCowBQYDK2VwAyEAsinglekey'", got[0].PublicKey)
		}
	})
}

func TestMetaServiceInitialized(t *testing.T) {
	client, err := NewClient("test-key")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if client.Meta == nil {
		t.Error("client.Meta is nil, expected initialized MetaService")
	}
}
