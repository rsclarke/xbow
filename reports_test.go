package xbow

import (
	"testing"
	"time"

	"github.com/rsclarke/xbow/internal/api"
)

func TestReportSummaryFromResponse(t *testing.T) {
	resp := &api.GetAPIV1ReportsReportIDSummaryResponse{
		Markdown: "XBOW identified **5 critical** findings in this assessment.",
	}

	got := reportSummaryFromResponse(resp)

	if got.Markdown != "XBOW identified **5 critical** findings in this assessment." {
		t.Errorf("Markdown = %q, want markdown summary", got.Markdown)
	}
}

func TestReportsPageFromResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	t.Run("converts items and pagination info", func(t *testing.T) {
		nextCursor := "cursor-abc"
		resp := &api.GetAPIV1AssetsAssetIDReportsResponse{
			Items: api.GetAPIV1AssetsAssetIDReports_Response_Items{
				{
					ID:        "report-1",
					Version:   1,
					CreatedAt: now,
				},
				{
					ID:        "report-2",
					Version:   2,
					CreatedAt: now.Add(time.Hour),
				},
			},
			NextCursor: &nextCursor,
		}

		got := reportsPageFromResponse(resp)

		if len(got.Items) != 2 {
			t.Fatalf("got %d items, want 2", len(got.Items))
		}
		if got.Items[0].ID != "report-1" {
			t.Errorf("Items[0].ID = %q, want 'report-1'", got.Items[0].ID)
		}
		if got.Items[0].Version != 1 {
			t.Errorf("Items[0].Version = %d, want 1", got.Items[0].Version)
		}
		if !got.Items[0].CreatedAt.Equal(now) {
			t.Errorf("Items[0].CreatedAt = %v, want %v", got.Items[0].CreatedAt, now)
		}
		if got.Items[1].ID != "report-2" {
			t.Errorf("Items[1].ID = %q, want 'report-2'", got.Items[1].ID)
		}
		if got.Items[1].Version != 2 {
			t.Errorf("Items[1].Version = %d, want 2", got.Items[1].Version)
		}
		if got.PageInfo.NextCursor == nil || *got.PageInfo.NextCursor != "cursor-abc" {
			t.Errorf("NextCursor = %v, want 'cursor-abc'", got.PageInfo.NextCursor)
		}
		if !got.PageInfo.HasMore {
			t.Error("HasMore = false, want true")
		}
	})

	t.Run("handles nil cursor", func(t *testing.T) {
		resp := &api.GetAPIV1AssetsAssetIDReportsResponse{
			Items:      api.GetAPIV1AssetsAssetIDReports_Response_Items{},
			NextCursor: nil,
		}

		got := reportsPageFromResponse(resp)

		if got.PageInfo.NextCursor != nil {
			t.Errorf("NextCursor = %v, want nil", got.PageInfo.NextCursor)
		}
		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false")
		}
	})

	t.Run("handles empty cursor string", func(t *testing.T) {
		empty := ""
		resp := &api.GetAPIV1AssetsAssetIDReportsResponse{
			Items:      api.GetAPIV1AssetsAssetIDReports_Response_Items{},
			NextCursor: &empty,
		}

		got := reportsPageFromResponse(resp)

		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false for empty cursor")
		}
	})
}

func TestReportListItemFields(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	resp := &api.GetAPIV1AssetsAssetIDReportsResponse{
		Items: api.GetAPIV1AssetsAssetIDReports_Response_Items{
			{
				ID:        "report-123",
				Version:   42,
				CreatedAt: now,
			},
		},
	}

	got := reportsPageFromResponse(resp)

	if len(got.Items) != 1 {
		t.Fatalf("got %d items, want 1", len(got.Items))
	}

	item := got.Items[0]
	if item.ID != "report-123" {
		t.Errorf("ID = %q, want 'report-123'", item.ID)
	}
	if item.Version != 42 {
		t.Errorf("Version = %d, want 42", item.Version)
	}
	if !item.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", item.CreatedAt, now)
	}
}

func TestReportSummaryEmptyMarkdown(t *testing.T) {
	resp := &api.GetAPIV1ReportsReportIDSummaryResponse{
		Markdown: "",
	}

	got := reportSummaryFromResponse(resp)

	if got.Markdown != "" {
		t.Errorf("Markdown = %q, want empty string", got.Markdown)
	}
}
