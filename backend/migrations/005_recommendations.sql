-- Recommandations générées par l'engine d'analyse
CREATE TABLE IF NOT EXISTS recommendations (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id   UUID NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at  TIMESTAMPTZ,
    priority     VARCHAR(16) NOT NULL,    -- HIGH, MEDIUM, LOW
    type         VARCHAR(64) NOT NULL,    -- rightsizing, missing_limits, add_hpa, ...
    title        TEXT NOT NULL,
    description  TEXT,
    target       TEXT,                    -- namespace/deployment
    savings_eur_annual DOUBLE PRECISION,
    confidence   DOUBLE PRECISION,
    patch_yaml   TEXT,                    -- Patch K8s prêt à appliquer
    status       VARCHAR(32) NOT NULL DEFAULT 'open'  -- open, applied, dismissed
);

CREATE INDEX IF NOT EXISTS idx_reco_cluster ON recommendations(cluster_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_reco_status ON recommendations(status) WHERE status = 'open';
