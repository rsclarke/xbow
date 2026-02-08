package xbow

import (
	"context"
	"encoding/json"
	"iter"
	"time"

	"github.com/rsclarke/xbow/internal/api"
)

// AssessmentsService handles assessment-related API calls.
type AssessmentsService struct {
	client *Client
}

// Get retrieves an assessment by ID.
func (s *AssessmentsService) Get(ctx context.Context, id string) (*Assessment, error) {
	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	opts := &api.GetAPIV1AssessmentsAssessmentIDRequestOptions{
		PathParams: &api.GetAPIV1AssessmentsAssessmentIDPath{
			AssessmentID: id,
		},
		Header: &api.GetAPIV1AssessmentsAssessmentIDHeaders{
			XXBOWAPIVersion: api.N20260201,
		},
	}

	resp, err := s.client.raw.GetAPIV1AssessmentsAssessmentID(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return assessmentFromGetResponse(resp), nil
}

// CreateAssessmentRequest specifies the parameters for creating an assessment.
type CreateAssessmentRequest struct {
	AttackCredits int64
	Objective     *string
}

// Create requests a new assessment for an asset.
func (s *AssessmentsService) Create(ctx context.Context, assetID string, req *CreateAssessmentRequest) (*Assessment, error) {
	if req == nil {
		return nil, &Error{Code: "ERR_INVALID_REQUEST", Message: "CreateAssessmentRequest cannot be nil"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	opts := &api.PostAPIV1AssetsAssetIDAssessmentsRequestOptions{
		PathParams: &api.PostAPIV1AssetsAssetIDAssessmentsPath{
			AssetID: assetID,
		},
		Header: &api.PostAPIV1AssetsAssetIDAssessmentsHeaders{
			XXBOWAPIVersion: api.PostAPIV1AssetsAssetIDAssessmentsHeaderXXBOWAPIVersionN20260201,
		},
		Body: &api.PostAPIV1AssetsAssetIDAssessmentsBody{
			AttackCredits: int(req.AttackCredits),
			Objective:     req.Objective,
		},
	}

	resp, err := s.client.raw.PostAPIV1AssetsAssetIDAssessments(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return assessmentFromCreateResponse(resp), nil
}

// ListByAsset returns a page of assessments for an asset.
func (s *AssessmentsService) ListByAsset(ctx context.Context, assetID string, opts *ListOptions) (*Page[AssessmentListItem], error) {
	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	reqOpts := &api.GetAPIV1AssetsAssetIDAssessmentsRequestOptions{
		PathParams: &api.GetAPIV1AssetsAssetIDAssessmentsPath{
			AssetID: assetID,
		},
		Header: &api.GetAPIV1AssetsAssetIDAssessmentsHeaders{
			XXBOWAPIVersion: api.GetAPIV1AssetsAssetIDAssessmentsHeaderXXBOWAPIVersionN20260201,
		},
	}

	if opts != nil {
		reqOpts.Query = &api.GetAPIV1AssetsAssetIDAssessmentsQuery{}
		if opts.Limit > 0 {
			reqOpts.Query.Limit = &opts.Limit
		}
		if opts.After != "" {
			reqOpts.Query.After = &opts.After
		}
	}

	resp, err := s.client.raw.GetAPIV1AssetsAssetIDAssessments(ctx, reqOpts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return assessmentsPageFromResponse(resp), nil
}

// AllByAsset returns an iterator over all assessments for an asset.
// Use this for automatic pagination:
//
//	for assessment, err := range client.Assessments.AllByAsset(ctx, assetID, nil) {
//	    if err != nil {
//	        return err
//	    }
//	    fmt.Println(assessment.Name)
//	}
func (s *AssessmentsService) AllByAsset(ctx context.Context, assetID string, opts *ListOptions) iter.Seq2[AssessmentListItem, error] {
	return paginate(ctx, opts, func(ctx context.Context, pageOpts *ListOptions) (*Page[AssessmentListItem], error) {
		return s.ListByAsset(ctx, assetID, pageOpts)
	})
}

// Cancel cancels a running assessment.
func (s *AssessmentsService) Cancel(ctx context.Context, id string) (*Assessment, error) {
	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	opts := &api.PostAPIV1AssessmentsAssessmentIDCancelRequestOptions{
		PathParams: &api.PostAPIV1AssessmentsAssessmentIDCancelPath{
			AssessmentID: id,
		},
		Header: &api.PostAPIV1AssessmentsAssessmentIDCancelHeaders{
			XXBOWAPIVersion: api.PostAPIV1AssessmentsAssessmentIDCancelHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.PostAPIV1AssessmentsAssessmentIDCancel(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return assessmentFromCancelResponse(resp), nil
}

// Pause pauses a running assessment.
func (s *AssessmentsService) Pause(ctx context.Context, id string) (*Assessment, error) {
	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	opts := &api.PostAPIV1AssessmentsAssessmentIDPauseRequestOptions{
		PathParams: &api.PostAPIV1AssessmentsAssessmentIDPausePath{
			AssessmentID: id,
		},
		Header: &api.PostAPIV1AssessmentsAssessmentIDPauseHeaders{
			XXBOWAPIVersion: api.PostAPIV1AssessmentsAssessmentIDPauseHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.PostAPIV1AssessmentsAssessmentIDPause(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return assessmentFromPauseResponse(resp), nil
}

// Resume resumes a paused assessment.
func (s *AssessmentsService) Resume(ctx context.Context, id string) (*Assessment, error) {
	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	opts := &api.PostAPIV1AssessmentsAssessmentIDResumeRequestOptions{
		PathParams: &api.PostAPIV1AssessmentsAssessmentIDResumePath{
			AssessmentID: id,
		},
		Header: &api.PostAPIV1AssessmentsAssessmentIDResumeHeaders{
			XXBOWAPIVersion: api.PostAPIV1AssessmentsAssessmentIDResumeHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.PostAPIV1AssessmentsAssessmentIDResume(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return assessmentFromResumeResponse(resp), nil
}

// Conversion functions from generated types to domain types

func assessmentFromGetResponse(r *api.GetAPIV1AssessmentsAssessmentIDResponse) *Assessment {
	return &Assessment{
		ID:             r.ID,
		Name:           r.Name,
		AssetID:        r.AssetID,
		OrganizationID: r.OrganizationID,
		State:          AssessmentState(r.State),
		Progress:       float64(r.Progress),
		AttackCredits:  int64(r.AttackCredits),
		RecentEvents: convertRecentEvents(r.RecentEvents, func(e api.GetAPIV1AssessmentsAssessmentID_Response_RecentEvents_Item) rawUnion {
			if e.GetAPIV1AssessmentsAssessmentID_Response_RecentEvents_OneOf == nil {
				return nil
			}
			return e.GetAPIV1AssessmentsAssessmentID_Response_RecentEvents_OneOf
		}),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func assessmentFromCreateResponse(r *api.PostAPIV1AssetsAssetIDAssessmentsResponse) *Assessment {
	return &Assessment{
		ID:             r.ID,
		Name:           r.Name,
		AssetID:        r.AssetID,
		OrganizationID: r.OrganizationID,
		State:          AssessmentState(r.State),
		Progress:       float64(r.Progress),
		AttackCredits:  int64(r.AttackCredits),
		RecentEvents: convertRecentEvents(r.RecentEvents, func(e api.PostAPIV1AssetsAssetIDAssessments_Response_RecentEvents_Item) rawUnion {
			if e.PostAPIV1AssetsAssetIDAssessments_Response_RecentEvents_OneOf == nil {
				return nil
			}
			return e.PostAPIV1AssetsAssetIDAssessments_Response_RecentEvents_OneOf
		}),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func assessmentFromCancelResponse(r *api.PostAPIV1AssessmentsAssessmentIDCancelResponse) *Assessment {
	return &Assessment{
		ID:             r.ID,
		Name:           r.Name,
		AssetID:        r.AssetID,
		OrganizationID: r.OrganizationID,
		State:          AssessmentState(r.State),
		Progress:       float64(r.Progress),
		AttackCredits:  int64(r.AttackCredits),
		RecentEvents: convertRecentEvents(r.RecentEvents, func(e api.PostAPIV1AssessmentsAssessmentIDCancel_Response_RecentEvents_Item) rawUnion {
			if e.PostAPIV1AssessmentsAssessmentIDCancel_Response_RecentEvents_OneOf == nil {
				return nil
			}
			return e.PostAPIV1AssessmentsAssessmentIDCancel_Response_RecentEvents_OneOf
		}),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func assessmentFromPauseResponse(r *api.PostAPIV1AssessmentsAssessmentIDPauseResponse) *Assessment {
	return &Assessment{
		ID:             r.ID,
		Name:           r.Name,
		AssetID:        r.AssetID,
		OrganizationID: r.OrganizationID,
		State:          AssessmentState(r.State),
		Progress:       float64(r.Progress),
		AttackCredits:  int64(r.AttackCredits),
		RecentEvents: convertRecentEvents(r.RecentEvents, func(e api.PostAPIV1AssessmentsAssessmentIDPause_Response_RecentEvents_Item) rawUnion {
			if e.PostAPIV1AssessmentsAssessmentIDPause_Response_RecentEvents_OneOf == nil {
				return nil
			}
			return e.PostAPIV1AssessmentsAssessmentIDPause_Response_RecentEvents_OneOf
		}),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func assessmentFromResumeResponse(r *api.PostAPIV1AssessmentsAssessmentIDResumeResponse) *Assessment {
	return &Assessment{
		ID:             r.ID,
		Name:           r.Name,
		AssetID:        r.AssetID,
		OrganizationID: r.OrganizationID,
		State:          AssessmentState(r.State),
		Progress:       float64(r.Progress),
		AttackCredits:  int64(r.AttackCredits),
		RecentEvents: convertRecentEvents(r.RecentEvents, func(e api.PostAPIV1AssessmentsAssessmentIDResume_Response_RecentEvents_Item) rawUnion {
			if e.PostAPIV1AssessmentsAssessmentIDResume_Response_RecentEvents_OneOf == nil {
				return nil
			}
			return e.PostAPIV1AssessmentsAssessmentIDResume_Response_RecentEvents_OneOf
		}),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func assessmentsPageFromResponse(r *api.GetAPIV1AssetsAssetIDAssessmentsResponse) *Page[AssessmentListItem] {
	items := make([]AssessmentListItem, 0, len(r.Items))
	for _, item := range r.Items {
		items = append(items, AssessmentListItem{
			ID:        item.ID,
			Name:      item.Name,
			State:     AssessmentState(item.State),
			Progress:  float64(item.Progress),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return &Page[AssessmentListItem]{
		Items: items,
		PageInfo: PageInfo{
			NextCursor: r.NextCursor,
			HasMore:    r.NextCursor != nil && *r.NextCursor != "",
		},
	}
}

type rawUnion interface {
	Raw() json.RawMessage
}

type recentEventJSON struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason,omitempty"`
}

func convertRecentEvents[Item any](items []Item, getOneOf func(Item) rawUnion) []AssessmentEvent {
	result := make([]AssessmentEvent, 0, len(items))
	for _, item := range items {
		oneOf := getOneOf(item)
		if oneOf == nil {
			continue
		}
		var ev recentEventJSON
		if err := json.Unmarshal(oneOf.Raw(), &ev); err != nil {
			continue
		}
		result = append(result, AssessmentEvent(ev))
	}
	return result
}

func ptrValue[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
