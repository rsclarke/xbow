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
	if got.StartURL == nil || *got.StartURL != "https://example.com" {
		t.Errorf("StartURL = %v, want 'https://example.com'", got.StartURL)
	}
	if got.MaxRequestsPerSecond == nil || *got.MaxRequestsPerSecond != 100 {
		t.Errorf("MaxRequestsPerSecond = %v, want 100", got.MaxRequestsPerSecond)
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
	if got.MaxRequestsPerSecond == nil || *got.MaxRequestsPerSecond != 500 {
		t.Errorf("MaxRequestsPerSecond = %v, want 500", got.MaxRequestsPerSecond)
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
	if got.StartURL != nil {
		t.Errorf("StartURL = %v, want nil (null)", got.StartURL)
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
		if got.PageInfo.NextCursor == nil || *got.PageInfo.NextCursor != "cursor-xyz" {
			t.Errorf("NextCursor = %v, want 'cursor-xyz'", got.PageInfo.NextCursor)
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

		if got.PageInfo.NextCursor != nil {
			t.Errorf("NextCursor = %v, want nil", got.PageInfo.NextCursor)
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

func TestConvertApprovedTimeWindowsFromGet(t *testing.T) {
	t.Run("converts time windows", func(t *testing.T) {
		atw := api.GetAPIV1AssetsAssetID_Response_ApprovedTimeWindows{
			Tz: "Europe/London",
			Entries: api.GetAPIV1AssetsAssetID_Response_ApprovedTimeWindows_AnyOf_Entries{
				{
					StartWeekday: 1,
					StartTime:    "09:00",
					EndWeekday:   5,
					EndTime:      "17:00",
				},
			},
		}

		got := convertApprovedTimeWindowsFromGet(atw)

		if got == nil {
			t.Fatal("expected non-nil result")
		}
		if got.Tz != "Europe/London" {
			t.Errorf("Tz = %q, want 'Europe/London'", got.Tz)
		}
		if len(got.Entries) != 1 {
			t.Fatalf("got %d entries, want 1", len(got.Entries))
		}
		if got.Entries[0].StartWeekday != 1 {
			t.Errorf("StartWeekday = %d, want 1", got.Entries[0].StartWeekday)
		}
		if got.Entries[0].StartTime != "09:00" {
			t.Errorf("StartTime = %q, want '09:00'", got.Entries[0].StartTime)
		}
	})

	t.Run("returns nil for empty", func(t *testing.T) {
		atw := api.GetAPIV1AssetsAssetID_Response_ApprovedTimeWindows{}

		got := convertApprovedTimeWindowsFromGet(atw)

		if got != nil {
			t.Errorf("expected nil, got %+v", got)
		}
	})
}

func TestConvertCredentialsFromGet(t *testing.T) {
	t.Run("converts credentials with all fields", func(t *testing.T) {
		authURI := "otpauth://totp/test"
		creds := api.GetAPIV1AssetsAssetID_Response_Credentials{
			{
				ID:               "cred-1",
				Name:             "Test Cred",
				Type:             api.UsernamePassword,
				Username:         "testuser",
				Password:         "testpass",
				AuthenticatorURI: &authURI,
			},
		}

		got := convertCredentialsFromGet(creds)

		if len(got) != 1 {
			t.Fatalf("got %d credentials, want 1", len(got))
		}
		if got[0].ID != "cred-1" {
			t.Errorf("ID = %q, want 'cred-1'", got[0].ID)
		}
		if got[0].Username != "testuser" {
			t.Errorf("Username = %q, want 'testuser'", got[0].Username)
		}
		if got[0].Type != "username-password" {
			t.Errorf("Type = %q, want 'username-password'", got[0].Type)
		}
		if got[0].AuthenticatorURI == nil || *got[0].AuthenticatorURI != authURI {
			t.Errorf("AuthenticatorURI = %v, want %q", got[0].AuthenticatorURI, authURI)
		}
	})

	t.Run("handles nil optional fields", func(t *testing.T) {
		creds := api.GetAPIV1AssetsAssetID_Response_Credentials{
			{
				ID:       "cred-2",
				Name:     "Basic Cred",
				Type:     api.UsernamePassword,
				Username: "user",
				Password: "pass",
			},
		}

		got := convertCredentialsFromGet(creds)

		if got[0].EmailAddress != nil {
			t.Errorf("EmailAddress = %v, want nil", got[0].EmailAddress)
		}
		if got[0].AuthenticatorURI != nil {
			t.Errorf("AuthenticatorURI = %v, want nil", got[0].AuthenticatorURI)
		}
	})
}

func TestConvertDNSBoundaryRulesFromGet(t *testing.T) {
	includeSubdomains := true
	rules := api.GetAPIV1AssetsAssetID_Response_DNSBoundaryRules{
		{
			ID:                "rule-1",
			Action:            api.AllowVisit,
			Type:              api.Glob,
			Filter:            "*.example.com",
			IncludeSubdomains: &includeSubdomains,
		},
	}

	got := convertDNSBoundaryRulesFromGet(rules)

	if len(got) != 1 {
		t.Fatalf("got %d rules, want 1", len(got))
	}
	if got[0].ID != "rule-1" {
		t.Errorf("ID = %q, want 'rule-1'", got[0].ID)
	}
	if got[0].Action != "allow-visit" {
		t.Errorf("Action = %q, want 'allow-visit'", got[0].Action)
	}
	if got[0].Filter != "*.example.com" {
		t.Errorf("Filter = %q, want '*.example.com'", got[0].Filter)
	}
	if got[0].IncludeSubdomains == nil || !*got[0].IncludeSubdomains {
		t.Errorf("IncludeSubdomains = %v, want true", got[0].IncludeSubdomains)
	}
}

func TestConvertHTTPBoundaryRulesFromGet(t *testing.T) {
	rules := api.GetAPIV1AssetsAssetID_Response_HTTPBoundaryRules{
		{
			ID:     "http-rule-1",
			Action: api.GetAPIV1AssetsAssetIDResponseHTTPBoundaryRulesAnyOfActionDeny,
			Type:   api.Exact,
			Filter: "https://blocked.example.com",
		},
	}

	got := convertHTTPBoundaryRulesFromGet(rules)

	if len(got) != 1 {
		t.Fatalf("got %d rules, want 1", len(got))
	}
	if got[0].Action != HTTPBoundaryRuleActionDeny {
		t.Errorf("Action = %q, want %q", got[0].Action, HTTPBoundaryRuleActionDeny)
	}
	if got[0].Filter != "https://blocked.example.com" {
		t.Errorf("Filter = %q, want 'https://blocked.example.com'", got[0].Filter)
	}
}

func TestConvertHeadersFromGet(t *testing.T) {
	t.Run("skips nil anyOf", func(t *testing.T) {
		headers := map[string]api.GetAPIV1AssetsAssetID_Response_Headers_AnyOf_AdditionalProperties{
			"X-Nil": {},
		}

		got := convertHeadersFromGet(headers)

		if _, exists := got["X-Nil"]; exists {
			t.Errorf("expected X-Nil to be skipped")
		}
	})

	t.Run("returns nil for empty", func(t *testing.T) {
		got := convertHeadersFromGet(nil)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
}

func TestConvertChecksFromGet(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	checks := api.GetAPIV1AssetsAssetID_Response_Checks{
		AssetReachable: api.GetAPIV1AssetsAssetID_Response_Checks_AssetReachable{
			State:   api.Valid,
			Message: "Asset is reachable",
		},
		Credentials: api.GetAPIV1AssetsAssetID_Response_Checks_Credentials{
			State:   api.GetAPIV1AssetsAssetIDResponseChecksCredentialsStateChecking,
			Message: "Validating credentials",
		},
		DNSBoundaryRules: api.GetAPIV1AssetsAssetID_Response_Checks_DNSBoundaryRules{
			State:   api.GetAPIV1AssetsAssetIDResponseChecksDNSBoundaryRulesStateUnchecked,
			Message: "",
		},
		UpdatedAt: now,
	}

	got := convertChecksFromGet(checks)

	if got.AssetReachable.State != AssetCheckStateValid {
		t.Errorf("AssetReachable.State = %q, want %q", got.AssetReachable.State, AssetCheckStateValid)
	}
	if got.AssetReachable.Message != "Asset is reachable" {
		t.Errorf("AssetReachable.Message = %q, want 'Asset is reachable'", got.AssetReachable.Message)
	}
	if got.Credentials.State != AssetCheckStateChecking {
		t.Errorf("Credentials.State = %q, want %q", got.Credentials.State, AssetCheckStateChecking)
	}
	if got.DNSBoundaryRules.State != AssetCheckStateUnchecked {
		t.Errorf("DNSBoundaryRules.State = %q, want %q", got.DNSBoundaryRules.State, AssetCheckStateUnchecked)
	}
	if !got.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, now)
	}
}

func TestAssetFromGetResponseWithAllFields(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	archiveAt := now.Add(30 * 24 * time.Hour)

	resp := &api.GetAPIV1AssetsAssetIDResponse{
		ID:                   "asset-full",
		Name:                 "Full Asset",
		OrganizationID:       "org-123",
		Lifecycle:            api.Active,
		Sku:                  "premium",
		StartURL:             "https://example.com",
		MaxRequestsPerSecond: 200,
		ArchiveAt:            archiveAt,
		ApprovedTimeWindows: api.GetAPIV1AssetsAssetID_Response_ApprovedTimeWindows{
			Tz: "UTC",
			Entries: api.GetAPIV1AssetsAssetID_Response_ApprovedTimeWindows_AnyOf_Entries{
				{StartWeekday: 1, StartTime: "00:00", EndWeekday: 7, EndTime: "23:59"},
			},
		},
		Credentials: api.GetAPIV1AssetsAssetID_Response_Credentials{
			{ID: "c1", Name: "cred", Type: api.UsernamePassword, Username: "u", Password: "p"},
		},
		DNSBoundaryRules: api.GetAPIV1AssetsAssetID_Response_DNSBoundaryRules{
			{ID: "d1", Action: api.AllowVisit, Type: api.Glob, Filter: "*.example.com"},
		},
		HTTPBoundaryRules: api.GetAPIV1AssetsAssetID_Response_HTTPBoundaryRules{
			{ID: "h1", Action: api.GetAPIV1AssetsAssetIDResponseHTTPBoundaryRulesAnyOfActionDeny, Type: api.Exact, Filter: "https://blocked.com"},
		},
		Checks: api.GetAPIV1AssetsAssetID_Response_Checks{
			AssetReachable:   api.GetAPIV1AssetsAssetID_Response_Checks_AssetReachable{State: api.Valid, Message: "ok"},
			Credentials:      api.GetAPIV1AssetsAssetID_Response_Checks_Credentials{State: api.GetAPIV1AssetsAssetIDResponseChecksCredentialsStateValid, Message: "ok"},
			DNSBoundaryRules: api.GetAPIV1AssetsAssetID_Response_Checks_DNSBoundaryRules{State: api.GetAPIV1AssetsAssetIDResponseChecksDNSBoundaryRulesStateValid, Message: "ok"},
			UpdatedAt:        now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	got := assetFromGetResponse(resp)

	if got.ArchiveAt == nil || !got.ArchiveAt.Equal(archiveAt) {
		t.Errorf("ArchiveAt = %v, want %v", got.ArchiveAt, archiveAt)
	}
	if got.ApprovedTimeWindows == nil {
		t.Fatal("ApprovedTimeWindows is nil")
	}
	if got.ApprovedTimeWindows.Tz != "UTC" {
		t.Errorf("ApprovedTimeWindows.Tz = %q, want 'UTC'", got.ApprovedTimeWindows.Tz)
	}
	if len(got.Credentials) != 1 {
		t.Errorf("got %d credentials, want 1", len(got.Credentials))
	}
	if len(got.DNSBoundaryRules) != 1 {
		t.Errorf("got %d DNS rules, want 1", len(got.DNSBoundaryRules))
	}
	if len(got.HTTPBoundaryRules) != 1 {
		t.Errorf("got %d HTTP rules, want 1", len(got.HTTPBoundaryRules))
	}
	if got.Checks == nil {
		t.Fatal("Checks is nil")
	}
	if got.Checks.AssetReachable.State != AssetCheckStateValid {
		t.Errorf("Checks.AssetReachable.State = %q, want %q", got.Checks.AssetReachable.State, AssetCheckStateValid)
	}
}

func TestAssetFromGetResponseArchiveAtZero(t *testing.T) {
	resp := &api.GetAPIV1AssetsAssetIDResponse{
		ID:             "asset-no-archive",
		Name:           "No Archive",
		OrganizationID: "org-123",
		Lifecycle:      api.Active,
		Sku:            "standard",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	got := assetFromGetResponse(resp)

	if got.ArchiveAt != nil {
		t.Errorf("ArchiveAt = %v, want nil for zero time", got.ArchiveAt)
	}
}

func TestConvertApprovedTimeWindowsToBody(t *testing.T) {
	t.Run("converts time windows", func(t *testing.T) {
		atw := &ApprovedTimeWindows{
			Tz: "America/New_York",
			Entries: []TimeWindowEntry{
				{StartWeekday: 1, StartTime: "09:00", EndWeekday: 5, EndTime: "17:00"},
			},
		}

		got := convertApprovedTimeWindowsToBody(atw)

		if got.Tz != "America/New_York" {
			t.Errorf("Tz = %q, want 'America/New_York'", got.Tz)
		}
		if len(got.Entries) != 1 {
			t.Fatalf("got %d entries, want 1", len(got.Entries))
		}
		if got.Entries[0].StartTime != "09:00" {
			t.Errorf("StartTime = %q, want '09:00'", got.Entries[0].StartTime)
		}
		if got.Entries[0].StartWeekday != 1 {
			t.Errorf("StartWeekday = %v, want 1", got.Entries[0].StartWeekday)
		}
	})

	t.Run("returns empty for nil", func(t *testing.T) {
		got := convertApprovedTimeWindowsToBody(nil)

		if got.Tz != "" {
			t.Errorf("Tz = %q, want empty", got.Tz)
		}
		if len(got.Entries) != 0 {
			t.Errorf("Entries = %v, want empty", got.Entries)
		}
	})
}

func TestConvertCredentialsToBody(t *testing.T) {
	t.Run("converts credentials", func(t *testing.T) {
		authURI := "otpauth://totp/test"
		creds := []Credential{
			{
				ID:               "cred-1",
				Name:             "Test",
				Type:             "username-password",
				Username:         "user",
				Password:         "pass",
				AuthenticatorURI: &authURI,
			},
		}

		got := convertCredentialsToBody(creds)

		if len(got) != 1 {
			t.Fatalf("got %d credentials, want 1", len(got))
		}
		if got[0].ID != "cred-1" {
			t.Errorf("ID = %q, want 'cred-1'", got[0].ID)
		}
		if string(got[0].Type) != "username-password" {
			t.Errorf("Type = %q, want 'username-password'", got[0].Type)
		}
		if got[0].AuthenticatorURI == nil || *got[0].AuthenticatorURI != authURI {
			t.Errorf("AuthenticatorURI = %v, want %q", got[0].AuthenticatorURI, authURI)
		}
	})

	t.Run("returns nil for empty", func(t *testing.T) {
		got := convertCredentialsToBody(nil)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
}

func TestConvertDNSBoundaryRulesToBody(t *testing.T) {
	t.Run("converts rules", func(t *testing.T) {
		includeSubdomains := true
		rules := []DNSBoundaryRule{
			{
				ID:                "rule-1",
				Action:            DNSBoundaryRuleActionDeny,
				Type:              "glob",
				Filter:            "*.blocked.com",
				IncludeSubdomains: &includeSubdomains,
			},
		}

		got := convertDNSBoundaryRulesToBody(rules)

		if len(got) != 1 {
			t.Fatalf("got %d rules, want 1", len(got))
		}
		if got[0].ID != "rule-1" {
			t.Errorf("ID = %q, want 'rule-1'", got[0].ID)
		}
		if string(got[0].Action) != "deny" {
			t.Errorf("Action = %q, want 'deny'", got[0].Action)
		}
		if got[0].IncludeSubdomains == nil || !*got[0].IncludeSubdomains {
			t.Errorf("IncludeSubdomains = %v, want true", got[0].IncludeSubdomains)
		}
	})

	t.Run("returns nil for empty", func(t *testing.T) {
		got := convertDNSBoundaryRulesToBody(nil)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
}

func TestConvertHeadersToBody(t *testing.T) {
	t.Run("converts single value as A", func(t *testing.T) {
		headers := map[string][]string{
			"X-Single": {"value1"},
		}

		got := convertHeadersToBody(headers)

		if got["X-Single"].PutAPIV1AssetsAssetIDBody_Headers_AnyOf_AdditionalProperties_AnyOf == nil {
			t.Fatal("anyOf is nil")
		}
		anyOf := got["X-Single"].PutAPIV1AssetsAssetIDBody_Headers_AnyOf_AdditionalProperties_AnyOf
		if !anyOf.IsA() {
			t.Error("expected IsA() to be true for single value")
		}
		if anyOf.A != "value1" {
			t.Errorf("A = %q, want 'value1'", anyOf.A)
		}
	})

	t.Run("converts multiple values as B", func(t *testing.T) {
		headers := map[string][]string{
			"X-Multi": {"val1", "val2"},
		}

		got := convertHeadersToBody(headers)

		anyOf := got["X-Multi"].PutAPIV1AssetsAssetIDBody_Headers_AnyOf_AdditionalProperties_AnyOf
		if !anyOf.IsB() {
			t.Error("expected IsB() to be true for multiple values")
		}
		if len(anyOf.B) != 2 {
			t.Fatalf("B has %d values, want 2", len(anyOf.B))
		}
		if anyOf.B[0] != "val1" || anyOf.B[1] != "val2" {
			t.Errorf("B = %v, want [val1, val2]", anyOf.B)
		}
	})

	t.Run("returns nil for empty", func(t *testing.T) {
		got := convertHeadersToBody(nil)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
}

func TestConvertHTTPBoundaryRulesToBody(t *testing.T) {
	t.Run("converts rules", func(t *testing.T) {
		rules := []HTTPBoundaryRule{
			{
				ID:     "http-1",
				Action: HTTPBoundaryRuleActionAllowAttack,
				Type:   "exact",
				Filter: "https://allowed.com",
			},
		}

		got := convertHTTPBoundaryRulesToBody(rules)

		if len(got) != 1 {
			t.Fatalf("got %d rules, want 1", len(got))
		}
		if got[0].ID != "http-1" {
			t.Errorf("ID = %q, want 'http-1'", got[0].ID)
		}
		if string(got[0].Action) != "allow-attack" {
			t.Errorf("Action = %q, want 'allow-attack'", got[0].Action)
		}
		if got[0].Filter != "https://allowed.com" {
			t.Errorf("Filter = %q, want 'https://allowed.com'", got[0].Filter)
		}
	})

	t.Run("returns nil for empty", func(t *testing.T) {
		got := convertHTTPBoundaryRulesToBody(nil)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
}

func TestUpdateAssetNilRequest(t *testing.T) {
	client, _ := NewClient(WithOrganizationKey("test-key"))

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
	client, _ := NewClient(WithOrganizationKey("test-key"))

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

func TestConvertSimpleErrorFromGet(t *testing.T) {
	t.Run("returns nil for empty type", func(t *testing.T) {
		got := convertSimpleErrorFromGet("")
		if got != nil {
			t.Errorf("expected nil, got %+v", got)
		}
	})

	t.Run("returns error for non-empty type", func(t *testing.T) {
		got := convertSimpleErrorFromGet("invalid-credentials")
		if got == nil {
			t.Fatal("expected non-nil error")
		}
		if got.Type != "invalid-credentials" {
			t.Errorf("Type = %q, want 'invalid-credentials'", got.Type)
		}
	})
}

func TestConvertAssetReachableErrorFromGet(t *testing.T) {
	t.Run("returns nil when oneOf is nil", func(t *testing.T) {
		errData := api.GetAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error{}
		got := convertAssetReachableErrorFromGet(errData)
		if got != nil {
			t.Errorf("expected nil, got %+v", got)
		}
	})
}

func TestDNSBoundaryRuleActionConstants(t *testing.T) {
	tests := []struct {
		action DNSBoundaryRuleAction
		want   string
	}{
		{DNSBoundaryRuleActionAllowAttack, "allow-attack"},
		{DNSBoundaryRuleActionAllowVisit, "allow-visit"},
		{DNSBoundaryRuleActionDeny, "deny"},
	}

	for _, tt := range tests {
		if string(tt.action) != tt.want {
			t.Errorf("DNSBoundaryRuleAction = %q, want %q", tt.action, tt.want)
		}
	}
}

func TestHTTPBoundaryRuleActionConstants(t *testing.T) {
	tests := []struct {
		action HTTPBoundaryRuleAction
		want   string
	}{
		{HTTPBoundaryRuleActionAllowAttack, "allow-attack"},
		{HTTPBoundaryRuleActionAllowAuth, "allow-auth"},
		{HTTPBoundaryRuleActionAllowVisit, "allow-visit"},
		{HTTPBoundaryRuleActionDeny, "deny"},
	}

	for _, tt := range tests {
		if string(tt.action) != tt.want {
			t.Errorf("HTTPBoundaryRuleAction = %q, want %q", tt.action, tt.want)
		}
	}
}
