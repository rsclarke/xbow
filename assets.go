package xbow

import (
	"context"
	"encoding/json"
	"iter"
	"reflect"
	"strings"
	"time"

	"github.com/rsclarke/xbow/internal/api"
)

// Nullable conversion helpers - these map zero values to nil to preserve JSON null semantics.
// The OpenAPI spec uses anyOf[T, null] for these fields, and the generated code represents
// null as the zero value of T.

func strPtrFromNullable(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtrFromNullable(n int) *int {
	if n == 0 {
		return nil
	}
	return &n
}

func timePtrFromNullable(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

// AssetsService handles asset-related API calls.
type AssetsService struct {
	client *Client
}

// Get retrieves an asset by ID.
func (s *AssetsService) Get(ctx context.Context, id string) (*Asset, error) {
	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

	opts := &api.GetAPIV1AssetsAssetIDRequestOptions{
		PathParams: &api.GetAPIV1AssetsAssetIDPath{
			AssetID: id,
		},
		Header: &api.GetAPIV1AssetsAssetIDHeaders{
			XXBOWAPIVersion: api.GetAPIV1AssetsAssetIDHeaderXXBOWAPIVersionN20260201,
		},
	}

	resp, err := s.client.raw.GetAPIV1AssetsAssetID(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return assetFromGetResponse(resp), nil
}

// UpdateAssetRequest specifies the parameters for updating an asset.
type UpdateAssetRequest struct {
	Name                 string               `json:"name"`
	StartURL             string               `json:"startUrl"`
	MaxRequestsPerSecond int                  `json:"maxRequestsPerSecond"`
	Sku                  *string              `json:"sku,omitempty"`
	ApprovedTimeWindows  *ApprovedTimeWindows `json:"approvedTimeWindows,omitempty"`
	Credentials          []Credential         `json:"credentials"`
	DNSBoundaryRules     []DNSBoundaryRule    `json:"dnsBoundaryRules"`
	Headers              map[string][]string  `json:"headers"`
	HTTPBoundaryRules    []HTTPBoundaryRule   `json:"httpBoundaryRules"`
}

// Update updates an asset.
func (s *AssetsService) Update(ctx context.Context, id string, req *UpdateAssetRequest) (*Asset, error) {
	if req == nil {
		return nil, &Error{Code: "ERR_INVALID_REQUEST", Message: "UpdateAssetRequest cannot be nil"}
	}

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
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
			ApprovedTimeWindows:  convertApprovedTimeWindowsToBody(req.ApprovedTimeWindows),
			Credentials:          convertCredentialsToBody(req.Credentials),
			DNSBoundaryRules:     convertDNSBoundaryRulesToBody(req.DNSBoundaryRules),
			Headers:              convertHeadersToBody(req.Headers),
			HTTPBoundaryRules:    convertHTTPBoundaryRulesToBody(req.HTTPBoundaryRules),
		},
	}

	resp, err := s.client.raw.PutAPIV1AssetsAssetID(ctx, opts, auth)
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

	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
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

	resp, err := s.client.raw.PostAPIV1OrganizationsOrganizationIDAssets(ctx, opts, auth)
	if err != nil {
		return nil, wrapError(err)
	}

	return assetFromCreateResponse(resp), nil
}

// ListByOrganization returns a page of assets for an organization.
func (s *AssetsService) ListByOrganization(ctx context.Context, organizationID string, opts *ListOptions) (*Page[AssetListItem], error) {
	auth, err := s.client.orgAuthEditor()
	if err != nil {
		return nil, err
	}

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

	resp, err := s.client.raw.GetAPIV1OrganizationsOrganizationIDAssets(ctx, reqOpts, auth)
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

// assetFromJSON marshals a generated response to JSON and converts it to an Asset.
// All three generated asset response types (Get, Put, Create) serialize to the
// same JSON wire format, so a single conversion path handles all of them.
//
// The generated MarshalJSON methods can panic on nil oneOf pointers, so we use
// structToMap (reflection-based, tag-aware) to build a plain map that bypasses
// custom marshalers entirely.
func assetFromJSON(resp any) *Asset {
	m := structToMap(reflect.ValueOf(resp))
	data, err := json.Marshal(m)
	if err != nil {
		return &Asset{}
	}

	var raw rawAssetJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return &Asset{}
	}

	return raw.toAsset()
}

func structToMap(v reflect.Value) any {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			return v.Interface()
		}
		m := make(map[string]any)
		t := v.Type()
		for i := range t.NumField() {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			tag := f.Tag.Get("json")
			name, _, _ := strings.Cut(tag, ",")
			if name == "" || name == "-" {
				continue
			}
			m[name] = structToMap(v.Field(i))
		}
		return m
	case reflect.Slice:
		if v.IsNil() {
			return nil
		}
		s := make([]any, v.Len())
		for i := range v.Len() {
			s[i] = structToMap(v.Index(i))
		}
		return s
	case reflect.Map:
		if v.IsNil() {
			return nil
		}
		m := make(map[string]any, v.Len())
		for _, k := range v.MapKeys() {
			m[k.String()] = structToMap(v.MapIndex(k))
		}
		return m
	default:
		return v.Interface()
	}
}

