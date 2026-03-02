package repository

import (
	"context"

	"github.com/k8s-green/backend/internal/domain"
)

type TenantRepo struct{ db *Postgres }

func (r *TenantRepo) List(ctx context.Context) ([]domain.Tenant, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT id, name, plan, created_at, active FROM tenants WHERE active = true ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tenants []domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		rows.Scan(&t.ID, &t.Name, &t.Plan, &t.CreatedAt, &t.Active)
		tenants = append(tenants, t)
	}
	return tenants, nil
}

func (r *TenantRepo) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, name, plan, created_at, active FROM tenants WHERE id = $1`, id,
	).Scan(&t.ID, &t.Name, &t.Plan, &t.CreatedAt, &t.Active)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TenantRepo) Create(ctx context.Context, t *domain.Tenant) error {
	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO tenants (id, name, plan, active) VALUES ($1, $2, $3, true)
		 ON CONFLICT (id) DO NOTHING`,
		t.ID, t.Name, t.Plan,
	)
	if err != nil {
		return err
	}
	// Créer le schema isolé du tenant
	_, err = r.db.Pool.Exec(ctx, `SELECT create_tenant_schema($1)`, t.ID)
	return err
}
