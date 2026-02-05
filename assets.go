package xbow

import (
	"context"
	"iter"
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
	ApprovedTimeWindows  *ApprovedTimeWindows
	Credentials          []Credential
	DNSBoundaryRules     []DNSBoundaryRule
	Headers              map[string][]string
	HTTPBoundaryRules    []HTTPBoundaryRule
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
			ApprovedTimeWindows:  convertApprovedTimeWindowsToBody(req.ApprovedTimeWindows),
			Credentials:          convertCredentialsToBody(req.Credentials),
			DNSBoundaryRules:     convertDNSBoundaryRulesToBody(req.DNSBoundaryRules),
			Headers:              convertHeadersToBody(req.Headers),
			HTTPBoundaryRules:    convertHTTPBoundaryRulesToBody(req.HTTPBoundaryRules),
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
		StartURL:             strPtrFromNullable(r.StartURL),
		MaxRequestsPerSecond: intPtrFromNullable(r.MaxRequestsPerSecond),
		ApprovedTimeWindows:  convertApprovedTimeWindowsFromGet(r.ApprovedTimeWindows),
		Credentials:          convertCredentialsFromGet(r.Credentials),
		DNSBoundaryRules:     convertDNSBoundaryRulesFromGet(r.DNSBoundaryRules),
		Headers:              convertHeadersFromGet(r.Headers),
		HTTPBoundaryRules:    convertHTTPBoundaryRulesFromGet(r.HTTPBoundaryRules),
		Checks:               convertChecksFromGet(r.Checks),
		ArchiveAt:            timePtrFromNullable(r.ArchiveAt),
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
		StartURL:             strPtrFromNullable(r.StartURL),
		MaxRequestsPerSecond: intPtrFromNullable(r.MaxRequestsPerSecond),
		ApprovedTimeWindows:  convertApprovedTimeWindowsFromPut(r.ApprovedTimeWindows),
		Credentials:          convertCredentialsFromPut(r.Credentials),
		DNSBoundaryRules:     convertDNSBoundaryRulesFromPut(r.DNSBoundaryRules),
		Headers:              convertHeadersFromPut(r.Headers),
		HTTPBoundaryRules:    convertHTTPBoundaryRulesFromPut(r.HTTPBoundaryRules),
		Checks:               convertChecksFromPut(r.Checks),
		ArchiveAt:            timePtrFromNullable(r.ArchiveAt),
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
		StartURL:             strPtrFromNullable(r.StartURL),
		MaxRequestsPerSecond: intPtrFromNullable(r.MaxRequestsPerSecond),
		ApprovedTimeWindows:  convertApprovedTimeWindowsFromCreate(r.ApprovedTimeWindows),
		Credentials:          convertCredentialsFromCreate(r.Credentials),
		DNSBoundaryRules:     convertDNSBoundaryRulesFromCreate(r.DNSBoundaryRules),
		Headers:              convertHeadersFromCreate(r.Headers),
		HTTPBoundaryRules:    convertHTTPBoundaryRulesFromCreate(r.HTTPBoundaryRules),
		Checks:               convertChecksFromCreate(r.Checks),
		ArchiveAt:            timePtrFromNullable(r.ArchiveAt),
		CreatedAt:            r.CreatedAt,
		UpdatedAt:            r.UpdatedAt,
	}
}

// Conversion helpers for GET response types

func convertApprovedTimeWindowsFromGet(atw api.GetAPIV1AssetsAssetID_Response_ApprovedTimeWindows) *ApprovedTimeWindows {
	if atw.Tz == "" && len(atw.Entries) == 0 {
		return nil
	}
	entries := make([]TimeWindowEntry, 0, len(atw.Entries))
	for _, e := range atw.Entries {
		entries = append(entries, TimeWindowEntry{
			StartWeekday: int(e.StartWeekday),
			StartTime:    e.StartTime,
			EndWeekday:   int(e.EndWeekday),
			EndTime:      e.EndTime,
		})
	}
	return &ApprovedTimeWindows{Tz: atw.Tz, Entries: entries}
}

