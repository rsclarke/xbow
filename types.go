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
	ID                   string         `json:"id"`
	Name                 string         `json:"name"`
	OrganizationID       string         `json:"organizationId"`
	Lifecycle            AssetLifecycle `json:"lifecycle"`
	Sku                  string         `json:"sku"`
	StartURL             string         `json:"startUrl"`
	MaxRequestsPerSecond int            `json:"maxRequestsPerSecond"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
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
	AttackCredits  int               `json:"attackCredits"`
	RecentEvents   []AssessmentEvent `json:"recentEvents"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

// AssessmentEvent represents an event in an assessment's history.
type AssessmentEvent struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason,omitempty"`
}
