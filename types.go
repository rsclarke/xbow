package xbow

import "time"

// AssetLifecycle represents the lifecycle state of an asset.
type AssetLifecycle string

// Possible values for AssetLifecycle.
const (
	AssetLifecycleActive   AssetLifecycle = "active"
	AssetLifecycleArchived AssetLifecycle = "archived"
)

// Asset represents a web application to be assessed.
type Asset struct {
	ID                   string               `json:"id"`
	Name                 string               `json:"name"`
	OrganizationID       string               `json:"organizationId"`
	Lifecycle            AssetLifecycle       `json:"lifecycle"`
	Sku                  string               `json:"sku"`
	StartURL             *string              `json:"startUrl"`
	MaxRequestsPerSecond *int                 `json:"maxRequestsPerSecond"`
	ApprovedTimeWindows  *ApprovedTimeWindows `json:"approvedTimeWindows"`
	Credentials          []Credential         `json:"credentials"`
	DNSBoundaryRules     []DNSBoundaryRule    `json:"dnsBoundaryRules"`
	Headers              map[string][]string  `json:"headers"`
	HTTPBoundaryRules    []HTTPBoundaryRule   `json:"httpBoundaryRules"`
	Checks               *AssetChecks         `json:"checks"`
	ArchiveAt            *time.Time           `json:"archiveAt"`
	CreatedAt            time.Time            `json:"createdAt"`
	UpdatedAt            time.Time            `json:"updatedAt"`
}

// ApprovedTimeWindows represents time windows when assessments can run.
type ApprovedTimeWindows struct {
	Tz      string            `json:"tz"`
	Entries []TimeWindowEntry `json:"entries"`
}

// TimeWindowEntry represents a single time window entry.
type TimeWindowEntry struct {
	StartWeekday int    `json:"startWeekday"`
	StartTime    string `json:"startTime"`
	EndWeekday   int    `json:"endWeekday"`
	EndTime      string `json:"endTime"`
}

// Credential represents authentication credentials for an asset.
type Credential struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	Username         string  `json:"username"`
	Password         string  `json:"password"`
	EmailAddress     *string `json:"emailAddress,omitempty"`
	AuthenticatorURI *string `json:"authenticatorUri,omitempty"`
}

// DNSBoundaryRuleAction represents the action for a DNS boundary rule.
// DNS rules do not support allow-auth (only HTTP rules do).
type DNSBoundaryRuleAction string

// Possible values for DNSBoundaryRuleAction.
const (
	DNSBoundaryRuleActionAllowAttack DNSBoundaryRuleAction = "allow-attack"
	DNSBoundaryRuleActionAllowVisit  DNSBoundaryRuleAction = "allow-visit"
	DNSBoundaryRuleActionDeny        DNSBoundaryRuleAction = "deny"
)

// HTTPBoundaryRuleAction represents the action for an HTTP boundary rule.
type HTTPBoundaryRuleAction string

// Possible values for HTTPBoundaryRuleAction.
const (
	HTTPBoundaryRuleActionAllowAttack HTTPBoundaryRuleAction = "allow-attack"
	HTTPBoundaryRuleActionAllowAuth   HTTPBoundaryRuleAction = "allow-auth"
	HTTPBoundaryRuleActionAllowVisit  HTTPBoundaryRuleAction = "allow-visit"
	HTTPBoundaryRuleActionDeny        HTTPBoundaryRuleAction = "deny"
)

// DNSBoundaryRule represents a DNS boundary rule for an asset.
type DNSBoundaryRule struct {
	ID                string                `json:"id"`
	Action            DNSBoundaryRuleAction `json:"action"`
	Type              string                `json:"type"`
	Filter            string                `json:"filter"`
	IncludeSubdomains *bool                 `json:"includeSubdomains,omitempty"`
}

// HTTPBoundaryRule represents an HTTP boundary rule for an asset.
type HTTPBoundaryRule struct {
	ID                string                 `json:"id"`
	Action            HTTPBoundaryRuleAction `json:"action"`
	Type              string                 `json:"type"`
	Filter            string                 `json:"filter"`
	IncludeSubdomains *bool                  `json:"includeSubdomains,omitempty"`
}

// AssetChecks represents validation checks for an asset.
type AssetChecks struct {
	AssetReachable   AssetCheck `json:"assetReachable"`
	Credentials      AssetCheck `json:"credentials"`
	DNSBoundaryRules AssetCheck `json:"dnsBoundaryRules"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}

// AssetCheckState represents the state of an asset check.
type AssetCheckState string

// Possible values for AssetCheckState.
const (
	AssetCheckStateUnchecked AssetCheckState = "unchecked"
	AssetCheckStateChecking  AssetCheckState = "checking"
	AssetCheckStateValid     AssetCheckState = "valid"
	AssetCheckStateInvalid   AssetCheckState = "invalid"
)

// AssetCheck represents a single validation check.
type AssetCheck struct {
	State   AssetCheckState  `json:"state"`
	Message string           `json:"message"`
	Error   *AssetCheckError `json:"error,omitempty"`
}

// AssetCheckError represents error details for a failed asset check.
type AssetCheckError struct {
	Type        string `json:"type"`
	Code        string `json:"code,omitempty"`
	Status      int    `json:"status,omitempty"`
	WafProvider string `json:"wafProvider,omitempty"`
}

// AssetListItem represents an asset in list responses (fewer fields).
type AssetListItem struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Lifecycle AssetLifecycle `json:"lifecycle"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// AssessmentState represents the current state of an assessment.
type AssessmentState string

// Possible values for AssessmentState.
const (
	AssessmentStateWaitingForCapacity   AssessmentState = "waiting-for-capacity"
	AssessmentStateRunning              AssessmentState = "running"
	AssessmentStateSucceeded            AssessmentState = "succeeded"
	AssessmentStateReportReady          AssessmentState = "report-ready"
	AssessmentStateFailed               AssessmentState = "failed"
	AssessmentStateCancelling           AssessmentState = "cancelling"
	AssessmentStateCancelled            AssessmentState = "cancelled"
	AssessmentStatePaused               AssessmentState = "paused"
	AssessmentStateWaitingForTimeWindow AssessmentState = "waiting-for-time-window"
)

// Assessment represents a security assessment.
type Assessment struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	AssetID        string            `json:"assetId"`
	OrganizationID string            `json:"organizationId"`
	State          AssessmentState   `json:"state"`
	Progress       float64           `json:"progress"`
	AttackCredits  int64             `json:"attackCredits"`
	RecentEvents   []AssessmentEvent `json:"recentEvents"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

// AssessmentListItem represents an assessment in list responses (fewer fields).
type AssessmentListItem struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	State     AssessmentState `json:"state"`
	Progress  float64         `json:"progress"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

// AssessmentEvent represents an event in an assessment's history.
type AssessmentEvent struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason,omitempty"`
}
