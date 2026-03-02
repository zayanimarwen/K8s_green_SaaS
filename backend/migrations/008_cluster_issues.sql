-- Migration 008 : Table pour les diagnostics K8s (k8sGPT-like)

CREATE TABLE IF NOT EXISTS cluster_issues (
    id              TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    cluster_id      TEXT NOT NULL,
    issue_type      TEXT NOT NULL,
    severity        TEXT NOT NULL CHECK (severity IN ('critical', 'warning', 'info')),
    namespace       TEXT NOT NULL DEFAULT '',
    resource_kind   TEXT NOT NULL DEFAULT 'Pod',
    resource_name   TEXT NOT NULL,
    message         TEXT NOT NULL,
    details         TEXT NOT NULL DEFAULT '',
    ai_explanation  TEXT,
    detected_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at     TIMESTAMPTZ,
    resolved        BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_cluster_issues_cluster   ON cluster_issues (cluster_id, detected_at DESC);
CREATE INDEX IF NOT EXISTS idx_cluster_issues_severity  ON cluster_issues (severity, resolved);
CREATE INDEX IF NOT EXISTS idx_cluster_issues_namespace ON cluster_issues (namespace, cluster_id);

COMMENT ON TABLE cluster_issues IS 'Problemes detectes dans les clusters - alimente par l agent, enrichi par LLM';
