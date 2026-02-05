package xbow

import (
	"context"
	"iter"

	"github.com/rsclarke/xbow/internal/api"
)

// WebhooksService handles webhook-related API calls.
type WebhooksService struct {
	client *Client
}

// CreateWebhookRequest contains the parameters for creating a webhook subscription.
type CreateWebhookRequest struct {
	APIVersion WebhookAPIVersion  `json:"apiVersion"`
	TargetURL  string             `json:"targetUrl"`
	Events     []WebhookEventType `json:"events"`
}

// UpdateWebhookRequest contains the parameters for updating a webhook subscription.
// All fields are optional; only provided fields will be updated.
type UpdateWebhookRequest struct {
	APIVersion *WebhookAPIVersion `json:"apiVersion,omitempty"`
	TargetURL  *string            `json:"targetUrl,omitempty"`
	Events     []WebhookEventType `json:"events,omitempty"`
}

// Get retrieves a webhook subscription by ID.
func (s *WebhooksService) Get(ctx context.Context, id string) (*Webhook, error) {
	if id == "" {
		return nil, &Error{Code: "ERR_INVALID_PARAM", Message: "webhook id is required"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	opts := &api.GetAPIV1WebhooksWebhookIDRequestOptions{
		PathParams: &api.GetAPIV1WebhooksWebhookIDPath{
			WebhookID: id,
		},
		Header: &api.GetAPIV1WebhooksWebhookIDHeaders{
			XXBOWAPIVersion: api.GetAPIV1WebhooksWebhookIDHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.GetAPIV1WebhooksWebhookID(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return webhookFromGetResponse(resp), nil
}

// Update updates an existing webhook subscription.
func (s *WebhooksService) Update(ctx context.Context, id string, req *UpdateWebhookRequest) (*Webhook, error) {
	if id == "" {
		return nil, &Error{Code: "ERR_INVALID_PARAM", Message: "webhook id is required"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	body := &api.PatchAPIV1WebhooksWebhookIDBody{}
	if req != nil {
		if req.APIVersion != nil {
			v := api.PatchAPIV1WebhooksWebhookIDBodyAPIVersion(*req.APIVersion)
			body.APIVersion = &v
		}
		if req.TargetURL != nil {
			body.TargetURL = req.TargetURL
		}
		if len(req.Events) > 0 {
			events := make(api.PatchAPIV1WebhooksWebhookIDBody_Events, 0, len(req.Events))
			for _, e := range req.Events {
				item := api.PatchAPIV1WebhooksWebhookIDBody_Events_Item{}
				item.PatchAPIV1WebhooksWebhookIDBody_Events_AnyOf = &api.PatchAPIV1WebhooksWebhookIDBody_Events_AnyOf{}
				_ = item.PatchAPIV1WebhooksWebhookIDBody_Events_AnyOf.FromString(string(e))
				events = append(events, item)
			}
			body.Events = &events
		}
	}

	opts := &api.PatchAPIV1WebhooksWebhookIDRequestOptions{
		PathParams: &api.PatchAPIV1WebhooksWebhookIDPath{
			WebhookID: id,
		},
		Body: body,
		Header: &api.PatchAPIV1WebhooksWebhookIDHeaders{
			XXBOWAPIVersion: api.PatchAPIV1WebhooksWebhookIDHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.PatchAPIV1WebhooksWebhookID(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return webhookFromPatchResponse(resp), nil
}

// Delete deletes a webhook subscription.
func (s *WebhooksService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return &Error{Code: "ERR_INVALID_PARAM", Message: "webhook id is required"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return err
	}

	opts := &api.DeleteAPIV1WebhooksWebhookIDRequestOptions{
		PathParams: &api.DeleteAPIV1WebhooksWebhookIDPath{
			WebhookID: id,
		},
		Header: &api.DeleteAPIV1WebhooksWebhookIDHeaders{
			XXBOWAPIVersion: api.DeleteAPIV1WebhooksWebhookIDHeaderXXBOWAPIVersionN20260201,
		},
	}

	_, err = s.client.raw.DeleteAPIV1WebhooksWebhookID(ctx, opts, auth)
	if err != nil {
		return wrapError(err)
	}

	return nil
}

// Ping sends a ping event to a webhook subscription to test connectivity.
func (s *WebhooksService) Ping(ctx context.Context, id string) error {
	if id == "" {
		return &Error{Code: "ERR_INVALID_PARAM", Message: "webhook id is required"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return err
	}

	opts := &api.PostAPIV1WebhooksWebhookIDPingRequestOptions{
		PathParams: &api.PostAPIV1WebhooksWebhookIDPingPath{
			WebhookID: id,
		},
		Header: &api.PostAPIV1WebhooksWebhookIDPingHeaders{
			XXBOWAPIVersion: api.PostAPIV1WebhooksWebhookIDPingHeaderXXBOWAPIVersionN20260201,
		},
	}

	_, err = s.client.raw.PostAPIV1WebhooksWebhookIDPing(ctx, opts, auth)
	if err != nil {
		return wrapError(err)
	}

	return nil
}

// ListByOrganization returns a page of webhook subscriptions for an organization.
func (s *WebhooksService) ListByOrganization(ctx context.Context, organizationID string, opts *ListOptions) (*Page[WebhookListItem], error) {
	if organizationID == "" {
		return nil, &Error{Code: "ERR_INVALID_PARAM", Message: "organization id is required"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	reqOpts := &api.GetAPIV1OrganizationsOrganizationIDWebhooksRequestOptions{
		PathParams: &api.GetAPIV1OrganizationsOrganizationIDWebhooksPath{
			OrganizationID: organizationID,
		},
		Header: &api.GetAPIV1OrganizationsOrganizationIDWebhooksHeaders{
			XXBOWAPIVersion: api.GetAPIV1OrganizationsOrganizationIDWebhooksHeaderXXBOWAPIVersionN20260201,
		},
	}

	if opts != nil {
		reqOpts.Query = &api.GetAPIV1OrganizationsOrganizationIDWebhooksQuery{}
		if opts.Limit > 0 {
			reqOpts.Query.Limit = &opts.Limit
		}
		if opts.After != "" {
			reqOpts.Query.After = &opts.After
		}
	}

	resp, err := s.client.raw.GetAPIV1OrganizationsOrganizationIDWebhooks(ctx, reqOpts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return webhooksPageFromResponse(resp), nil
}

// AllByOrganization returns an iterator over all webhook subscriptions for an organization.
// Use this for automatic pagination:
//
//	for webhook, err := range client.Webhooks.AllByOrganization(ctx, orgID, nil) {
//	    if err != nil {
//	        return err
//	    }
//	    fmt.Println(webhook.TargetURL)
//	}
func (s *WebhooksService) AllByOrganization(ctx context.Context, organizationID string, opts *ListOptions) iter.Seq2[WebhookListItem, error] {
	return paginate(ctx, opts, func(ctx context.Context, pageOpts *ListOptions) (*Page[WebhookListItem], error) {
		return s.ListByOrganization(ctx, organizationID, pageOpts)
	})
}

// Create creates a new webhook subscription for an organization.
func (s *WebhooksService) Create(ctx context.Context, organizationID string, req *CreateWebhookRequest) (*Webhook, error) {
	if organizationID == "" {
		return nil, &Error{Code: "ERR_INVALID_PARAM", Message: "organization id is required"}
	}
	if req == nil {
		return nil, &Error{Code: "ERR_INVALID_PARAM", Message: "request is required"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	events := make(api.PostAPIV1OrganizationsOrganizationIDWebhooksBody_Events, 0, len(req.Events))
	for _, e := range req.Events {
		item := api.PostAPIV1OrganizationsOrganizationIDWebhooksBody_Events_Item{}
		item.PostAPIV1OrganizationsOrganizationIDWebhooksBody_Events_AnyOf = &api.PostAPIV1OrganizationsOrganizationIDWebhooksBody_Events_AnyOf{}
		_ = item.PostAPIV1OrganizationsOrganizationIDWebhooksBody_Events_AnyOf.FromString(string(e))
		events = append(events, item)
	}

	opts := &api.PostAPIV1OrganizationsOrganizationIDWebhooksRequestOptions{
		PathParams: &api.PostAPIV1OrganizationsOrganizationIDWebhooksPath{
			OrganizationID: organizationID,
		},
		Body: &api.PostAPIV1OrganizationsOrganizationIDWebhooksBody{
			APIVersion: api.PostAPIV1OrganizationsOrganizationIDWebhooksBodyAPIVersion(req.APIVersion),
			TargetURL:  req.TargetURL,
			Events:     events,
		},
		Header: &api.PostAPIV1OrganizationsOrganizationIDWebhooksHeaders{
			XXBOWAPIVersion: api.PostAPIV1OrganizationsOrganizationIDWebhooksHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.PostAPIV1OrganizationsOrganizationIDWebhooks(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return webhookFromCreateResponse(resp), nil
}

// ListDeliveries returns a page of delivery history for a webhook subscription.
func (s *WebhooksService) ListDeliveries(ctx context.Context, webhookID string, opts *ListOptions) (*Page[WebhookDelivery], error) {
	if webhookID == "" {
		return nil, &Error{Code: "ERR_INVALID_PARAM", Message: "webhook id is required"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	reqOpts := &api.GetAPIV1WebhooksWebhookIDDeliveriesRequestOptions{
		PathParams: &api.GetAPIV1WebhooksWebhookIDDeliveriesPath{
			WebhookID: webhookID,
		},
		Header: &api.GetAPIV1WebhooksWebhookIDDeliveriesHeaders{
			XXBOWAPIVersion: api.GetAPIV1WebhooksWebhookIDDeliveriesHeaderXXBOWAPIVersionN20260201,
		},
	}

	if opts != nil {
		reqOpts.Query = &api.GetAPIV1WebhooksWebhookIDDeliveriesQuery{}
		if opts.Limit > 0 {
			reqOpts.Query.Limit = &opts.Limit
		}
		if opts.After != "" {
			reqOpts.Query.After = &opts.After
		}
	}

	resp, err := s.client.raw.GetAPIV1WebhooksWebhookIDDeliveries(ctx, reqOpts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return deliveriesPageFromResponse(resp), nil
}

// AllDeliveries returns an iterator over all deliveries for a webhook subscription.
// Use this for automatic pagination:
//
//	for delivery, err := range client.Webhooks.AllDeliveries(ctx, webhookID, nil) {
//	    if err != nil {
//	        return err
//	    }
//	    fmt.Printf("Delivery at %s: success=%v\n", delivery.SentAt, delivery.Success)
//	}
func (s *WebhooksService) AllDeliveries(ctx context.Context, webhookID string, opts *ListOptions) iter.Seq2[WebhookDelivery, error] {
	return paginate(ctx, opts, func(ctx context.Context, pageOpts *ListOptions) (*Page[WebhookDelivery], error) {
		return s.ListDeliveries(ctx, webhookID, pageOpts)
	})
}

// Conversion functions from generated types to domain types

func webhookFromGetResponse(r *api.GetAPIV1WebhooksWebhookIDResponse) *Webhook {
	return &Webhook{
		ID:         r.ID,
		APIVersion: WebhookAPIVersion(r.APIVersion),
		TargetURL:  r.TargetURL,
		Events:     convertEventsFromGet(r.Events),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func webhookFromPatchResponse(r *api.PatchAPIV1WebhooksWebhookIDResponse) *Webhook {
	return &Webhook{
		ID:         r.ID,
		APIVersion: WebhookAPIVersion(r.APIVersion),
		TargetURL:  r.TargetURL,
		Events:     convertEventsFromPatch(r.Events),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func webhookFromCreateResponse(r *api.PostAPIV1OrganizationsOrganizationIDWebhooksResponse) *Webhook {
	return &Webhook{
		ID:         r.ID,
		APIVersion: WebhookAPIVersion(r.APIVersion),
		TargetURL:  r.TargetURL,
		Events:     convertEventsFromCreate(r.Events),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func webhooksPageFromResponse(r *api.GetAPIV1OrganizationsOrganizationIDWebhooksResponse) *Page[WebhookListItem] {
	items := make([]WebhookListItem, 0, len(r.Items))
	for _, item := range r.Items {
		items = append(items, WebhookListItem{
			ID:         item.ID,
			APIVersion: WebhookAPIVersion(item.APIVersion),
			TargetURL:  item.TargetURL,
			Events:     convertEventsFromList(item.Events),
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		})
	}

	return &Page[WebhookListItem]{
		Items: items,
		PageInfo: PageInfo{
			NextCursor: r.NextCursor,
			HasMore:    r.NextCursor != nil && *r.NextCursor != "",
		},
	}
}

func deliveriesPageFromResponse(r *api.GetAPIV1WebhooksWebhookIDDeliveriesResponse) *Page[WebhookDelivery] {
	items := make([]WebhookDelivery, 0, len(r.Items))
	for _, item := range r.Items {
		items = append(items, WebhookDelivery{
			Payload: item.Payload,
			Request: WebhookDeliveryRequest{
				Body:    item.Request.Body,
				Headers: item.Request.Headers,
			},
			Response: WebhookDeliveryResponse{
				Body:    item.Response.Body,
				Headers: item.Response.Headers,
				Status:  item.Response.Status,
			},
			SentAt:  item.SentAt,
			Success: item.Success,
		})
	}

	return &Page[WebhookDelivery]{
		Items: items,
		PageInfo: PageInfo{
			NextCursor: r.NextCursor,
			HasMore:    r.NextCursor != nil && *r.NextCursor != "",
		},
	}
}

// Event conversion helpers for different response types

func convertEventsFromGet(events api.GetAPIV1WebhooksWebhookID_Response_Events) []WebhookEventType {
	result := make([]WebhookEventType, 0, len(events))
	for _, e := range events {
		if e.GetAPIV1WebhooksWebhookID_Response_Events_AnyOf == nil {
			continue
		}
		if s, err := e.GetAPIV1WebhooksWebhookID_Response_Events_AnyOf.AsString(); err == nil {
			result = append(result, WebhookEventType(s))
		}
	}
	return result
}

func convertEventsFromPatch(events api.PatchAPIV1WebhooksWebhookID_Response_Events) []WebhookEventType {
	result := make([]WebhookEventType, 0, len(events))
	for _, e := range events {
		if e.PatchAPIV1WebhooksWebhookID_Response_Events_AnyOf == nil {
			continue
		}
		if s, err := e.PatchAPIV1WebhooksWebhookID_Response_Events_AnyOf.AsString(); err == nil {
			result = append(result, WebhookEventType(s))
		}
	}
	return result
}

func convertEventsFromCreate(events api.PostAPIV1OrganizationsOrganizationIDWebhooks_Response_Events) []WebhookEventType {
	result := make([]WebhookEventType, 0, len(events))
	for _, e := range events {
		if e.PostAPIV1OrganizationsOrganizationIDWebhooks_Response_Events_AnyOf == nil {
			continue
		}
		if s, err := e.PostAPIV1OrganizationsOrganizationIDWebhooks_Response_Events_AnyOf.AsString(); err == nil {
			result = append(result, WebhookEventType(s))
		}
	}
	return result
}

func convertEventsFromList(events api.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items_Events) []WebhookEventType {
	result := make([]WebhookEventType, 0, len(events))
	for _, e := range events {
		if e.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items_Events_AnyOf == nil {
			continue
		}
		if s, err := e.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items_Events_AnyOf.AsString(); err == nil {
			result = append(result, WebhookEventType(s))
		}
	}
	return result
}
