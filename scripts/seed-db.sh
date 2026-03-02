#!/bin/bash
set -e
echo "🌱 Seed données de test..."
docker compose exec -T postgres psql -U green -d green_saas << 'SQL'
-- Tenant de démo
INSERT INTO tenants (id, name, plan) VALUES
  ('tenant-macif', 'MACIF', 'enterprise'),
  ('tenant-maif',  'MAIF',  'enterprise')
ON CONFLICT (id) DO NOTHING;

-- Créer schemas tenants
SELECT create_tenant_schema('macif');
SELECT create_tenant_schema('maif');

-- Clusters de démo
INSERT INTO clusters (id, name, provider, region, environment)
VALUES
  (gen_random_uuid(), 'macif-prod-k8s1', 'on-prem', 'fr-paris', 'production'),
  (gen_random_uuid(), 'macif-staging-k8s1', 'on-prem', 'fr-paris', 'staging')
ON CONFLICT DO NOTHING;
SQL
echo "✅ Seed terminé"