func assetFromGetResponse(r *api.GetAPIV1AssetsAssetIDResponse) *Asset {
	return assetFromJSON(r)
}

func assetFromPutResponse(r *api.PutAPIV1AssetsAssetIDResponse) *Asset {
	return assetFromJSON(r)
}

func assetFromCreateResponse(r *api.PostAPIV1OrganizationsOrganizationIDAssetsResponse) *Asset {
	return assetFromJSON(r)
}

type rawAssetJSON struct {
	ID                   string                          `json:"id"`
	Name                 string                          `json:"name"`
	OrganizationID       string                          `json:"organizationId"`
	Lifecycle            string                          `json:"lifecycle"`
	Sku                  string                          `json:"sku"`
	StartURL             string                          `json:"startUrl"`
	MaxRequestsPerSecond int                             `json:"maxRequestsPerSecond"`
	ApprovedTimeWindows  *rawApprovedTimeWindowsJSON     `json:"approvedTimeWindows,omitempty"`
	Credentials          []rawCredentialJSON             `json:"credentials"`
	DNSBoundaryRules     []rawDNSBoundaryRuleJSON        `json:"dnsBoundaryRules"`
	Headers              map[string]json.RawMessage      `json:"headers"`
	HTTPBoundaryRules    []rawHTTPBoundaryRuleJSON       `json:"httpBoundaryRules"`
	Checks               rawChecksJSON                   `json:"checks"`
	ArchiveAt            time.Time                       `json:"archiveAt"`
	CreatedAt            time.Time                       `json:"createdAt"`
	UpdatedAt            time.Time                       `json:"updatedAt"`
}

type rawApprovedTimeWindowsJSON struct {
	Tz      string              `json:"tz"`
	Entries []rawTimeWindowJSON `json:"entries"`
}

type rawTimeWindowJSON struct {
	StartWeekday int    `json:"startWeekday"`
	StartTime    string `json:"startTime"`
	EndWeekday   int    `json:"endWeekday"`
	EndTime      string `json:"endTime"`
}

type rawCredentialJSON struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	Username         string  `json:"username"`
	Password         string  `json:"password"`
	EmailAddress     *string `json:"emailAddress,omitempty"`
	AuthenticatorURI *string `json:"authenticatorUri,omitempty"`
}

type rawDNSBoundaryRuleJSON struct {
	ID                string `json:"id"`
	Action            string `json:"action"`
	Type              string `json:"type"`
	Filter            string `json:"filter"`
	IncludeSubdomains *bool  `json:"includeSubdomains,omitempty"`
}

type rawHTTPBoundaryRuleJSON struct {
	ID                string `json:"id"`
	Action            string `json:"action"`
	Type              string `json:"type"`
	Filter            string `json:"filter"`
	IncludeSubdomains *bool  `json:"includeSubdomains,omitempty"`
}

type rawChecksJSON struct {
	AssetReachable   rawCheckJSON `json:"assetReachable"`
	Credentials      rawCheckJSON `json:"credentials"`
	DNSBoundaryRules rawCheckJSON `json:"dnsBoundaryRules"`
	UpdatedAt        time.Time    `json:"updatedAt"`
}

type rawCheckJSON struct {
	State   string          `json:"state"`
	Message string          `json:"message"`
	Error   json.RawMessage `json:"error"`
}

type rawAssetCheckErrorJSON struct {
	Type        string  `json:"type"`
	Code        string  `json:"code,omitempty"`
	Status      int     `json:"status,omitempty"`
	WafProvider *string `json:"wafProvider,omitempty"`
}

func (r *rawAssetJSON) toAsset() *Asset {
	a := &Asset{
		ID:                   r.ID,
		Name:                 r.Name,
		OrganizationID:       r.OrganizationID,
		Lifecycle:            AssetLifecycle(r.Lifecycle),
		Sku:                  r.Sku,
		StartURL:             strPtrFromNullable(r.StartURL),
		MaxRequestsPerSecond: intPtrFromNullable(r.MaxRequestsPerSecond),
		ArchiveAt:            timePtrFromNullable(r.ArchiveAt),
		CreatedAt:            r.CreatedAt,
		UpdatedAt:            r.UpdatedAt,
	}

	a.ApprovedTimeWindows = convertApprovedTimeWindows(r.ApprovedTimeWindows)
	a.Credentials = convertCredentials(r.Credentials)
	a.DNSBoundaryRules = convertDNSBoundaryRules(r.DNSBoundaryRules)
	a.Headers = convertHeaders(r.Headers)
	a.HTTPBoundaryRules = convertHTTPBoundaryRules(r.HTTPBoundaryRules)
	a.Checks = convertChecks(r.Checks)

	return a
}

