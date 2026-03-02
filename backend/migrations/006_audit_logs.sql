-- Logs d'audit RGPD — append-only, jamais de DELETE/UPDATE
CREATE TABLE IF NOT EXISTS audit_logs (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    tenant_id   VARCHAR(64) NOT NULL,
    user_id     VARCHAR(128),
    user_email  VARCHAR(255),
    action      VARCHAR(128) NOT NULL,
    resource    VARCHAR(128),
    resource_id VARCHAR(128),
    ip_address  INET,
    user_agent  TEXT,
    status_code INTEGER,
    duration_ms INTEGER,
    metadata    JSONB
);

-- Append-only : interdire UPDATE et DELETE
CREATE RULE audit_no_update AS ON UPDATE TO audit_logs DO INSTEAD NOTHING;
CREATE RULE audit_no_delete AS ON DELETE TO audit_logs DO INSTEAD NOTHING;

CREATE INDEX IF NOT EXISTS idx_audit_tenant ON audit_logs(tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_user   ON audit_logs(user_id, created_at DESC);
-- Rétention 3 ans (RGPD)
SELECT add_retention_policy('audit_logs', INTERVAL '3 years');
