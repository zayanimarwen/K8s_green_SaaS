package domain

import "time"

type Cluster struct {
	ID           string     `json:"id" db:"id"`
	TenantID     string     `json:"tenant_id"`
	Name         string     `json:"name" db:"name"`
	Provider     string     `json:"provider" db:"provider"`
	Region       string     `json:"region" db:"region"`
	Environment  string     `json:"environment" db:"environment"`
	K8sVersion   string     `json:"k8s_version,omitempty" db:"k8s_version"`
	AgentVersion string     `json:"agent_version,omitempty" db:"agent_version"`
	LastSeenAt   *time.Time `json:"last_seen_at,omitempty" db:"last_seen_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	Active       bool       `json:"active" db:"active"`
}
