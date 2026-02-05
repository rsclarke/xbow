package xbow

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rsclarke/xbow/internal/api"
)

func TestOrganizationFromGetResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	t.Run("converts all fields", func(t *testing.T) {
		resp := &api.GetAPIV1OrganizationsOrganizationIDResponse{
			ID:         "org-123",
			Name:       "Test Organization",
			ExternalID: "ext-456",
			State:      api.GetAPIV1OrganizationsOrganizationIDResponseStateActive,
			CreatedAt:  now,
			UpdatedAt:  now.Add(time.Hour),
		}

		got := organizationFromGetResponse(resp)

		if got.ID != "org-123" {
			t.Errorf("ID = %q, want 'org-123'", got.ID)
		}
		if got.Name != "Test Organization" {
			t.Errorf("Name = %q, want 'Test Organization'", got.Name)
		}
		if got.ExternalID == nil || *got.ExternalID != "ext-456" {
			t.Errorf("ExternalID = %v, want 'ext-456'", got.ExternalID)
		}
		if got.State != OrganizationStateActive {
			t.Errorf("State = %q, want %q", got.State, OrganizationStateActive)
		}
		if !got.CreatedAt.Equal(now) {
			t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, now)
		}
		if !got.UpdatedAt.Equal(now.Add(time.Hour)) {
			t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, now.Add(time.Hour))
		}
	})

	t.Run("handles empty external ID as nil", func(t *testing.T) {
		resp := &api.GetAPIV1OrganizationsOrganizationIDResponse{
			ID:         "org-123",
			Name:       "Test Org",
			ExternalID: "",
			State:      api.GetAPIV1OrganizationsOrganizationIDResponseStateDisabled,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		got := organizationFromGetResponse(resp)

		if got.ExternalID != nil {
			t.Errorf("ExternalID = %v, want nil", got.ExternalID)
		}
		if got.State != OrganizationStateDisabled {
			t.Errorf("State = %q, want %q", got.State, OrganizationStateDisabled)
		}
	})
}

func TestOrganizationFromPutResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	resp := &api.PutAPIV1OrganizationsOrganizationIDResponse{
		ID:         "org-456",
		Name:       "Updated Organization",
		ExternalID: "updated-ext",
		State:      api.PutAPIV1OrganizationsOrganizationIDResponseStateActive,
		CreatedAt:  now,
		UpdatedAt:  now.Add(time.Hour),
	}

	got := organizationFromPutResponse(resp)

	if got.ID != "org-456" {
		t.Errorf("ID = %q, want 'org-456'", got.ID)
	}
	if got.Name != "Updated Organization" {
		t.Errorf("Name = %q, want 'Updated Organization'", got.Name)
	}
	if got.ExternalID == nil || *got.ExternalID != "updated-ext" {
		t.Errorf("ExternalID = %v, want 'updated-ext'", got.ExternalID)
	}
	if got.State != OrganizationStateActive {
		t.Errorf("State = %q, want %q", got.State, OrganizationStateActive)
	}
}

func TestOrganizationFromCreateResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	resp := &api.PostAPIV1IntegrationsIntegrationIDOrganizationsResponse{
		ID:         "org-new",
		Name:       "New Organization",
		ExternalID: "new-ext",
		State:      api.PostAPIV1IntegrationsIntegrationIDOrganizationsResponseStateActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	got := organizationFromCreateResponse(resp)

	if got.ID != "org-new" {
		t.Errorf("ID = %q, want 'org-new'", got.ID)
	}
	if got.Name != "New Organization" {
		t.Errorf("Name = %q, want 'New Organization'", got.Name)
	}
	if got.ExternalID == nil || *got.ExternalID != "new-ext" {
		t.Errorf("ExternalID = %v, want 'new-ext'", got.ExternalID)
	}
	if got.State != OrganizationStateActive {
		t.Errorf("State = %q, want %q", got.State, OrganizationStateActive)
	}
}

func TestOrganizationsPageFromResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	t.Run("converts items and pagination info", func(t *testing.T) {
		nextCursor := "cursor-abc"
		resp := &api.GetAPIV1IntegrationsIntegrationIDOrganizationsResponse{
			Items: api.GetAPIV1IntegrationsIntegrationIDOrganizations_Response_Items{
				{
					ID:         "org-1",
					Name:       "Org 1",
					ExternalID: "ext-1",
					State:      api.GetAPIV1IntegrationsIntegrationIDOrganizationsResponseItemsStateActive,
					CreatedAt:  now,
					UpdatedAt:  now,
				},
				{
					ID:         "org-2",
					Name:       "Org 2",
					ExternalID: "",
					State:      api.Disabled,
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
			NextCursor: &nextCursor,
		}

		got := organizationsPageFromResponse(resp)

		if len(got.Items) != 2 {
			t.Fatalf("got %d items, want 2", len(got.Items))
		}
		if got.Items[0].ID != "org-1" {
			t.Errorf("Items[0].ID = %q, want 'org-1'", got.Items[0].ID)
		}
		if got.Items[0].ExternalID == nil || *got.Items[0].ExternalID != "ext-1" {
			t.Errorf("Items[0].ExternalID = %v, want 'ext-1'", got.Items[0].ExternalID)
		}
		if got.Items[0].State != OrganizationStateActive {
			t.Errorf("Items[0].State = %q, want %q", got.Items[0].State, OrganizationStateActive)
		}
		if got.Items[1].ExternalID != nil {
			t.Errorf("Items[1].ExternalID = %v, want nil", got.Items[1].ExternalID)
		}
		if got.Items[1].State != OrganizationStateDisabled {
			t.Errorf("Items[1].State = %q, want %q", got.Items[1].State, OrganizationStateDisabled)
		}
		if got.PageInfo.NextCursor == nil || *got.PageInfo.NextCursor != "cursor-abc" {
			t.Errorf("NextCursor = %v, want 'cursor-abc'", got.PageInfo.NextCursor)
		}
		if !got.PageInfo.HasMore {
			t.Error("HasMore = false, want true")
		}
	})

	t.Run("handles nil cursor", func(t *testing.T) {
		resp := &api.GetAPIV1IntegrationsIntegrationIDOrganizationsResponse{
			Items:      api.GetAPIV1IntegrationsIntegrationIDOrganizations_Response_Items{},
			NextCursor: nil,
		}

		got := organizationsPageFromResponse(resp)

		if got.PageInfo.NextCursor != nil {
			t.Errorf("NextCursor = %v, want nil", got.PageInfo.NextCursor)
		}
		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false")
		}
	})

	t.Run("handles empty cursor string", func(t *testing.T) {
		empty := ""
		resp := &api.GetAPIV1IntegrationsIntegrationIDOrganizationsResponse{
			Items:      api.GetAPIV1IntegrationsIntegrationIDOrganizations_Response_Items{},
			NextCursor: &empty,
		}

		got := organizationsPageFromResponse(resp)

		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false for empty cursor")
		}
	})
}

func TestAPIKeyFromResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	expiresAt := now.Add(30 * 24 * time.Hour)

	t.Run("converts all fields", func(t *testing.T) {
		resp := &api.PostAPIV1OrganizationsOrganizationIDKeysResponse{
			ID:        "key-123",
			Name:      "Test Key",
			Key:       "xbl-org-abc123",
			ExpiresAt: expiresAt,
			CreatedAt: now,
			UpdatedAt: now,
		}

		got := apiKeyFromResponse(resp)

		if got.ID != "key-123" {
			t.Errorf("ID = %q, want 'key-123'", got.ID)
		}
		if got.Name != "Test Key" {
			t.Errorf("Name = %q, want 'Test Key'", got.Name)
		}
		if got.Key != "xbl-org-abc123" {
			t.Errorf("Key = %q, want 'xbl-org-abc123'", got.Key)
		}
		if got.ExpiresAt == nil || !got.ExpiresAt.Equal(expiresAt) {
			t.Errorf("ExpiresAt = %v, want %v", got.ExpiresAt, expiresAt)
		}
		if !got.CreatedAt.Equal(now) {
			t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, now)
		}
		if !got.UpdatedAt.Equal(now) {
			t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, now)
		}
	})

	t.Run("handles zero ExpiresAt as nil", func(t *testing.T) {
		resp := &api.PostAPIV1OrganizationsOrganizationIDKeysResponse{
			ID:        "key-456",
			Name:      "Never Expires",
			Key:       "xbl-org-def456",
			ExpiresAt: time.Time{},
			CreatedAt: now,
			UpdatedAt: now,
		}

		got := apiKeyFromResponse(resp)

		if got.ExpiresAt != nil {
			t.Errorf("ExpiresAt = %v, want nil", got.ExpiresAt)
		}
	})
}

func TestUpdateOrganizationNilRequest(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.Organizations.Update(context.TODO(), "org-123", nil)
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

func TestCreateOrganizationNilRequest(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.Organizations.Create(context.TODO(), "integration-123", nil)
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

func TestCreateOrganizationNoMembers(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.Organizations.Create(context.TODO(), "integration-123", &CreateOrganizationRequest{
		Name:    "Test Org",
		Members: []OrganizationMember{},
	})
	if err == nil {
		t.Fatal("expected error for empty members")
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.Code != "ERR_INVALID_REQUEST" {
		t.Errorf("Code = %q, want 'ERR_INVALID_REQUEST'", apiErr.Code)
	}
	if apiErr.Message != "at least one member is required" {
		t.Errorf("Message = %q, want 'at least one member is required'", apiErr.Message)
	}
}

func TestCreateKeyNilRequest(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.Organizations.CreateKey(context.TODO(), "org-123", nil)
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

func TestOrganizationStateConstants(t *testing.T) {
	tests := []struct {
		state OrganizationState
		want  string
	}{
		{OrganizationStateActive, "active"},
		{OrganizationStateDisabled, "disabled"},
	}

	for _, tt := range tests {
		if string(tt.state) != tt.want {
			t.Errorf("OrganizationState = %q, want %q", tt.state, tt.want)
		}
	}
}