func convertApprovedTimeWindows(raw *rawApprovedTimeWindowsJSON) *ApprovedTimeWindows {
	if raw == nil || (raw.Tz == "" && len(raw.Entries) == 0) {
		return nil
	}
	entries := make([]TimeWindowEntry, 0, len(raw.Entries))
	for _, e := range raw.Entries {
		entries = append(entries, TimeWindowEntry{
			StartWeekday: e.StartWeekday,
			StartTime:    e.StartTime,
			EndWeekday:   e.EndWeekday,
			EndTime:      e.EndTime,
		})
	}
	return &ApprovedTimeWindows{Tz: raw.Tz, Entries: entries}
}

func convertCredentials(raw []rawCredentialJSON) []Credential {
	if raw == nil {
		return nil
	}
	result := make([]Credential, 0, len(raw))
	for _, c := range raw {
		result = append(result, Credential{
			ID:               c.ID,
			Name:             c.Name,
			Type:             c.Type,
			Username:         c.Username,
			Password:         c.Password,
			EmailAddress:     c.EmailAddress,
			AuthenticatorURI: c.AuthenticatorURI,
		})
	}
	return result
}

func convertDNSBoundaryRules(raw []rawDNSBoundaryRuleJSON) []DNSBoundaryRule {
	if raw == nil {
		return nil
	}
	result := make([]DNSBoundaryRule, 0, len(raw))
	for _, r := range raw {
		result = append(result, DNSBoundaryRule{
			ID:                r.ID,
			Action:            DNSBoundaryRuleAction(r.Action),
			Type:              r.Type,
			Filter:            r.Filter,
			IncludeSubdomains: r.IncludeSubdomains,
		})
	}
	return result
}

func convertHeaders(raw map[string]json.RawMessage) map[string][]string {
	if raw == nil {
		return nil
	}
	result := make(map[string][]string, len(raw))
	for k, v := range raw {
		var arr []string
		if json.Unmarshal(v, &arr) == nil {
			result[k] = arr
			continue
		}
		var s string
		if json.Unmarshal(v, &s) == nil {
			result[k] = []string{s}
		}
	}
	return result
}

func convertHTTPBoundaryRules(raw []rawHTTPBoundaryRuleJSON) []HTTPBoundaryRule {
	if raw == nil {
		return nil
	}
	result := make([]HTTPBoundaryRule, 0, len(raw))
	for _, r := range raw {
		result = append(result, HTTPBoundaryRule{
			ID:                r.ID,
			Action:            HTTPBoundaryRuleAction(r.Action),
			Type:              r.Type,
			Filter:            r.Filter,
			IncludeSubdomains: r.IncludeSubdomains,
		})
	}
	return result
}

func convertChecks(raw rawChecksJSON) *AssetChecks {
	return &AssetChecks{
		AssetReachable:   convertCheck(raw.AssetReachable),
		Credentials:      convertCheck(raw.Credentials),
		DNSBoundaryRules: convertCheck(raw.DNSBoundaryRules),
		UpdatedAt:        timePtrFromNullable(raw.UpdatedAt),
	}
}

func convertCheck(raw rawCheckJSON) AssetCheck {
	return AssetCheck{
		State:   AssetCheckState(raw.State),
		Message: raw.Message,
		Error:   convertCheckError(raw.Error),
	}
}

func convertCheckError(raw json.RawMessage) *AssetCheckError {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var errJSON rawAssetCheckErrorJSON
	if err := json.Unmarshal(raw, &errJSON); err != nil {
		return nil
	}
	if errJSON.Type == "" {
		return nil
	}

	result := &AssetCheckError{
		Type:   errJSON.Type,
		Code:   errJSON.Code,
		Status: errJSON.Status,
	}
	if errJSON.WafProvider != nil {
		result.WafProvider = *errJSON.WafProvider
	}
	return result
}

// Conversion helpers for request body types (domain -> generated)

