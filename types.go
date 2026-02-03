package xbow

import "time"

// AssessmentState represents the current state of an assessment.
type AssessmentState string

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
