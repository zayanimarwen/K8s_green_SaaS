package domain

import "time"

type Tenant struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Plan      string    `json:"plan" db:"plan"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Active    bool      `json:"active" db:"active"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