func convertCredentialsFromGet(creds api.GetAPIV1AssetsAssetID_Response_Credentials) []Credential {
	if creds == nil {
		return nil
	}
	result := make([]Credential, 0, len(creds))
	for _, c := range creds {
		cred := Credential{
			ID:       c.ID,
			Name:     c.Name,
			Type:     string(c.Type),
			Username: c.Username,
			Password: c.Password,
		}
		if c.EmailAddress != nil {
			email := string(*c.EmailAddress)
			cred.EmailAddress = &email
		}
		cred.AuthenticatorURI = c.AuthenticatorURI
		result = append(result, cred)
	}
	return result
}

func convertDNSBoundaryRulesFromGet(rules api.GetAPIV1AssetsAssetID_Response_DNSBoundaryRules) []DNSBoundaryRule {
	if rules == nil {
		return nil
	}
	result := make([]DNSBoundaryRule, 0, len(rules))
	for _, r := range rules {
		result = append(result, DNSBoundaryRule{
			ID:                r.ID,
			Action:            DNSBoundaryRuleAction(r.Action),
			Type:              string(r.Type),
			Filter:            r.Filter,
			IncludeSubdomains: r.IncludeSubdomains,
		})
	}
	return result
}

func convertHeadersFromGet(headers map[string]api.GetAPIV1AssetsAssetID_Response_Headers_AnyOf_AdditionalProperties) map[string][]string {
	if headers == nil {
		return nil
	}
	result := make(map[string][]string, len(headers))
	for k, v := range headers {
		if v.GetAPIV1AssetsAssetID_Response_Headers_AnyOf_AdditionalProperties_AnyOf == nil {
			continue
		}
		anyOf := v.GetAPIV1AssetsAssetID_Response_Headers_AnyOf_AdditionalProperties_AnyOf
		if anyOf.IsB() {
			result[k] = []string(anyOf.B)
		} else if anyOf.IsA() {
			result[k] = []string{anyOf.A}
		}
	}
	return result
}

func convertHTTPBoundaryRulesFromGet(rules api.GetAPIV1AssetsAssetID_Response_HTTPBoundaryRules) []HTTPBoundaryRule {
	if rules == nil {
		return nil
	}
	result := make([]HTTPBoundaryRule, 0, len(rules))
	for _, r := range rules {
		result = append(result, HTTPBoundaryRule{
			ID:                r.ID,
			Action:            HTTPBoundaryRuleAction(r.Action),
			Type:              string(r.Type),
			Filter:            r.Filter,
			IncludeSubdomains: r.IncludeSubdomains,
		})
	}
	return result
}

func convertChecksFromGet(checks api.GetAPIV1AssetsAssetID_Response_Checks) *AssetChecks {
	return &AssetChecks{
		AssetReachable: AssetCheck{
			State:   AssetCheckState(checks.AssetReachable.State),
			Message: checks.AssetReachable.Message,
			Error:   convertAssetReachableErrorFromGet(checks.AssetReachable.ErrorData),
		},
		Credentials: AssetCheck{
			State:   AssetCheckState(checks.Credentials.State),
			Message: checks.Credentials.Message,
			Error:   convertSimpleErrorFromGet(checks.Credentials.ErrorData.Type),
		},
		DNSBoundaryRules: AssetCheck{
			State:   AssetCheckState(checks.DNSBoundaryRules.State),
			Message: checks.DNSBoundaryRules.Message,
			Error:   convertSimpleErrorFromGet(checks.DNSBoundaryRules.ErrorData.Type),
		},
		UpdatedAt: timePtrFromNullable(checks.UpdatedAt),
	}
}

func convertAssetReachableErrorFromGet(errData api.GetAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error) *AssetCheckError {
	oneOf := errData.GetAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf
	if oneOf == nil {
		return nil
	}

	// Try each variant (dns, timeout, network, http, waf)
	if v, err := oneOf.AsGetAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf_0(); err == nil {
		return &AssetCheckError{Type: string(v.Type), Code: v.Code}
	}
	if v, err := oneOf.AsGetAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf_1(); err == nil {
		return &AssetCheckError{Type: string(v.Type)}
	}
	if v, err := oneOf.AsGetAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf_2(); err == nil {
		return &AssetCheckError{Type: string(v.Type), Code: v.Code}
	}
	if v, err := oneOf.AsGetAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf_3(); err == nil {
		return &AssetCheckError{Type: string(v.Type), Status: int(v.Status)}
	}
	if v, err := oneOf.AsGetAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf_4(); err == nil {
		result := &AssetCheckError{Type: string(v.Type)}
		if v.WafProvider != nil {
			result.WafProvider = string(*v.WafProvider)
		}
		return result
	}
	return nil
}

