-- Schéma multi-tenant : chaque tenant a son propre schema PostgreSQL
CREATE TABLE IF NOT EXISTS tenants (
    id          VARCHAR(64) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    plan        VARCHAR(32) NOT NULL DEFAULT 'starter',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    metadata    JSONB
);

-- Création dynamique du schema isolé par tenant (appelé au onboarding)
CREATE OR REPLACE FUNCTION create_tenant_schema(tenant_id TEXT) RETURNS void AS $$
BEGIN
    EXECUTE format('CREATE SCHEMA IF NOT EXISTS tenant_%s', tenant_id);
    EXECUTE format('GRANT USAGE ON SCHEMA tenant_%s TO green_api', tenant_id);
    EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA tenant_%s GRANT ALL ON TABLES TO green_api', tenant_id);
END;
$$ LANGUAGE plpgsql;
