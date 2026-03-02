package domain

import "time"

type User struct {
	ID          string     `json:"id" db:"id"`
	KeycloakID  string     `json:"keycloak_id" db:"keycloak_id"`
	TenantID    string     `json:"tenant_id" db:"tenant_id"`
	Email       string     `json:"email" db:"email"`
	DisplayName string     `json:"display_name" db:"display_name"`
	Role        string     `json:"role" db:"role"`
	Preferences map[string]interface{} `json:"preferences"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	LastLogin   *time.Time `json:"last_login,omitempty" db:"last_login"`
}