func convertSimpleErrorFromGet(errType string) *AssetCheckError {
	if errType == "" {
		return nil
	}
	return &AssetCheckError{Type: errType}
}

// Conversion helpers for PUT response types

func convertApprovedTimeWindowsFromPut(atw api.PutAPIV1AssetsAssetID_Response_ApprovedTimeWindows) *ApprovedTimeWindows {
	if atw.Tz == "" && len(atw.Entries) == 0 {
		return nil
	}
	entries := make([]TimeWindowEntry, 0, len(atw.Entries))
	for _, e := range atw.Entries {
		entries = append(entries, TimeWindowEntry{
			StartWeekday: int(e.StartWeekday),
			StartTime:    e.StartTime,
			EndWeekday:   int(e.EndWeekday),
			EndTime:      e.EndTime,
		})
	}
	return &ApprovedTimeWindows{Tz: atw.Tz, Entries: entries}
}

func convertCredentialsFromPut(creds api.PutAPIV1AssetsAssetID_Response_Credentials) []Credential {
	if creds == nil {
		return nil
	}
	result := make([]Credential, 0, len(creds))
	for _, c := range creds {
		cred := Credential{
			ID:       c.ID,
			Name:     c.Name,
			Type:     string(c.Type),
			Username: c.Username,
			Password: c.Password,
		}
		if c.EmailAddress != nil {
			email := string(*c.EmailAddress)
			cred.EmailAddress = &email
		}
		cred.AuthenticatorURI = c.AuthenticatorURI
		result = append(result, cred)
	}
	return result
}

func convertDNSBoundaryRulesFromPut(rules api.PutAPIV1AssetsAssetID_Response_DNSBoundaryRules) []DNSBoundaryRule {
	if rules == nil {
		return nil
	}
	result := make([]DNSBoundaryRule, 0, len(rules))
	for _, r := range rules {
		result = append(result, DNSBoundaryRule{
			ID:                r.ID,
			Action:            DNSBoundaryRuleAction(r.Action),
			Type:              string(r.Type),
			Filter:            r.Filter,
			IncludeSubdomains: r.IncludeSubdomains,
		})
	}
	return result
}

func convertHeadersFromPut(headers map[string]api.PutAPIV1AssetsAssetID_Response_Headers_AnyOf_AdditionalProperties) map[string][]string {
	if headers == nil {
		return nil
	}
	result := make(map[string][]string, len(headers))
	for k, v := range headers {
		if v.PutAPIV1AssetsAssetID_Response_Headers_AnyOf_AdditionalProperties_AnyOf == nil {
			continue
		}
		anyOf := v.PutAPIV1AssetsAssetID_Response_Headers_AnyOf_AdditionalProperties_AnyOf
		if anyOf.IsB() {
			result[k] = []string(anyOf.B)
		} else if anyOf.IsA() {
			result[k] = []string{anyOf.A}
		}
	}
	return result
}

func convertHTTPBoundaryRulesFromPut(rules api.PutAPIV1AssetsAssetID_Response_HTTPBoundaryRules) []HTTPBoundaryRule {
	if rules == nil {
		return nil
	}
	result := make([]HTTPBoundaryRule, 0, len(rules))
	for _, r := range rules {
		result = append(result, HTTPBoundaryRule{
			ID:                r.ID,
			Action:            HTTPBoundaryRuleAction(r.Action),
			Type:              string(r.Type),
			Filter:            r.Filter,
			IncludeSubdomains: r.IncludeSubdomains,
		})
	}
	return result
}

