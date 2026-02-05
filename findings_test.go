package xbow

import (
	"testing"
	"time"

	"github.com/rsclarke/xbow/internal/api"
)

func TestFindingFromGetResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	resp := &api.GetAPIV1FindingsFindingIDResponse{
		ID:          "finding-123",
		Name:        "SQL Injection",
		Severity:    api.GetAPIV1FindingsFindingIDResponseSeverityCritical,
		State:       api.GetAPIV1FindingsFindingIDResponseStateOpen,
		Summary:     "A **critical** SQL injection vulnerability was found.",
		Impact:      "Attackers can access the database.",
		Mitigations: "Use parameterized queries.",
		Recipe:      "1. Visit /search\n2. Enter `' OR 1=1 --`",
		Evidence:    "' OR '1'='1' --",
		CreatedAt:   now,
		UpdatedAt:   now.Add(time.Hour),
	}

	got := findingFromGetResponse(resp)

	if got.ID != "finding-123" {
		t.Errorf("ID = %q, want 'finding-123'", got.ID)
	}
	if got.Name != "SQL Injection" {
		t.Errorf("Name = %q, want 'SQL Injection'", got.Name)
	}
	if got.Severity != FindingSeverityCritical {
		t.Errorf("Severity = %q, want %q", got.Severity, FindingSeverityCritical)
	}
	if got.State != FindingStateOpen {
		t.Errorf("State = %q, want %q", got.State, FindingStateOpen)
	}
	if got.Summary != "A **critical** SQL injection vulnerability was found." {
		t.Errorf("Summary = %q, want markdown summary", got.Summary)
	}
	if got.Impact != "Attackers can access the database." {
		t.Errorf("Impact = %q, want impact description", got.Impact)
	}
	if got.Mitigations != "Use parameterized queries." {
		t.Errorf("Mitigations = %q, want mitigation description", got.Mitigations)
	}
	if got.Recipe != "1. Visit /search\n2. Enter `' OR 1=1 --`" {
		t.Errorf("Recipe = %q, want recipe steps", got.Recipe)
	}
	if got.Evidence != "' OR '1'='1' --" {
		t.Errorf("Evidence = %q, want evidence", got.Evidence)
	}
	if !got.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, now)
	}
	if !got.UpdatedAt.Equal(now.Add(time.Hour)) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, now.Add(time.Hour))
	}
}

func TestFindingsPageFromResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	t.Run("converts items and pagination info", func(t *testing.T) {
		nextCursor := "cursor-abc"
		resp := &api.GetAPIV1AssetsAssetIDFindingsResponse{
			Items: api.GetAPIV1AssetsAssetIDFindings_Response_Items{
				{
					ID:        "f1",
					Name:      "XSS",
					Severity:  api.High,
					State:     api.Open,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "f2",
					Name:      "CSRF",
					Severity:  api.Medium,
					State:     api.Fixed,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			NextCursor: &nextCursor,
		}

		got := findingsPageFromResponse(resp)

		if len(got.Items) != 2 {
			t.Fatalf("got %d items, want 2", len(got.Items))
		}
		if got.Items[0].ID != "f1" {
			t.Errorf("Items[0].ID = %q, want 'f1'", got.Items[0].ID)
		}
		if got.Items[0].Severity != FindingSeverityHigh {
			t.Errorf("Items[0].Severity = %q, want %q", got.Items[0].Severity, FindingSeverityHigh)
		}
		if got.Items[0].State != FindingStateOpen {
			t.Errorf("Items[0].State = %q, want %q", got.Items[0].State, FindingStateOpen)
		}
		if got.Items[1].Severity != FindingSeverityMedium {
			t.Errorf("Items[1].Severity = %q, want %q", got.Items[1].Severity, FindingSeverityMedium)
		}
		if got.Items[1].State != FindingStateFixed {
			t.Errorf("Items[1].State = %q, want %q", got.Items[1].State, FindingStateFixed)
		}
		if got.PageInfo.NextCursor == nil || *got.PageInfo.NextCursor != "cursor-abc" {
			t.Errorf("NextCursor = %v, want 'cursor-abc'", got.PageInfo.NextCursor)
		}
		if !got.PageInfo.HasMore {
			t.Error("HasMore = false, want true")
		}
	})

	t.Run("handles nil cursor", func(t *testing.T) {
		resp := &api.GetAPIV1AssetsAssetIDFindingsResponse{
			Items:      api.GetAPIV1AssetsAssetIDFindings_Response_Items{},
			NextCursor: nil,
		}

		got := findingsPageFromResponse(resp)

		if got.PageInfo.NextCursor != nil {
			t.Errorf("NextCursor = %v, want nil", got.PageInfo.NextCursor)
		}
		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false")
		}
	})

	t.Run("handles empty cursor string", func(t *testing.T) {
		empty := ""
		resp := &api.GetAPIV1AssetsAssetIDFindingsResponse{
			Items:      api.GetAPIV1AssetsAssetIDFindings_Response_Items{},
			NextCursor: &empty,
		}

		got := findingsPageFromResponse(resp)

		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false for empty cursor")
		}
	})
}

