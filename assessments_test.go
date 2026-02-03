package xbow

import (
	"testing"
	"time"

	"github.com/rsclarke/xbow/internal/api"
)

func TestAssessmentFromGetResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	resp := &api.GetAPIV1AssessmentsAssessmentIDResponse{
		ID:             "assess-123",
		Name:           "Test Assessment",
		AssetID:        "asset-456",
		OrganizationID: "org-789",
		State:          api.Running,
		Progress:       0.75,
		AttackCredits:  100,
		CreatedAt:      now,
		UpdatedAt:      now.Add(time.Hour),
	}

	got := assessmentFromGetResponse(resp)

	if got.ID != "assess-123" {
		t.Errorf("ID = %q, want 'assess-123'", got.ID)
	}
	if got.Name != "Test Assessment" {
		t.Errorf("Name = %q, want 'Test Assessment'", got.Name)
	}
	if got.AssetID != "asset-456" {
		t.Errorf("AssetID = %q, want 'asset-456'", got.AssetID)
	}
	if got.OrganizationID != "org-789" {
		t.Errorf("OrganizationID = %q, want 'org-789'", got.OrganizationID)
	}
	if got.State != AssessmentStateRunning {
		t.Errorf("State = %q, want %q", got.State, AssessmentStateRunning)
	}
	if got.Progress != 0.75 {
		t.Errorf("Progress = %f, want 0.75", got.Progress)
	}
	if got.AttackCredits != 100 {
		t.Errorf("AttackCredits = %d, want 100", got.AttackCredits)
	}
	if !got.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, now)
	}
	if !got.UpdatedAt.Equal(now.Add(time.Hour)) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, now.Add(time.Hour))
	}
}

func TestAssessmentsPageFromResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	t.Run("converts items and pagination info", func(t *testing.T) {
		nextCursor := "cursor-abc"
		resp := &api.GetAPIV1AssetsAssetIDAssessmentsResponse{
			Items: api.GetAPIV1AssetsAssetIDAssessments_Response_Items{
				{
					ID:        "a1",
					Name:      "Assessment 1",
					State:     api.GetAPIV1AssetsAssetIDAssessmentsResponseItemsStateSucceeded,
					Progress:  1.0,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "a2",
					Name:      "Assessment 2",
					State:     api.GetAPIV1AssetsAssetIDAssessmentsResponseItemsStatePaused,
					Progress:  0.5,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			NextCursor: &nextCursor,
		}

		got := assessmentsPageFromResponse(resp)

		if len(got.Items) != 2 {
			t.Fatalf("got %d items, want 2", len(got.Items))
		}
		if got.Items[0].ID != "a1" {
			t.Errorf("Items[0].ID = %q, want 'a1'", got.Items[0].ID)
		}
		if got.Items[0].State != AssessmentStateSucceeded {
			t.Errorf("Items[0].State = %q, want %q", got.Items[0].State, AssessmentStateSucceeded)
		}
		if got.Items[1].State != AssessmentStatePaused {
			t.Errorf("Items[1].State = %q, want %q", got.Items[1].State, AssessmentStatePaused)
		}
		if got.PageInfo.NextCursor != "cursor-abc" {
			t.Errorf("NextCursor = %q, want 'cursor-abc'", got.PageInfo.NextCursor)
		}
		if !got.PageInfo.HasMore {
			t.Error("HasMore = false, want true")
		}
	})

	t.Run("handles nil cursor", func(t *testing.T) {
		resp := &api.GetAPIV1AssetsAssetIDAssessmentsResponse{
			Items:      api.GetAPIV1AssetsAssetIDAssessments_Response_Items{},
			NextCursor: nil,
		}

		got := assessmentsPageFromResponse(resp)

		if got.PageInfo.NextCursor != "" {
			t.Errorf("NextCursor = %q, want empty", got.PageInfo.NextCursor)
		}
		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false")
		}
	})

	t.Run("handles empty cursor string", func(t *testing.T) {
		empty := ""
		resp := &api.GetAPIV1AssetsAssetIDAssessmentsResponse{
			Items:      api.GetAPIV1AssetsAssetIDAssessments_Response_Items{},
			NextCursor: &empty,
		}

		got := assessmentsPageFromResponse(resp)

		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false for empty cursor")
		}
	})
}

func TestCreateAssessmentNilRequest(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.Assessments.Create(nil, "asset-123", nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.Code != "ERR_INVALID_REQUEST" {
		t.Errorf("Code = %q, want 'ERR_INVALID_REQUEST'", apiErr.Code)
	}
}

func TestPtrValue(t *testing.T) {
	t.Run("returns value when not nil", func(t *testing.T) {
		s := "hello"
		if got := ptrValue(&s); got != "hello" {
			t.Errorf("ptrValue() = %q, want 'hello'", got)
		}
	})

	t.Run("returns zero when nil", func(t *testing.T) {
		var p *string
		if got := ptrValue(p); got != "" {
			t.Errorf("ptrValue(nil) = %q, want empty", got)
		}
	})

	t.Run("works with int", func(t *testing.T) {
		n := 42
		if got := ptrValue(&n); got != 42 {
			t.Errorf("ptrValue() = %d, want 42", got)
		}
	})
}
