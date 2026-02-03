package xbow

import (
	"context"
	"iter"

	"github.com/rsclarke/xbow/internal/api"
)

// AssetsService handles asset-related API calls.
type AssetsService struct {
	client *Client
}

// Get retrieves an asset by ID.
func (s *AssetsService) Get(ctx context.Context, id string) (*Asset, error) {
	opts := &api.GetAPIV1AssetsAssetIDRequestOptions{
		PathParams: &api.GetAPIV1AssetsAssetIDPath{
			AssetID: id,
		},
		Header: &api.GetAPIV1AssetsAssetIDHeaders{
			XXBOWAPIVersion: api.GetAPIV1AssetsAssetIDHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.GetAPIV1AssetsAssetID(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return assetFromGetResponse(resp), nil
}

// UpdateAssetRequest specifies the parameters for updating an asset.
type UpdateAssetRequest struct {
	Name                 string
	StartURL             string
	MaxRequestsPerSecond int
	Sku                  *string
}

// Update updates an asset.
func (s *AssetsService) Update(ctx context.Context, id string, req *UpdateAssetRequest) (*Asset, error) {
	if req == nil {
		return nil, &Error{Code: "ERR_INVALID_REQUEST", Message: "UpdateAssetRequest cannot be nil"}
	}

	opts := &api.PutAPIV1AssetsAssetIDRequestOptions{
		PathParams: &api.PutAPIV1AssetsAssetIDPath{
			AssetID: id,
		},
		Header: &api.PutAPIV1AssetsAssetIDHeaders{
			XXBOWAPIVersion: api.PutAPIV1AssetsAssetIDHeaderXXBOWAPIVersionN20260201,
		},
		Body: &api.PutAPIV1AssetsAssetIDBody{
			Name:                 req.Name,
			StartURL:             req.StartURL,
			MaxRequestsPerSecond: req.MaxRequestsPerSecond,
			Sku:                  req.Sku,
		},
	}

	resp, err := s.client.raw.PutAPIV1AssetsAssetID(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return assetFromPutResponse(resp), nil
}

// CreateAssetRequest specifies the parameters for creating an asset.
type CreateAssetRequest struct {
	Name string
	Sku  string
}

// Create creates a new asset in an organization.
func (s *AssetsService) Create(ctx context.Context, organizationID string, req *CreateAssetRequest) (*Asset, error) {
	if req == nil {
		return nil, &Error{Code: "ERR_INVALID_REQUEST", Message: "CreateAssetRequest cannot be nil"}
	}

	opts := &api.PostAPIV1OrganizationsOrganizationIDAssetsRequestOptions{
		PathParams: &api.PostAPIV1OrganizationsOrganizationIDAssetsPath{
			OrganizationID: organizationID,
		},
		Header: &api.PostAPIV1OrganizationsOrganizationIDAssetsHeaders{
			XXBOWAPIVersion: api.PostAPIV1OrganizationsOrganizationIDAssetsHeaderXXBOWAPIVersionN20260201,
		},
		Body: &api.PostAPIV1OrganizationsOrganizationIDAssetsBody{
			Name: req.Name,
			Sku:  req.Sku,
		},
	}

	resp, err := s.client.raw.PostAPIV1OrganizationsOrganizationIDAssets(ctx, opts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return assetFromCreateResponse(resp), nil
}

// ListByOrganization returns a page of assets for an organization.
func (s *AssetsService) ListByOrganization(ctx context.Context, organizationID string, opts *ListOptions) (*Page[AssetListItem], error) {
	reqOpts := &api.GetAPIV1OrganizationsOrganizationIDAssetsRequestOptions{
		PathParams: &api.GetAPIV1OrganizationsOrganizationIDAssetsPath{
			OrganizationID: organizationID,
		},
		Header: &api.GetAPIV1OrganizationsOrganizationIDAssetsHeaders{
			XXBOWAPIVersion: api.GetAPIV1OrganizationsOrganizationIDAssetsHeaderXXBOWAPIVersionN20260201,
		},
	}

	if opts != nil {
		reqOpts.Query = &api.GetAPIV1OrganizationsOrganizationIDAssetsQuery{}
		if opts.Limit > 0 {
			reqOpts.Query.Limit = &opts.Limit
		}
		if opts.After != "" {
			reqOpts.Query.After = &opts.After
		}
	}

	resp, err := s.client.raw.GetAPIV1OrganizationsOrganizationIDAssets(ctx, reqOpts, s.client.authEditor())
	if err != nil {
		return nil, wrapError(err)
	}

	return assetsPageFromResponse(resp), nil
}

// AllByOrganization returns an iterator over all assets for an organization.
func (s *AssetsService) AllByOrganization(ctx context.Context, organizationID string, opts *ListOptions) iter.Seq2[AssetListItem, error] {
	return paginate(ctx, opts, func(ctx context.Context, pageOpts *ListOptions) (*Page[AssetListItem], error) {
		return s.ListByOrganization(ctx, organizationID, pageOpts)
	})
}

// Conversion functions from generated types to domain types

func assetFromGetResponse(r *api.GetAPIV1AssetsAssetIDResponse) *Asset {
	return &Asset{
		ID:                   r.ID,
		Name:                 r.Name,
		OrganizationID:       r.OrganizationID,
		Lifecycle:            AssetLifecycle(r.Lifecycle),
		Sku:                  r.Sku,
		StartURL:             r.StartURL,
		MaxRequestsPerSecond: r.MaxRequestsPerSecond,
		CreatedAt:            r.CreatedAt,
		UpdatedAt:            r.UpdatedAt,
	}
}

func assetFromPutResponse(r *api.PutAPIV1AssetsAssetIDResponse) *Asset {
	return &Asset{
		ID:                   r.ID,
		Name:                 r.Name,
		OrganizationID:       r.OrganizationID,
		Lifecycle:            AssetLifecycle(r.Lifecycle),
		Sku:                  r.Sku,
		StartURL:             r.StartURL,
		MaxRequestsPerSecond: r.MaxRequestsPerSecond,
		CreatedAt:            r.CreatedAt,
		UpdatedAt:            r.UpdatedAt,
	}
}

func assetFromCreateResponse(r *api.PostAPIV1OrganizationsOrganizationIDAssetsResponse) *Asset {
	return &Asset{
		ID:                   r.ID,
		Name:                 r.Name,
		OrganizationID:       r.OrganizationID,
		Lifecycle:            AssetLifecycle(r.Lifecycle),
		Sku:                  r.Sku,
		StartURL:             r.StartURL,
		MaxRequestsPerSecond: r.MaxRequestsPerSecond,
		CreatedAt:            r.CreatedAt,
		UpdatedAt:            r.UpdatedAt,
	}
}

func assetsPageFromResponse(r *api.GetAPIV1OrganizationsOrganizationIDAssetsResponse) *Page[AssetListItem] {
	items := make([]AssetListItem, 0, len(r.Items))
	for _, item := range r.Items {
		items = append(items, AssetListItem{
			ID:        item.ID,
			Name:      item.Name,
			Lifecycle: AssetLifecycle(item.Lifecycle),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return &Page[AssetListItem]{
		Items: items,
		PageInfo: PageInfo{
			NextCursor: ptrValue(r.NextCursor),
			HasMore:    r.NextCursor != nil && *r.NextCursor != "",
		},
	}
}
