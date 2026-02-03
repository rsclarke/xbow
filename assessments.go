package xbow

import (
	"context"
	"iter"

	"github.com/rsclarke/xbow/internal/api"
)

// AssessmentsService handles assessment-related API calls.
type AssessmentsService struct {
	client *Client
}

// Get retrieves an assessment by ID.
func (s *AssessmentsService) Get(ctx context.Context, id string) (*Assessment, error) {
	opts := &api.GetAPIV1AssessmentsAssessmentIDRequestOptions{
		PathParams: &api.GetAPIV1AssessmentsAssessmentIDPath{
			AssessmentID: id,
		},
		Header: &api.GetAPIV1AssessmentsAssessmentIDHeaders{
			XXBOWAPIVersion: api.N20260201,
		},
	}

	resp, err := s.client.raw.GetAPIV1AssessmentsAssessmentID(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return assessmentFromGetResponse(resp), nil
}

// CreateAssessmentRequest specifies the parameters for creating an assessment.
type CreateAssessmentRequest struct {
	AttackCredits int
	Objective     *string
}

// Create requests a new assessment for an asset.
func (s *AssessmentsService) Create(ctx context.Context, assetID string, req *CreateAssessmentRequest) (*Assessment, error) {
	if req == nil {
		return nil, &Error{Code: "ERR_INVALID_REQUEST", Message: "CreateAssessmentRequest cannot be nil"}
	}

	opts := &api.PostAPIV1AssetsAssetIDAssessmentsRequestOptions{
		PathParams: &api.PostAPIV1AssetsAssetIDAssessmentsPath{
			AssetID: assetID,
		},
		Header: &api.PostAPIV1AssetsAssetIDAssessmentsHeaders{
			XXBOWAPIVersion: api.PostAPIV1AssetsAssetIDAssessmentsHeaderXXBOWAPIVersionN20260201,
		},
		Body: &api.PostAPIV1AssetsAssetIDAssessmentsBody{
			AttackCredits: req.AttackCredits,
			Objective:     req.Objective,
		},
	}

	resp, err := s.client.raw.PostAPIV1AssetsAssetIDAssessments(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return assessmentFromCreateResponse(resp), nil
}

// ListByAsset returns a page of assessments for an asset.
func (s *AssessmentsService) ListByAsset(ctx context.Context, assetID string, opts *ListOptions) (*Page[Assessment], error) {
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

	resp, err := s.client.raw.GetAPIV1AssetsAssetIDAssessments(ctx, reqOpts, s.client.authEditor())
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
func (s *AssessmentsService) AllByAsset(ctx context.Context, assetID string, opts *ListOptions) iter.Seq2[Assessment, error] {
	return paginate(ctx, opts, func(ctx context.Context, pageOpts *ListOptions) (*Page[Assessment], error) {
		return s.ListByAsset(ctx, assetID, pageOpts)
	})
}

// Cancel cancels a running assessment.
func (s *AssessmentsService) Cancel(ctx context.Context, id string) (*Assessment, error) {
	opts := &api.PostAPIV1AssessmentsAssessmentIDCancelRequestOptions{
		PathParams: &api.PostAPIV1AssessmentsAssessmentIDCancelPath{
			AssessmentID: id,
		},
		Header: &api.PostAPIV1AssessmentsAssessmentIDCancelHeaders{
			XXBOWAPIVersion: api.PostAPIV1AssessmentsAssessmentIDCancelHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.PostAPIV1AssessmentsAssessmentIDCancel(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return assessmentFromCancelResponse(resp), nil
}

// Pause pauses a running assessment.
func (s *AssessmentsService) Pause(ctx context.Context, id string) (*Assessment, error) {
	opts := &api.PostAPIV1AssessmentsAssessmentIDPauseRequestOptions{
		PathParams: &api.PostAPIV1AssessmentsAssessmentIDPausePath{
			AssessmentID: id,
		},
		Header: &api.PostAPIV1AssessmentsAssessmentIDPauseHeaders{
			XXBOWAPIVersion: api.PostAPIV1AssessmentsAssessmentIDPauseHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.PostAPIV1AssessmentsAssessmentIDPause(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return assessmentFromPauseResponse(resp), nil
}

// Resume resumes a paused assessment.
func (s *AssessmentsService) Resume(ctx context.Context, id string) (*Assessment, error) {
	opts := &api.PostAPIV1AssessmentsAssessmentIDResumeRequestOptions{
		PathParams: &api.PostAPIV1AssessmentsAssessmentIDResumePath{
			AssessmentID: id,
		},
		Header: &api.PostAPIV1AssessmentsAssessmentIDResumeHeaders{
			XXBOWAPIVersion: api.PostAPIV1AssessmentsAssessmentIDResumeHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.PostAPIV1AssessmentsAssessmentIDResume(ctx, opts, s.client.authEditor())
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
		AttackCredits:  r.AttackCredits,
		RecentEvents:   convertRecentEventsFromGet(r.RecentEvents),
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
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
		AttackCredits:  r.AttackCredits,
		RecentEvents:   convertRecentEventsFromCreate(r.RecentEvents),
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
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
		AttackCredits:  r.AttackCredits,
		RecentEvents:   convertRecentEventsFromCancel(r.RecentEvents),
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
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
		AttackCredits:  r.AttackCredits,
		RecentEvents:   convertRecentEventsFromPause(r.RecentEvents),
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
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
		AttackCredits:  r.AttackCredits,
		RecentEvents:   convertRecentEventsFromResume(r.RecentEvents),
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

func assessmentsPageFromResponse(r *api.GetAPIV1AssetsAssetIDAssessmentsResponse) *Page[Assessment] {
	items := make([]Assessment, 0, len(r.Items))
	for _, item := range r.Items {
		items = append(items, Assessment{
			ID:        item.ID,
			Name:      item.Name,
			State:     AssessmentState(item.State),
			Progress:  float64(item.Progress),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return &Page[Assessment]{
		Items: items,
		PageInfo: PageInfo{
			NextCursor: ptrValue(r.NextCursor),
			HasMore:    r.NextCursor != nil && *r.NextCursor != "",
		},
	}
}

func convertRecentEventsFromGet(events api.GetAPIV1AssessmentsAssessmentID_Response_RecentEvents) []AssessmentEvent {
	result := make([]AssessmentEvent, 0, len(events))
	for _, e := range events {
		if e.GetAPIV1AssessmentsAssessmentID_Response_RecentEvents_OneOf == nil {
			continue
		}
		oneOf := e.GetAPIV1AssessmentsAssessmentID_Response_RecentEvents_OneOf

		// Try paused event
		if v, err := oneOf.AsGetAPIV1AssessmentsAssessmentID_Response_RecentEvents_OneOf_0(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}

		// Try auto-paused event (with reason)
		if v, err := oneOf.AsGetAPIV1AssessmentsAssessmentID_Response_RecentEvents_OneOf_1(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
				Reason:    string(v.Reason),
			})
			continue
		}

		// Try resumed event
		if v, err := oneOf.AsGetAPIV1AssessmentsAssessmentID_Response_RecentEvents_OneOf_2(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}
	}
	return result
}

func convertRecentEventsFromCreate(events api.PostAPIV1AssetsAssetIDAssessments_Response_RecentEvents) []AssessmentEvent {
	result := make([]AssessmentEvent, 0, len(events))
	for _, e := range events {
		if e.PostAPIV1AssetsAssetIDAssessments_Response_RecentEvents_OneOf == nil {
			continue
		}
		oneOf := e.PostAPIV1AssetsAssetIDAssessments_Response_RecentEvents_OneOf

		if v, err := oneOf.AsPostAPIV1AssetsAssetIDAssessments_Response_RecentEvents_OneOf_0(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}

		if v, err := oneOf.AsPostAPIV1AssetsAssetIDAssessments_Response_RecentEvents_OneOf_1(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
				Reason:    string(v.Reason),
			})
			continue
		}

		if v, err := oneOf.AsPostAPIV1AssetsAssetIDAssessments_Response_RecentEvents_OneOf_2(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}
	}
	return result
}

func convertRecentEventsFromCancel(events api.PostAPIV1AssessmentsAssessmentIDCancel_Response_RecentEvents) []AssessmentEvent {
	result := make([]AssessmentEvent, 0, len(events))
	for _, e := range events {
		if e.PostAPIV1AssessmentsAssessmentIDCancel_Response_RecentEvents_OneOf == nil {
			continue
		}
		oneOf := e.PostAPIV1AssessmentsAssessmentIDCancel_Response_RecentEvents_OneOf

		if v, err := oneOf.AsPostAPIV1AssessmentsAssessmentIDCancel_Response_RecentEvents_OneOf_0(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}

		if v, err := oneOf.AsPostAPIV1AssessmentsAssessmentIDCancel_Response_RecentEvents_OneOf_1(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
				Reason:    string(v.Reason),
			})
			continue
		}

		if v, err := oneOf.AsPostAPIV1AssessmentsAssessmentIDCancel_Response_RecentEvents_OneOf_2(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}
	}
	return result
}

func convertRecentEventsFromPause(events api.PostAPIV1AssessmentsAssessmentIDPause_Response_RecentEvents) []AssessmentEvent {
	result := make([]AssessmentEvent, 0, len(events))
	for _, e := range events {
		if e.PostAPIV1AssessmentsAssessmentIDPause_Response_RecentEvents_OneOf == nil {
			continue
		}
		oneOf := e.PostAPIV1AssessmentsAssessmentIDPause_Response_RecentEvents_OneOf

		if v, err := oneOf.AsPostAPIV1AssessmentsAssessmentIDPause_Response_RecentEvents_OneOf_0(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}

		if v, err := oneOf.AsPostAPIV1AssessmentsAssessmentIDPause_Response_RecentEvents_OneOf_1(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
				Reason:    string(v.Reason),
			})
			continue
		}

		if v, err := oneOf.AsPostAPIV1AssessmentsAssessmentIDPause_Response_RecentEvents_OneOf_2(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}
	}
	return result
}

func convertRecentEventsFromResume(events api.PostAPIV1AssessmentsAssessmentIDResume_Response_RecentEvents) []AssessmentEvent {
	result := make([]AssessmentEvent, 0, len(events))
	for _, e := range events {
		if e.PostAPIV1AssessmentsAssessmentIDResume_Response_RecentEvents_OneOf == nil {
			continue
		}
		oneOf := e.PostAPIV1AssessmentsAssessmentIDResume_Response_RecentEvents_OneOf

		if v, err := oneOf.AsPostAPIV1AssessmentsAssessmentIDResume_Response_RecentEvents_OneOf_0(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}

		if v, err := oneOf.AsPostAPIV1AssessmentsAssessmentIDResume_Response_RecentEvents_OneOf_1(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
				Reason:    string(v.Reason),
			})
			continue
		}

		if v, err := oneOf.AsPostAPIV1AssessmentsAssessmentIDResume_Response_RecentEvents_OneOf_2(); err == nil {
			result = append(result, AssessmentEvent{
				Name:      string(v.Name),
				Timestamp: v.Timestamp,
			})
			continue
		}
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
