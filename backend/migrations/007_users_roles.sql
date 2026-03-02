-- Mapping utilisateurs / rôles par tenant
-- Les utilisateurs sont dans Keycloak, on stocke ici les préférences et la config
CREATE TABLE IF NOT EXISTS user_settings (
    user_id      VARCHAR(128) PRIMARY KEY,
    tenant_id    VARCHAR(64) NOT NULL,
    email        VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    role         VARCHAR(32) NOT NULL DEFAULT 'viewer',
    preferences  JSONB DEFAULT '{}',
    notification_email    BOOLEAN DEFAULT TRUE,
    notification_slack    BOOLEAN DEFAULT FALSE,
    slack_webhook_url     TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_user_tenant ON user_settings(tenant_id);
