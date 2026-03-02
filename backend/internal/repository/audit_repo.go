package repository

import (
	"context"
	"time"
)

type AuditRepo struct{ db *Postgres }

type AuditEntry struct {
	TenantID   string
	UserID     string
	UserEmail  string
	Action     string
	Resource   string
	ResourceID string
	IPAddress  string
	UserAgent  string
	StatusCode int
	DurationMs int64
	Metadata   map[string]interface{}
	CreatedAt  time.Time
}

func (r *AuditRepo) Insert(ctx context.Context, e AuditEntry) error {
	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO audit_logs
		 (tenant_id, user_id, user_email, action, resource, resource_id,
		  ip_address, user_agent, status_code, duration_ms, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		e.TenantID, e.UserID, e.UserEmail, e.Action, e.Resource, e.ResourceID,
		e.IPAddress, e.UserAgent, e.StatusCode, e.DurationMs,
		time.Now(),
	)
	return err
}

func (r *AuditRepo) ListByTenant(ctx context.Context, tenantID string, limit int) ([]AuditEntry, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT tenant_id, user_id, user_email, action, resource, ip_address, status_code, duration_ms, created_at
		 FROM audit_logs WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2`,
		tenantID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []AuditEntry
	for rows.Next() {
		var e AuditEntry
		rows.Scan(&e.TenantID, &e.UserID, &e.UserEmail, &e.Action, &e.Resource,
			&e.IPAddress, &e.StatusCode, &e.DurationMs, &e.CreatedAt)
		entries = append(entries, e)
	}
	return entries, nil
}
