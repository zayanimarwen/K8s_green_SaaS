package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// Postgres encapsule le pool de connexions PostgreSQL/TimescaleDB
type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(url string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 20
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Ping
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	log.Info().Msg("PostgreSQL connecté")

	return &Postgres{Pool: pool}, nil
}

func (p *Postgres) Close() {
	p.Pool.Close()
}

// Repos — accès centralisé
func (p *Postgres) ClusterRepo()     *ClusterRepo     { return &ClusterRepo{p} }
func (p *Postgres) TenantRepo()      *TenantRepo      { return &TenantRepo{p} }
func (p *Postgres) TimescaleRepo()   *TimescaleRepo   { return &TimescaleRepo{p} }
func (p *Postgres) AuditRepo()       *AuditRepo       { return &AuditRepo{p} }