func convertChecksFromPut(checks api.PutAPIV1AssetsAssetID_Response_Checks) *AssetChecks {
	return &AssetChecks{
		AssetReachable: AssetCheck{
			State:   AssetCheckState(checks.AssetReachable.State),
			Message: checks.AssetReachable.Message,
			Error:   convertAssetReachableErrorFromPut(checks.AssetReachable.ErrorData),
		},
		Credentials: AssetCheck{
			State:   AssetCheckState(checks.Credentials.State),
			Message: checks.Credentials.Message,
			Error:   convertSimpleErrorFromPut(checks.Credentials.ErrorData.Type),
		},
		DNSBoundaryRules: AssetCheck{
			State:   AssetCheckState(checks.DNSBoundaryRules.State),
			Message: checks.DNSBoundaryRules.Message,
			Error:   convertSimpleErrorFromPut(checks.DNSBoundaryRules.ErrorData.Type),
		},
		UpdatedAt: timePtrFromNullable(checks.UpdatedAt),
	}
}

func convertAssetReachableErrorFromPut(errData api.PutAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error) *AssetCheckError {
	oneOf := errData.PutAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf
	if oneOf == nil {
		return nil
	}

	if v, err := oneOf.AsPutAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf_0(); err == nil {
		return &AssetCheckError{Type: string(v.Type), Code: v.Code}
	}
	if v, err := oneOf.AsPutAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf_1(); err == nil {
		return &AssetCheckError{Type: string(v.Type)}
	}
	if v, err := oneOf.AsPutAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf_2(); err == nil {
		return &AssetCheckError{Type: string(v.Type), Code: v.Code}
	}
	if v, err := oneOf.AsPutAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf_3(); err == nil {
		return &AssetCheckError{Type: string(v.Type), Status: int(v.Status)}
	}
	if v, err := oneOf.AsPutAPIV1AssetsAssetID_Response_Checks_AssetReachable_Error_AnyOf_OneOf_4(); err == nil {
		result := &AssetCheckError{Type: string(v.Type)}
		if v.WafProvider != nil {
			result.WafProvider = string(*v.WafProvider)
		}
		return result
	}
	return nil
}

func convertSimpleErrorFromPut(errType string) *AssetCheckError {
	if errType == "" {
		return nil
	}
	return &AssetCheckError{Type: errType}
}

// Conversion helpers for POST (create) response types

func convertApprovedTimeWindowsFromCreate(atw api.PostAPIV1OrganizationsOrganizationIDAssets_Response_ApprovedTimeWindows) *ApprovedTimeWindows {
	if atw.Tz == "" && len(atw.Entries) == 0 {
		return nil
	}
	entries := make([]TimeWindowEntry, 0, len(atw.Entries))
	for _, e := range atw.Entries {
		entries = append(entries, TimeWindowEntry{
			StartWeekday: int(e.StartWeekday),
			StartTime:    e.StartTime,
			EndWeekday:   int(e.EndWeekday),
			EndTime:      e.EndTime,
		})
	}
	return &ApprovedTimeWindows{Tz: atw.Tz, Entries: entries}
}

func convertCredentialsFromCreate(creds api.PostAPIV1OrganizationsOrganizationIDAssets_Response_Credentials) []Credential {
	if creds == nil {
		return nil
	}
	result := make([]Credential, 0, len(creds))
	for _, c := range creds {
		cred := Credential{
			ID:       c.ID,
			Name:     c.Name,
			Type:     string(c.Type),
			Username: c.Username,
			Password: c.Password,
		}
		if c.EmailAddress != nil {
			email := string(*c.EmailAddress)
			cred.EmailAddress = &email
		}
		cred.AuthenticatorURI = c.AuthenticatorURI
		result = append(result, cred)
	}
	return result
}

func convertDNSBoundaryRulesFromCreate(rules api.PostAPIV1OrganizationsOrganizationIDAssets_Response_DNSBoundaryRules) []DNSBoundaryRule {
	if rules == nil {
		return nil
	}
	result := make([]DNSBoundaryRule, 0, len(rules))
	for _, r := range rules {
		result = append(result, DNSBoundaryRule{
			ID:                r.ID,
			Action:            DNSBoundaryRuleAction(r.Action),
			Type:              string(r.Type),
			Filter:            r.Filter,
			IncludeSubdomains: r.IncludeSubdomains,
		})
	}
	return result
}

