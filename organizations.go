package xbow

import (
	"context"
	"iter"

	"github.com/doordash-oss/oapi-codegen-dd/v3/pkg/runtime"
	"github.com/rsclarke/xbow/internal/api"
)

// OrganizationsService handles organization-related API calls.
type OrganizationsService struct {
	client *Client
}

// Get retrieves an organization by ID.
func (s *OrganizationsService) Get(ctx context.Context, id string) (*Organization, error) {
	opts := &api.GetAPIV1OrganizationsOrganizationIDRequestOptions{
		PathParams: &api.GetAPIV1OrganizationsOrganizationIDPath{
			OrganizationID: id,
		},
		Header: &api.GetAPIV1OrganizationsOrganizationIDHeaders{
			XXBOWAPIVersion: api.GetAPIV1OrganizationsOrganizationIDHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.GetAPIV1OrganizationsOrganizationID(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return organizationFromGetResponse(resp), nil
}

// UpdateOrganizationRequest specifies the parameters for updating an organization.
// Both Name and ExternalID are required. Set ExternalID to nil to clear it.
type UpdateOrganizationRequest struct {
	Name       string  // Required
	ExternalID *string // Required (use nil to set null, or pointer to string for a value)
}

// Update updates an organization.
func (s *OrganizationsService) Update(ctx context.Context, id string, req *UpdateOrganizationRequest) (*Organization, error) {
	if req == nil {
		return nil, &Error{Code: "ERR_INVALID_REQUEST", Message: "UpdateOrganizationRequest cannot be nil"}
	}

	externalID := ""
	if req.ExternalID != nil {
		externalID = *req.ExternalID
	}

	opts := &api.PutAPIV1OrganizationsOrganizationIDRequestOptions{
		PathParams: &api.PutAPIV1OrganizationsOrganizationIDPath{
			OrganizationID: id,
		},
		Header: &api.PutAPIV1OrganizationsOrganizationIDHeaders{
			XXBOWAPIVersion: api.PutAPIV1OrganizationsOrganizationIDHeaderXXBOWAPIVersionN20260201,
		},
		Body: &api.PutAPIV1OrganizationsOrganizationIDBody{
			Name:       req.Name,
			ExternalID: externalID,
		},
	}

	resp, err := s.client.raw.PutAPIV1OrganizationsOrganizationID(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return organizationFromPutResponse(resp), nil
}

// CreateOrganizationRequest specifies the parameters for creating an organization.
// All fields are required per the API specification.
type CreateOrganizationRequest struct {
	Name       string               // Required - name of the organization
	ExternalID *string              // Required (use nil to set null, or pointer to string for a value)
	Members    []OrganizationMember // Required - at least one member
}

// Create creates a new organization in an integration.
func (s *OrganizationsService) Create(ctx context.Context, integrationID string, req *CreateOrganizationRequest) (*Organization, error) {
	if req == nil {
		return nil, &Error{Code: "ERR_INVALID_REQUEST", Message: "CreateOrganizationRequest cannot be nil"}
	}

	if len(req.Members) == 0 {
		return nil, &Error{Code: "ERR_INVALID_REQUEST", Message: "at least one member is required"}
	}

	externalID := ""
	if req.ExternalID != nil {
		externalID = *req.ExternalID
	}

	members := make(api.PostAPIV1IntegrationsIntegrationIDOrganizationsBody_Members, 0, len(req.Members))
	for _, m := range req.Members {
		members = append(members, api.PostAPIV1IntegrationsIntegrationIDOrganizationsBody_Members_Item{
			Email: runtime.Email(m.Email),
			Name:  m.Name,
		})
	}

	opts := &api.PostAPIV1IntegrationsIntegrationIDOrganizationsRequestOptions{
		PathParams: &api.PostAPIV1IntegrationsIntegrationIDOrganizationsPath{
			IntegrationID: integrationID,
		},
		Header: &api.PostAPIV1IntegrationsIntegrationIDOrganizationsHeaders{
			XXBOWAPIVersion: api.PostAPIV1IntegrationsIntegrationIDOrganizationsHeaderXXBOWAPIVersionN20260201,
		},
		Body: &api.PostAPIV1IntegrationsIntegrationIDOrganizationsBody{
			Name:       req.Name,
			ExternalID: externalID,
			Members:    members,
		},
	}

	resp, err := s.client.raw.PostAPIV1IntegrationsIntegrationIDOrganizations(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return organizationFromCreateResponse(resp), nil
}

// ListByIntegration returns a page of organizations for an integration.
func (s *OrganizationsService) ListByIntegration(ctx context.Context, integrationID string, opts *ListOptions) (*Page[OrganizationListItem], error) {
	reqOpts := &api.GetAPIV1IntegrationsIntegrationIDOrganizationsRequestOptions{
		PathParams: &api.GetAPIV1IntegrationsIntegrationIDOrganizationsPath{
			IntegrationID: integrationID,
		},
		Header: &api.GetAPIV1IntegrationsIntegrationIDOrganizationsHeaders{
			XXBOWAPIVersion: api.GetAPIV1IntegrationsIntegrationIDOrganizationsHeaderXXBOWAPIVersionN20260201,
		},
	}

	if opts != nil {
		reqOpts.Query = &api.GetAPIV1IntegrationsIntegrationIDOrganizationsQuery{}
		if opts.Limit > 0 {
			reqOpts.Query.Limit = &opts.Limit
		}
		if opts.After != "" {
			reqOpts.Query.After = &opts.After
		}
	}

	resp, err := s.client.raw.GetAPIV1IntegrationsIntegrationIDOrganizations(ctx, reqOpts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return organizationsPageFromResponse(resp), nil
}

// AllByIntegration returns an iterator over all organizations for an integration.
func (s *OrganizationsService) AllByIntegration(ctx context.Context, integrationID string, opts *ListOptions) iter.Seq2[OrganizationListItem, error] {
	return paginate(ctx, opts, func(ctx context.Context, pageOpts *ListOptions) (*Page[OrganizationListItem], error) {
		return s.ListByIntegration(ctx, integrationID, pageOpts)
	})
}

// CreateKeyRequest specifies the parameters for creating an organization API key.
type CreateKeyRequest struct {
	Name          string
	ExpiresInDays *int
}

// CreateKey creates a new API key for an organization.
func (s *OrganizationsService) CreateKey(ctx context.Context, organizationID string, req *CreateKeyRequest) (*OrganizationAPIKey, error) {
	if req == nil {
		return nil, &Error{Code: "ERR_INVALID_REQUEST", Message: "CreateKeyRequest cannot be nil"}
	}

	opts := &api.PostAPIV1OrganizationsOrganizationIDKeysRequestOptions{
		PathParams: &api.PostAPIV1OrganizationsOrganizationIDKeysPath{
			OrganizationID: organizationID,
		},
		Header: &api.PostAPIV1OrganizationsOrganizationIDKeysHeaders{
			XXBOWAPIVersion: api.PostAPIV1OrganizationsOrganizationIDKeysHeaderXXBOWAPIVersionN20260201,
		},
		Body: &api.PostAPIV1OrganizationsOrganizationIDKeysBody{
			Name:          req.Name,
			ExpiresInDays: req.ExpiresInDays,
		},
	}

	resp, err := s.client.raw.PostAPIV1OrganizationsOrganizationIDKeys(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return apiKeyFromResponse(resp), nil
}

// RevokeKey revokes an organization API key.
func (s *OrganizationsService) RevokeKey(ctx context.Context, keyID string) error {
	opts := &api.DeleteAPIV1KeysKeyIDRequestOptions{
		PathParams: &api.DeleteAPIV1KeysKeyIDPath{
			KeyID: keyID,
		},
		Header: &api.DeleteAPIV1KeysKeyIDHeaders{
			XXBOWAPIVersion: api.DeleteAPIV1KeysKeyIDHeaderXXBOWAPIVersionN20260201,
		},
	}

	_, err := s.client.raw.DeleteAPIV1KeysKeyID(ctx, opts, s.client.authEditor())
	if err != nil {
		return wrapError(err)
	}

	return nil
}

// Conversion functions from generated types to domain types

func organizationFromGetResponse(r *api.GetAPIV1OrganizationsOrganizationIDResponse) *Organization {
	return &Organization{
		ID:         r.ID,
		Name:       r.Name,
		ExternalID: strPtrFromNullable(r.ExternalID),
		State:      OrganizationState(r.State),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func organizationFromPutResponse(r *api.PutAPIV1OrganizationsOrganizationIDResponse) *Organization {
	return &Organization{
		ID:         r.ID,
		Name:       r.Name,
		ExternalID: strPtrFromNullable(r.ExternalID),
		State:      OrganizationState(r.State),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func organizationFromCreateResponse(r *api.PostAPIV1IntegrationsIntegrationIDOrganizationsResponse) *Organization {
	return &Organization{
		ID:         r.ID,
		Name:       r.Name,
		ExternalID: strPtrFromNullable(r.ExternalID),
		State:      OrganizationState(r.State),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func organizationsPageFromResponse(r *api.GetAPIV1IntegrationsIntegrationIDOrganizationsResponse) *Page[OrganizationListItem] {
	items := make([]OrganizationListItem, 0, len(r.Items))
	for _, item := range r.Items {
		items = append(items, OrganizationListItem{
			ID:         item.ID,
			Name:       item.Name,
			ExternalID: strPtrFromNullable(item.ExternalID),
			State:      OrganizationState(item.State),
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		})
	}

	return &Page[OrganizationListItem]{
		Items: items,
		PageInfo: PageInfo{
			NextCursor: r.NextCursor,
			HasMore:    r.NextCursor != nil && *r.NextCursor != "",
		},
	}
}

func apiKeyFromResponse(r *api.PostAPIV1OrganizationsOrganizationIDKeysResponse) *OrganizationAPIKey {
	return &OrganizationAPIKey{
		ID:        r.ID,
		Name:      r.Name,
		Key:       r.Key,
		ExpiresAt: timePtrFromNullable(r.ExpiresAt),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