func TestFindingSeverityValues(t *testing.T) {
	tests := []struct {
		generated api.GetAPIV1FindingsFindingIDResponseSeverity
		want      FindingSeverity
	}{
		{api.GetAPIV1FindingsFindingIDResponseSeverityCritical, FindingSeverityCritical},
		{api.GetAPIV1FindingsFindingIDResponseSeverityHigh, FindingSeverityHigh},
		{api.GetAPIV1FindingsFindingIDResponseSeverityMedium, FindingSeverityMedium},
		{api.GetAPIV1FindingsFindingIDResponseSeverityLow, FindingSeverityLow},
		{api.GetAPIV1FindingsFindingIDResponseSeverityInformational, FindingSeverityInformational},
	}

	for _, tt := range tests {
		t.Run(string(tt.want), func(t *testing.T) {
			resp := &api.GetAPIV1FindingsFindingIDResponse{
				ID:          "f1",
				Name:        "Test",
				Severity:    tt.generated,
				State:       api.GetAPIV1FindingsFindingIDResponseStateOpen,
				Summary:     "test",
				Impact:      "test",
				Mitigations: "test",
				Recipe:      "test",
				Evidence:    "test",
			}

			got := findingFromGetResponse(resp)
			if got.Severity != tt.want {
				t.Errorf("Severity = %q, want %q", got.Severity, tt.want)
			}
		})
	}
}

func TestFindingStateValues(t *testing.T) {
	tests := []struct {
		generated api.GetAPIV1FindingsFindingIDResponseState
		want      FindingState
	}{
		{api.GetAPIV1FindingsFindingIDResponseStateOpen, FindingStateOpen},
		{api.GetAPIV1FindingsFindingIDResponseStateChallenged, FindingStateChallenged},
		{api.GetAPIV1FindingsFindingIDResponseStateConfirmed, FindingStateConfirmed},
		{api.GetAPIV1FindingsFindingIDResponseStateInvalid, FindingStateInvalid},
		{api.GetAPIV1FindingsFindingIDResponseStateFixed, FindingStateFixed},
	}

	for _, tt := range tests {
		t.Run(string(tt.want), func(t *testing.T) {
			resp := &api.GetAPIV1FindingsFindingIDResponse{
				ID:          "f1",
				Name:        "Test",
				Severity:    api.GetAPIV1FindingsFindingIDResponseSeverityMedium,
				State:       tt.generated,
				Summary:     "test",
				Impact:      "test",
				Mitigations: "test",
				Recipe:      "test",
				Evidence:    "test",
			}

			got := findingFromGetResponse(resp)
			if got.State != tt.want {
				t.Errorf("State = %q, want %q", got.State, tt.want)
			}
		})
	}
}

func TestFindingListItemFields(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	resp := &api.GetAPIV1AssetsAssetIDFindingsResponse{
		Items: api.GetAPIV1AssetsAssetIDFindings_Response_Items{
			{
				ID:        "f1",
				Name:      "Test Finding",
				Severity:  api.Critical,
				State:     api.Challenged,
				CreatedAt: now,
				UpdatedAt: now.Add(time.Hour),
			},
		},
	}

	got := findingsPageFromResponse(resp)

	item := got.Items[0]
	if item.ID != "f1" {
		t.Errorf("ID = %q, want 'f1'", item.ID)
	}
	if item.Name != "Test Finding" {
		t.Errorf("Name = %q, want 'Test Finding'", item.Name)
	}
	if item.Severity != FindingSeverityCritical {
		t.Errorf("Severity = %q, want %q", item.Severity, FindingSeverityCritical)
	}
	if item.State != FindingStateChallenged {
		t.Errorf("State = %q, want %q", item.State, FindingStateChallenged)
	}
	if !item.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", item.CreatedAt, now)
	}
	if !item.UpdatedAt.Equal(now.Add(time.Hour)) {
		t.Errorf("UpdatedAt = %v, want %v", item.UpdatedAt, now.Add(time.Hour))
	}
}

func TestAssessmentFromVerifyFixResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	resp := &api.PostAPIV1FindingsFindingIDVerifyFixResponse{
		ID:             "assess-123",
		Name:           "Fix Verification - SQL Injection",
		AssetID:        "asset-456",
		OrganizationID: "org-789",
		State:          api.PostAPIV1FindingsFindingIDVerifyFixResponseStateWaitingForCapacity,
		Progress:       0.0,
		AttackCredits:  40,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	got := assessmentFromVerifyFixResponse(resp)

	if got.ID != "assess-123" {
		t.Errorf("ID = %q, want 'assess-123'", got.ID)
	}
	if got.Name != "Fix Verification - SQL Injection" {
		t.Errorf("Name = %q, want 'Fix Verification - SQL Injection'", got.Name)
	}
	if got.AssetID != "asset-456" {
		t.Errorf("AssetID = %q, want 'asset-456'", got.AssetID)
	}
	if got.OrganizationID != "org-789" {
		t.Errorf("OrganizationID = %q, want 'org-789'", got.OrganizationID)
	}
	if got.State != AssessmentStateWaitingForCapacity {
		t.Errorf("State = %q, want %q", got.State, AssessmentStateWaitingForCapacity)
	}
	if got.Progress != 0.0 {
		t.Errorf("Progress = %f, want 0.0", got.Progress)
	}
	if got.AttackCredits != 40 {
		t.Errorf("AttackCredits = %d, want 40", got.AttackCredits)
	}
}
