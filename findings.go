package xbow

import (
	"context"
	"iter"

	"github.com/rsclarke/xbow/internal/api"
)

// FindingsService handles finding-related API calls.
type FindingsService struct {
	client *Client
}

// Get retrieves a finding by ID.
func (s *FindingsService) Get(ctx context.Context, id string) (*Finding, error) {
	opts := &api.GetAPIV1FindingsFindingIDRequestOptions{
		PathParams: &api.GetAPIV1FindingsFindingIDPath{
			FindingID: id,
		},
		Header: &api.GetAPIV1FindingsFindingIDHeaders{
			XXBOWAPIVersion: api.GetAPIV1FindingsFindingIDHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.GetAPIV1FindingsFindingID(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return findingFromGetResponse(resp), nil
}

// ListByAsset returns a page of findings for an asset.
func (s *FindingsService) ListByAsset(ctx context.Context, assetID string, opts *ListOptions) (*Page[FindingListItem], error) {
	reqOpts := &api.GetAPIV1AssetsAssetIDFindingsRequestOptions{
		PathParams: &api.GetAPIV1AssetsAssetIDFindingsPath{
			AssetID: assetID,
		},
		Header: &api.GetAPIV1AssetsAssetIDFindingsHeaders{
			XXBOWAPIVersion: api.GetAPIV1AssetsAssetIDFindingsHeaderXXBOWAPIVersionN20260201,
		},
	}

	if opts != nil {
		reqOpts.Query = &api.GetAPIV1AssetsAssetIDFindingsQuery{}
		if opts.Limit > 0 {
			reqOpts.Query.Limit = &opts.Limit
		}
		if opts.After != "" {
			reqOpts.Query.After = &opts.After
		}
	}

	resp, err := s.client.raw.GetAPIV1AssetsAssetIDFindings(ctx, reqOpts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return findingsPageFromResponse(resp), nil
}

// AllByAsset returns an iterator over all findings for an asset.
// Use this for automatic pagination:
//
//	for finding, err := range client.Findings.AllByAsset(ctx, assetID, nil) {
//	    if err != nil {
//	        return err
//	    }
//	    fmt.Println(finding.Name)
//	}
func (s *FindingsService) AllByAsset(ctx context.Context, assetID string, opts *ListOptions) iter.Seq2[FindingListItem, error] {
	return paginate(ctx, opts, func(ctx context.Context, pageOpts *ListOptions) (*Page[FindingListItem], error) {
		return s.ListByAsset(ctx, assetID, pageOpts)
	})
}

// VerifyFix requests verification that a finding has been fixed.
// This triggers a targeted assessment to verify the vulnerability has been mitigated.
// Returns the assessment created for the verification.
func (s *FindingsService) VerifyFix(ctx context.Context, id string) (*Assessment, error) {
	opts := &api.PostAPIV1FindingsFindingIDVerifyFixRequestOptions{
		PathParams: &api.PostAPIV1FindingsFindingIDVerifyFixPath{
			FindingID: id,
		},
		Header: &api.PostAPIV1FindingsFindingIDVerifyFixHeaders{
			XXBOWAPIVersion: api.PostAPIV1FindingsFindingIDVerifyFixHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.PostAPIV1FindingsFindingIDVerifyFix(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return assessmentFromVerifyFixResponse(resp), nil
}

// Conversion functions from generated types to domain types

func findingFromGetResponse(r *api.GetAPIV1FindingsFindingIDResponse) *Finding {
	return &Finding{
		ID:          r.ID,
		Name:        r.Name,
		Severity:    FindingSeverity(r.Severity),
		State:       FindingState(r.State),
		Summary:     r.Summary,
		Impact:      r.Impact,
		Mitigations: r.Mitigations,
		Recipe:      r.Recipe,
		Evidence:    r.Evidence,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func findingsPageFromResponse(r *api.GetAPIV1AssetsAssetIDFindingsResponse) *Page[FindingListItem] {
	items := make([]FindingListItem, 0, len(r.Items))
	for _, item := range r.Items {
		items = append(items, FindingListItem{
			ID:        item.ID,
			Name:      item.Name,
			Severity:  FindingSeverity(item.Severity),
			State:     FindingState(item.State),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return &Page[FindingListItem]{
		Items: items,
		PageInfo: PageInfo{
			NextCursor: r.NextCursor,
			HasMore:    r.NextCursor != nil && *r.NextCursor != "",
		},
	}
}

func assessmentFromVerifyFixResponse(r *api.PostAPIV1FindingsFindingIDVerifyFixResponse) *Assessment {
	return &Assessment{
		ID:             r.ID,
		Name:           r.Name,
		AssetID:        r.AssetID,
		OrganizationID: r.OrganizationID,
		State:          AssessmentState(r.State),
		Progress:       float64(r.Progress),
		AttackCredits:  int64(r.AttackCredits),
		RecentEvents:   convertRecentEventsFromVerifyFix(r.RecentEvents),
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

func convertRecentEventsFromVerifyFix(events api.PostAPIV1FindingsFindingIDVerifyFix_Response_RecentEvents) []AssessmentEvent {
	result := make([]AssessmentEvent, 0, len(events))
	for _, e := range events {
		if e.PostAPIV1FindingsFindingIDVerifyFix_Response_RecentEvents_OneOf == nil {
			continue
		}
		oneOf := e.PostAPIV1FindingsFindingIDVerifyFix_Response_RecentEvents_OneOf

		// Try paused event
		if v, err := oneOf.AsPostAPIV1FindingsFindingIDVerifyFix_Response_RecentEvents_OneOf_0(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}

		// Try auto-paused event (with reason)
		if v, err := oneOf.AsPostAPIV1FindingsFindingIDVerifyFix_Response_RecentEvents_OneOf_1(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
				Reason:    string(v.Reason),
			})
			continue
		}

		// Try resumed event
		if v, err := oneOf.AsPostAPIV1FindingsFindingIDVerifyFix_Response_RecentEvents_OneOf_2(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}
	}
	return result
}
