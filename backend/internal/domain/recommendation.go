package domain

import "time"

type RecommendationStatus string
const (
	StatusPending  RecommendationStatus = "pending"
	StatusApplied  RecommendationStatus = "applied"
	StatusSkipped  RecommendationStatus = "skipped"
	StatusRejected RecommendationStatus = "rejected"
)

type Recommendation struct {
	ID               string               `json:"id" db:"id"`
	ClusterID        string               `json:"cluster_id" db:"cluster_id"`
	CreatedAt        time.Time            `json:"created_at" db:"created_at"`
	Priority         string               `json:"priority" db:"priority"`
	Type             string               `json:"type" db:"type"`
	Title            string               `json:"title" db:"title"`
	Description      string               `json:"description" db:"description"`
	TargetNS         string               `json:"target_ns" db:"target_ns"`
	TargetResource   string               `json:"target_resource" db:"target_resource"`
	SavingsEurAnnual float64              `json:"savings_eur_annual" db:"savings_eur_annual"`
	Confidence       float64              `json:"confidence" db:"confidence"`
	Status           RecommendationStatus `json:"status" db:"status"`
	AppliedAt        *time.Time           `json:"applied_at,omitempty" db:"applied_at"`
}
