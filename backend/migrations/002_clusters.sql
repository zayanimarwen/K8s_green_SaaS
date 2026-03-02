-- Clusters (dans chaque schema tenant)
-- Exemple pour tenant_macif — même structure pour tous les tenants
CREATE TABLE IF NOT EXISTS clusters (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(255) NOT NULL,
    provider     VARCHAR(64) NOT NULL,  -- aws, gcp, azure, on-prem
    region       VARCHAR(64),
    environment  VARCHAR(32) NOT NULL DEFAULT 'production',
    k8s_version  VARCHAR(32),
    agent_version VARCHAR(32),
    last_seen_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    active       BOOLEAN NOT NULL DEFAULT TRUE,
    metadata     JSONB,
    UNIQUE(name)
);

CREATE INDEX IF NOT EXISTS idx_clusters_active ON clusters(active) WHERE active = TRUE;
