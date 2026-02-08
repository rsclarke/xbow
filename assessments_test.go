package xbow

import (
	"context"
	"encoding/json"
	"errors"
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
		if got.PageInfo.NextCursor == nil || *got.PageInfo.NextCursor != "cursor-abc" {
			t.Errorf("NextCursor = %v, want 'cursor-abc'", got.PageInfo.NextCursor)
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

		if got.PageInfo.NextCursor != nil {
			t.Errorf("NextCursor = %v, want nil", got.PageInfo.NextCursor)
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
	client, _ := NewClient(WithOrganizationKey("test-key"))

	_, err := client.Assessments.Create(context.TODO(), "asset-123", nil)
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

func TestAssessmentFromGetResponseLargeAttackCredits(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	// Test max value from OpenAPI spec: 9007199254740991
	largeCredits := 9007199254740991

	resp := &api.GetAPIV1AssessmentsAssessmentIDResponse{
		ID:             "assess-123",
		Name:           "Test Assessment",
		AssetID:        "asset-456",
		OrganizationID: "org-789",
		State:          api.Running,
		Progress:       0.5,
		AttackCredits:  largeCredits,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	got := assessmentFromGetResponse(resp)

	if got.AttackCredits != int64(largeCredits) {
		t.Errorf("AttackCredits = %d, want %d", got.AttackCredits, largeCredits)
	}
}

func TestAssessmentListItemFields(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	resp := &api.GetAPIV1AssetsAssetIDAssessmentsResponse{
		Items: api.GetAPIV1AssetsAssetIDAssessments_Response_Items{
			{
				ID:        "a1",
				Name:      "Assessment 1",
				State:     api.GetAPIV1AssetsAssetIDAssessmentsResponseItemsStateRunning,
				Progress:  0.5,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	got := assessmentsPageFromResponse(resp)

	item := got.Items[0]
	if item.ID != "a1" {
		t.Errorf("ID = %q, want 'a1'", item.ID)
	}
	if item.Name != "Assessment 1" {
		t.Errorf("Name = %q, want 'Assessment 1'", item.Name)
	}
	if item.State != AssessmentStateRunning {
		t.Errorf("State = %q, want %q", item.State, AssessmentStateRunning)
	}
	if item.Progress != 0.5 {
		t.Errorf("Progress = %f, want 0.5", item.Progress)
	}
	if !item.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", item.CreatedAt, now)
	}
	if !item.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", item.UpdatedAt, now)
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

type fakeOneOf struct {
	data json.RawMessage
}

func (f *fakeOneOf) Raw() json.RawMessage { return f.data }

type fakeItem struct {
	oneOf *fakeOneOf
}

func TestConvertRecentEvents(t *testing.T) {
	now := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	getOneOf := func(item fakeItem) rawUnion {
		if item.oneOf == nil {
			return nil
		}
		return item.oneOf
	}

	t.Run("paused event", func(t *testing.T) {
		items := []fakeItem{
			{oneOf: &fakeOneOf{data: json.RawMessage(`{"name":"paused","timestamp":"2025-06-15T12:00:00Z"}`)}},
		}
		got := convertRecentEvents(items, getOneOf)
		if len(got) != 1 {
			t.Fatalf("len = %d, want 1", len(got))
		}
		if got[0].Name != "paused" {
			t.Errorf("Name = %q, want 'paused'", got[0].Name)
		}
		if !got[0].Timestamp.Equal(now) {
			t.Errorf("Timestamp = %v, want %v", got[0].Timestamp, now)
		}
		if got[0].Reason != "" {
			t.Errorf("Reason = %q, want empty", got[0].Reason)
		}
	})

	t.Run("auto-paused event with reason", func(t *testing.T) {
		items := []fakeItem{
			{oneOf: &fakeOneOf{data: json.RawMessage(`{"name":"auto-paused","timestamp":"2025-06-15T12:00:00Z","reason":"out-of-scope"}`)}},
		}
		got := convertRecentEvents(items, getOneOf)
		if len(got) != 1 {
			t.Fatalf("len = %d, want 1", len(got))
		}
		if got[0].Name != "auto-paused" {
			t.Errorf("Name = %q, want 'auto-paused'", got[0].Name)
		}
		if got[0].Reason != "out-of-scope" {
			t.Errorf("Reason = %q, want 'out-of-scope'", got[0].Reason)
		}
	})

	t.Run("nil oneOf skipped", func(t *testing.T) {
		items := []fakeItem{
			{oneOf: nil},
			{oneOf: &fakeOneOf{data: json.RawMessage(`{"name":"resumed","timestamp":"2025-06-15T12:00:00Z"}`)}},
		}
		got := convertRecentEvents(items, getOneOf)
		if len(got) != 1 {
			t.Fatalf("len = %d, want 1", len(got))
		}
		if got[0].Name != "resumed" {
			t.Errorf("Name = %q, want 'resumed'", got[0].Name)
		}
	})

	t.Run("invalid JSON skipped", func(t *testing.T) {
		items := []fakeItem{
			{oneOf: &fakeOneOf{data: json.RawMessage(`{invalid`)}},
		}
		got := convertRecentEvents(items, getOneOf)
		if len(got) != 0 {
			t.Fatalf("len = %d, want 0", len(got))
		}
	})

	t.Run("empty input", func(t *testing.T) {
		got := convertRecentEvents([]fakeItem{}, getOneOf)
		if len(got) != 0 {
			t.Fatalf("len = %d, want 0", len(got))
		}
	})

	t.Run("multiple events", func(t *testing.T) {
		items := []fakeItem{
			{oneOf: &fakeOneOf{data: json.RawMessage(`{"name":"paused","timestamp":"2025-06-15T12:00:00Z"}`)}},
			{oneOf: &fakeOneOf{data: json.RawMessage(`{"name":"resumed","timestamp":"2025-06-15T13:00:00Z"}`)}},
		}
		got := convertRecentEvents(items, getOneOf)
		if len(got) != 2 {
			t.Fatalf("len = %d, want 2", len(got))
		}
		if got[0].Name != "paused" {
			t.Errorf("got[0].Name = %q, want 'paused'", got[0].Name)
		}
		if got[1].Name != "resumed" {
			t.Errorf("got[1].Name = %q, want 'resumed'", got[1].Name)
		}
	})
}