func convertHeadersFromCreate(headers map[string]api.PostAPIV1OrganizationsOrganizationIDAssets_Response_Headers_AnyOf_AdditionalProperties) map[string][]string {
	if headers == nil {
		return nil
	}
	result := make(map[string][]string, len(headers))
	for k, v := range headers {
		if v.PostAPIV1OrganizationsOrganizationIDAssets_Response_Headers_AnyOf_AdditionalProperties_AnyOf == nil {
			continue
		}
		anyOf := v.PostAPIV1OrganizationsOrganizationIDAssets_Response_Headers_AnyOf_AdditionalProperties_AnyOf
		if anyOf.IsB() {
			result[k] = []string(anyOf.B)
		} else if anyOf.IsA() {
			result[k] = []string{anyOf.A}
		}
	}
	return result
}

func convertHTTPBoundaryRulesFromCreate(rules api.PostAPIV1OrganizationsOrganizationIDAssets_Response_HTTPBoundaryRules) []HTTPBoundaryRule {
	if rules == nil {
		return nil
	}
	result := make([]HTTPBoundaryRule, 0, len(rules))
	for _, r := range rules {
		result = append(result, HTTPBoundaryRule{
			ID:                r.ID,
			Action:            HTTPBoundaryRuleAction(r.Action),
			Type:              string(r.Type),
			Filter:            r.Filter,
			IncludeSubdomains: r.IncludeSubdomains,
		})
	}
	return result
}

func convertChecksFromCreate(checks api.PostAPIV1OrganizationsOrganizationIDAssets_Response_Checks) *AssetChecks {
	return &AssetChecks{
		AssetReachable: AssetCheck{
			State:   AssetCheckState(checks.AssetReachable.State),
			Message: checks.AssetReachable.Message,
			Error:   convertAssetReachableErrorFromCreate(checks.AssetReachable.ErrorData),
		},
		Credentials: AssetCheck{
			State:   AssetCheckState(checks.Credentials.State),
			Message: checks.Credentials.Message,
			Error:   convertSimpleErrorFromCreate(checks.Credentials.ErrorData.Type),
		},
		DNSBoundaryRules: AssetCheck{
			State:   AssetCheckState(checks.DNSBoundaryRules.State),
			Message: checks.DNSBoundaryRules.Message,
			Error:   convertSimpleErrorFromCreate(checks.DNSBoundaryRules.ErrorData.Type),
		},
		UpdatedAt: timePtrFromNullable(checks.UpdatedAt),
	}
}

func convertAssetReachableErrorFromCreate(errData api.PostAPIV1OrganizationsOrganizationIDAssets_Response_Checks_AssetReachable_Error) *AssetCheckError {
	oneOf := errData.PostAPIV1OrganizationsOrganizationIDAssets_Response_Checks_AssetReachable_Error_AnyOf_OneOf
	if oneOf == nil {
		return nil
	}

	if v, err := oneOf.AsPostAPIV1OrganizationsOrganizationIDAssets_Response_Checks_AssetReachable_Error_AnyOf_OneOf_0(); err == nil {
		return &AssetCheckError{Type: string(v.Type), Code: v.Code}
	}
	if v, err := oneOf.AsPostAPIV1OrganizationsOrganizationIDAssets_Response_Checks_AssetReachable_Error_AnyOf_OneOf_1(); err == nil {
		return &AssetCheckError{Type: string(v.Type)}
	}
	if v, err := oneOf.AsPostAPIV1OrganizationsOrganizationIDAssets_Response_Checks_AssetReachable_Error_AnyOf_OneOf_2(); err == nil {
		return &AssetCheckError{Type: string(v.Type), Code: v.Code}
	}
	if v, err := oneOf.AsPostAPIV1OrganizationsOrganizationIDAssets_Response_Checks_AssetReachable_Error_AnyOf_OneOf_3(); err == nil {
		return &AssetCheckError{Type: string(v.Type), Status: int(v.Status)}
	}
	if v, err := oneOf.AsPostAPIV1OrganizationsOrganizationIDAssets_Response_Checks_AssetReachable_Error_AnyOf_OneOf_4(); err == nil {
		result := &AssetCheckError{Type: string(v.Type)}
		if v.WafProvider != nil {
			result.WafProvider = string(*v.WafProvider)
		}
		return result
	}
	return nil
}

func convertSimpleErrorFromCreate(errType string) *AssetCheckError {
	if errType == "" {
		return nil
	}
	return &AssetCheckError{Type: errType}
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
