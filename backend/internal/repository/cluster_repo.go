package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/k8s-green/backend/internal/domain"
)

type ClusterRepo struct{ db *Postgres }

func (r *ClusterRepo) ListByTenant(ctx context.Context, tenantID string) ([]domain.Cluster, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT id, name, provider, region, environment, k8s_version, agent_version, last_seen_at, created_at, active
		 FROM clusters WHERE active = true ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clusters []domain.Cluster
	for rows.Next() {
		var c domain.Cluster
		if err := rows.Scan(&c.ID, &c.Name, &c.Provider, &c.Region, &c.Environment,
			&c.K8sVersion, &c.AgentVersion, &c.LastSeenAt, &c.CreatedAt, &c.Active); err != nil {
			continue
		}
		c.TenantID = tenantID
		clusters = append(clusters, c)
	}
	return clusters, nil
}

func (r *ClusterRepo) Create(ctx context.Context, c *domain.Cluster) error {
	c.ID = uuid.New().String()
	c.CreatedAt = time.Now()
	c.Active = true
	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO clusters (id, name, provider, region, environment, created_at, active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		c.ID, c.Name, c.Provider, c.Region, c.Environment, c.CreatedAt, c.Active,
	)
	return err
}

func (r *ClusterRepo) UpdateLastSeen(ctx context.Context, clusterID string) error {
	_, err := r.db.Pool.Exec(ctx,
		`UPDATE clusters SET last_seen_at = NOW() WHERE id = $1`,
		clusterID,
	)
	return err
}
