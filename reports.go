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

// Get downloads a report as PDF bytes by ID.
// The returned bytes are the raw PDF file content.
func (s *ReportsService) Get(ctx context.Context, id string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v1/reports/%s", s.client.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.apiKey)
	req.Header.Set("X-XBOW-API-Version", APIVersion)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &Error{
			Code:    fmt.Sprintf("HTTP_%d", resp.StatusCode),
			Message: string(body),
		}
	}

	return body, nil
}

// GetSummary retrieves the markdown summary of a report by ID.
func (s *ReportsService) GetSummary(ctx context.Context, id string) (*ReportSummary, error) {
	opts := &api.GetAPIV1ReportsReportIDSummaryRequestOptions{
		PathParams: &api.GetAPIV1ReportsReportIDSummaryPath{
			ReportID: id,
		},
		Header: &api.GetAPIV1ReportsReportIDSummaryHeaders{
			XXBOWAPIVersion: api.GetAPIV1ReportsReportIDSummaryHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.GetAPIV1ReportsReportIDSummary(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return reportSummaryFromResponse(resp), nil
}

// ListByAsset returns a page of reports for an asset.
func (s *ReportsService) ListByAsset(ctx context.Context, assetID string, opts *ListOptions) (*Page[ReportListItem], error) {
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

	resp, err := s.client.raw.GetAPIV1AssetsAssetIDReports(ctx, reqOpts, s.client.authEditor())
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
