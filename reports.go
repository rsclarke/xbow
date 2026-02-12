package xbow

import (
	"context"
	"fmt"
	"io"
	"iter"
	"net/http"

	"github.com/rsclarke/xbow/internal/api"
)

// ReportsService handles report-related API calls.
type ReportsService struct {
	client *Client
}

// Get returns the report PDF as a streaming reader.
// The caller is responsible for closing the returned ReadCloser.
//
// Example:
//
//	rc, err := client.Reports.Get(ctx, "report-id")
//	if err != nil {
//	    return err
//	}
//	defer rc.Close()
//
//	// Write to file
//	f, _ := os.Create("report.pdf")
//	defer f.Close()
//	io.Copy(f, rc)
func (s *ReportsService) Get(ctx context.Context, id string) (io.ReadCloser, error) {
	if id == "" {
		return nil, &Error{Code: "ERR_INVALID_PARAM", Message: "report id is required"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/reports/%s", id)
	return s.client.doStream(ctx, http.MethodGet, path, auth)
}

// GetSummary retrieves the markdown summary of a report by ID.
func (s *ReportsService) GetSummary(ctx context.Context, id string) (*ReportSummary, error) {
	if id == "" {
		return nil, &Error{Code: "ERR_INVALID_PARAM", Message: "report id is required"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	opts := &api.GetAPIV1ReportsReportIDSummaryRequestOptions{
		PathParams: &api.GetAPIV1ReportsReportIDSummaryPath{
			ReportID: id,
		},
		Header: &api.GetAPIV1ReportsReportIDSummaryHeaders{
			XXBOWAPIVersion: api.GetAPIV1ReportsReportIDSummaryHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.GetAPIV1ReportsReportIDSummary(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return reportSummaryFromResponse(resp), nil
}

// ListByAsset returns a page of reports for an asset.
func (s *ReportsService) ListByAsset(ctx context.Context, assetID string, opts *ListOptions) (*Page[ReportListItem], error) {
	if assetID == "" {
		return nil, &Error{Code: "ERR_INVALID_PARAM", Message: "asset id is required"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	reqOpts := &api.GetAPIV1AssetsAssetIDReportsRequestOptions{
		PathParams: &api.GetAPIV1AssetsAssetIDReportsPath{
			AssetID: assetID,
		},
		Header: &api.GetAPIV1AssetsAssetIDReportsHeaders{
			XXBOWAPIVersion: api.GetAPIV1AssetsAssetIDReportsHeaderXXBOWAPIVersionN20260201,
		},
	}

	if opts != nil {
		reqOpts.Query = &api.GetAPIV1AssetsAssetIDReportsQuery{}
		if opts.Limit > 0 {
			reqOpts.Query.Limit = &opts.Limit
		}
		if opts.After != "" {
			reqOpts.Query.After = &opts.After
		}
	}

	resp, err := s.client.raw.GetAPIV1AssetsAssetIDReports(ctx, reqOpts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return reportsPageFromResponse(resp), nil
}

// AllByAsset returns an iterator over all reports for an asset.
// Use this for automatic pagination:
//
//	for report, err := range client.Reports.AllByAsset(ctx, assetID, nil) {
//	    if err != nil {
//	        return err
//	    }
//	    fmt.Println(report.ID)
//	}
func (s *ReportsService) AllByAsset(ctx context.Context, assetID string, opts *ListOptions) iter.Seq2[ReportListItem, error] {
	return paginate(ctx, opts, func(ctx context.Context, pageOpts *ListOptions) (*Page[ReportListItem], error) {
		return s.ListByAsset(ctx, assetID, pageOpts)
	})
}

// Conversion functions from generated types to domain types

func reportSummaryFromResponse(r *api.GetAPIV1ReportsReportIDSummaryResponse) *ReportSummary {
	return &ReportSummary{
		Markdown: r.Markdown,
	}
}

func reportsPageFromResponse(r *api.GetAPIV1AssetsAssetIDReportsResponse) *Page[ReportListItem] {
	items := make([]ReportListItem, 0, len(r.Items))
	for _, item := range r.Items {
		items = append(items, ReportListItem{
			ID:        item.ID,
			Version:   int64(item.Version),
			CreatedAt: item.CreatedAt,
		})
	}

	return &Page[ReportListItem]{
		Items: items,
		PageInfo: PageInfo{
			NextCursor: r.NextCursor,
			HasMore:    r.NextCursor != nil && *r.NextCursor != "",
		},
	}
}