func convertApprovedTimeWindowsToBody(atw *ApprovedTimeWindows) api.PutAPIV1AssetsAssetIDBody_ApprovedTimeWindows {
	if atw == nil {
		return api.PutAPIV1AssetsAssetIDBody_ApprovedTimeWindows{}
	}
	entries := make(api.PutAPIV1AssetsAssetIDBody_ApprovedTimeWindows_AnyOf_Entries, 0, len(atw.Entries))
	for _, e := range atw.Entries {
		entries = append(entries, api.PutAPIV1AssetsAssetIDBody_ApprovedTimeWindows_AnyOf_Entries_Item{
			StartWeekday: api.PutAPIV1AssetsAssetIDBodyApprovedTimeWindowsAnyOfEntriesStartWeekday(e.StartWeekday),
			StartTime:    e.StartTime,
			EndWeekday:   api.PutAPIV1AssetsAssetIDBodyApprovedTimeWindowsAnyOfEntriesEndWeekday(e.EndWeekday),
			EndTime:      e.EndTime,
		})
	}
	return api.PutAPIV1AssetsAssetIDBody_ApprovedTimeWindows{
		Tz:      atw.Tz,
		Entries: entries,
	}
}

func convertCredentialsToBody(creds []Credential) api.PutAPIV1AssetsAssetIDBody_Credentials {
	if len(creds) == 0 {
		return nil
	}
	result := make(api.PutAPIV1AssetsAssetIDBody_Credentials, 0, len(creds))
	for _, c := range creds {
		item := api.PutAPIV1AssetsAssetIDBody_Credentials_AnyOf_Item{
			ID:               c.ID,
			Name:             c.Name,
			Type:             api.PutAPIV1AssetsAssetIDBodyCredentialsAnyOfType(c.Type),
			Username:         c.Username,
			Password:         c.Password,
			AuthenticatorURI: c.AuthenticatorURI,
		}
		result = append(result, item)
	}
	return result
}

func convertDNSBoundaryRulesToBody(rules []DNSBoundaryRule) api.PutAPIV1AssetsAssetIDBody_DNSBoundaryRules {
	if len(rules) == 0 {
		return nil
	}
	result := make(api.PutAPIV1AssetsAssetIDBody_DNSBoundaryRules, 0, len(rules))
	for _, r := range rules {
		result = append(result, api.PutAPIV1AssetsAssetIDBody_DNSBoundaryRules_AnyOf_Item{
			ID:                r.ID,
			Action:            api.PutAPIV1AssetsAssetIDBodyDNSBoundaryRulesAnyOfAction(r.Action),
			Type:              api.PutAPIV1AssetsAssetIDBodyDNSBoundaryRulesAnyOfType(r.Type),
			Filter:            r.Filter,
			IncludeSubdomains: r.IncludeSubdomains,
		})
	}
	return result
}

func convertHeadersToBody(headers map[string][]string) map[string]api.PutAPIV1AssetsAssetIDBody_Headers_AnyOf_AdditionalProperties {
	if len(headers) == 0 {
		return nil
	}
	result := make(map[string]api.PutAPIV1AssetsAssetIDBody_Headers_AnyOf_AdditionalProperties, len(headers))
	for k, v := range headers {
		var anyOf api.PutAPIV1AssetsAssetIDBody_Headers_AnyOf_AdditionalProperties_AnyOf
		if len(v) == 1 {
			anyOf.A = v[0]
			anyOf.N = 1
		} else {
			anyOf.B = api.PutAPIV1AssetsAssetIDBody_Headers_AnyOf_AdditionalProperties_AnyOf_1(v)
			anyOf.N = 2
		}
		result[k] = api.PutAPIV1AssetsAssetIDBody_Headers_AnyOf_AdditionalProperties{
			PutAPIV1AssetsAssetIDBody_Headers_AnyOf_AdditionalProperties_AnyOf: &anyOf,
		}
	}
	return result
}

func convertHTTPBoundaryRulesToBody(rules []HTTPBoundaryRule) api.PutAPIV1AssetsAssetIDBody_HTTPBoundaryRules {
	if len(rules) == 0 {
		return nil
	}
	result := make(api.PutAPIV1AssetsAssetIDBody_HTTPBoundaryRules, 0, len(rules))
	for _, r := range rules {
		result = append(result, api.PutAPIV1AssetsAssetIDBody_HTTPBoundaryRules_AnyOf_Item{
			ID:                r.ID,
			Action:            api.PutAPIV1AssetsAssetIDBodyHTTPBoundaryRulesAnyOfAction(r.Action),
			Type:              api.PutAPIV1AssetsAssetIDBodyHTTPBoundaryRulesAnyOfType(r.Type),
			Filter:            r.Filter,
			IncludeSubdomains: r.IncludeSubdomains,
		})
	}
	return result
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
			NextCursor: r.NextCursor,
			HasMore:    r.NextCursor != nil && *r.NextCursor != "",
		},
	}
}
