package xbow

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rsclarke/xbow/internal/api"
)

func TestAssetFromGetResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	resp := &api.GetAPIV1AssetsAssetIDResponse{
		ID:                   "asset-123",
		Name:                 "Test Asset",
		OrganizationID:       "org-456",
		Lifecycle:            api.Active,
		Sku:                  "standard-sku",
		StartURL:             "https://example.com",
		MaxRequestsPerSecond: 100,
		CreatedAt:            now,
		UpdatedAt:            now.Add(time.Hour),
	}

	got := assetFromGetResponse(resp)

	if got.ID != "asset-123" {
		t.Errorf("ID = %q, want 'asset-123'", got.ID)
	}
	if got.Name != "Test Asset" {
		t.Errorf("Name = %q, want 'Test Asset'", got.Name)
	}
	if got.OrganizationID != "org-456" {
		t.Errorf("OrganizationID = %q, want 'org-456'", got.OrganizationID)
	}
	if got.Lifecycle != AssetLifecycleActive {
		t.Errorf("Lifecycle = %q, want %q", got.Lifecycle, AssetLifecycleActive)
	}
	if got.Sku != "standard-sku" {
		t.Errorf("Sku = %q, want 'standard-sku'", got.Sku)
	}
	if got.StartURL != "https://example.com" {
		t.Errorf("StartURL = %q, want 'https://example.com'", got.StartURL)
	}
	if got.MaxRequestsPerSecond != 100 {
		t.Errorf("MaxRequestsPerSecond = %d, want 100", got.MaxRequestsPerSecond)
	}
	if !got.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, now)
	}
	if !got.UpdatedAt.Equal(now.Add(time.Hour)) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, now.Add(time.Hour))
	}
}

func TestAssetFromPutResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	resp := &api.PutAPIV1AssetsAssetIDResponse{
		ID:                   "asset-123",
		Name:                 "Updated Asset",
		OrganizationID:       "org-456",
		Lifecycle:            api.PutAPIV1AssetsAssetIDResponseLifecycleArchived,
		Sku:                  "premium-sku",
		StartURL:             "https://updated.example.com",
		MaxRequestsPerSecond: 500,
		CreatedAt:            now,
		UpdatedAt:            now.Add(time.Hour),
	}

	got := assetFromPutResponse(resp)

	if got.ID != "asset-123" {
		t.Errorf("ID = %q, want 'asset-123'", got.ID)
	}
	if got.Name != "Updated Asset" {
		t.Errorf("Name = %q, want 'Updated Asset'", got.Name)
	}
	if got.Lifecycle != AssetLifecycleArchived {
		t.Errorf("Lifecycle = %q, want %q", got.Lifecycle, AssetLifecycleArchived)
	}
	if got.Sku != "premium-sku" {
		t.Errorf("Sku = %q, want 'premium-sku'", got.Sku)
	}
	if got.MaxRequestsPerSecond != 500 {
		t.Errorf("MaxRequestsPerSecond = %d, want 500", got.MaxRequestsPerSecond)
	}
}

func TestAssetFromCreateResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	resp := &api.PostAPIV1OrganizationsOrganizationIDAssetsResponse{
		ID:                   "asset-new",
		Name:                 "New Asset",
		OrganizationID:       "org-789",
		Lifecycle:            api.PostAPIV1OrganizationsOrganizationIDAssetsResponseLifecycleActive,
		Sku:                  "basic-sku",
		StartURL:             "",
		MaxRequestsPerSecond: 0,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	got := assetFromCreateResponse(resp)

	if got.ID != "asset-new" {
		t.Errorf("ID = %q, want 'asset-new'", got.ID)
	}
	if got.Name != "New Asset" {
		t.Errorf("Name = %q, want 'New Asset'", got.Name)
	}
	if got.OrganizationID != "org-789" {
		t.Errorf("OrganizationID = %q, want 'org-789'", got.OrganizationID)
	}
	if got.Lifecycle != AssetLifecycleActive {
		t.Errorf("Lifecycle = %q, want %q", got.Lifecycle, AssetLifecycleActive)
	}
	if got.StartURL != "" {
		t.Errorf("StartURL = %q, want empty", got.StartURL)
	}
}

func TestAssetsPageFromResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	t.Run("converts items and pagination info", func(t *testing.T) {
		nextCursor := "cursor-xyz"
		resp := &api.GetAPIV1OrganizationsOrganizationIDAssetsResponse{
			Items: api.GetAPIV1OrganizationsOrganizationIDAssets_Response_Items{
				{
					ID:        "a1",
					Name:      "Asset 1",
					Lifecycle: api.GetAPIV1OrganizationsOrganizationIDAssetsResponseItemsLifecycleActive,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "a2",
					Name:      "Asset 2",
					Lifecycle: api.GetAPIV1OrganizationsOrganizationIDAssetsResponseItemsLifecycleArchived,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			NextCursor: &nextCursor,
		}

		got := assetsPageFromResponse(resp)

		if len(got.Items) != 2 {
			t.Fatalf("got %d items, want 2", len(got.Items))
		}
		if got.Items[0].ID != "a1" {
			t.Errorf("Items[0].ID = %q, want 'a1'", got.Items[0].ID)
		}
		if got.Items[0].Lifecycle != AssetLifecycleActive {
			t.Errorf("Items[0].Lifecycle = %q, want %q", got.Items[0].Lifecycle, AssetLifecycleActive)
		}
		if got.Items[1].Lifecycle != AssetLifecycleArchived {
			t.Errorf("Items[1].Lifecycle = %q, want %q", got.Items[1].Lifecycle, AssetLifecycleArchived)
		}
		if got.PageInfo.NextCursor != "cursor-xyz" {
			t.Errorf("NextCursor = %q, want 'cursor-xyz'", got.PageInfo.NextCursor)
		}
		if !got.PageInfo.HasMore {
			t.Error("HasMore = false, want true")
		}
	})

	t.Run("handles nil cursor", func(t *testing.T) {
		resp := &api.GetAPIV1OrganizationsOrganizationIDAssetsResponse{
			Items:      api.GetAPIV1OrganizationsOrganizationIDAssets_Response_Items{},
			NextCursor: nil,
		}

		got := assetsPageFromResponse(resp)

		if got.PageInfo.NextCursor != "" {
			t.Errorf("NextCursor = %q, want empty", got.PageInfo.NextCursor)
		}
		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false")
		}
	})

	t.Run("handles empty cursor string", func(t *testing.T) {
		empty := ""
		resp := &api.GetAPIV1OrganizationsOrganizationIDAssetsResponse{
			Items:      api.GetAPIV1OrganizationsOrganizationIDAssets_Response_Items{},
			NextCursor: &empty,
		}

		got := assetsPageFromResponse(resp)

		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false for empty cursor")
		}
	})
}

func TestUpdateAssetNilRequest(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.Assets.Update(context.TODO(), "asset-123", nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.Code != "ERR_INVALID_REQUEST" {
		t.Errorf("Code = %q, want 'ERR_INVALID_REQUEST'", apiErr.Code)
	}
}

func TestCreateAssetNilRequest(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.Assets.Create(context.TODO(), "org-123", nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.Code != "ERR_INVALID_REQUEST" {
		t.Errorf("Code = %q, want 'ERR_INVALID_REQUEST'", apiErr.Code)
	}
}
