#!/bin/bash
set -e
TENANT_ID=$1; TENANT_NAME=$2; PLAN=${3:-starter}
if [ -z "$TENANT_ID" ] || [ -z "$TENANT_NAME" ]; then
  echo "Usage: $0 <tenant-id> <tenant-name> [plan]"
  exit 1
fi
echo "🏢 Création tenant: $TENANT_NAME ($TENANT_ID)"
docker compose exec -T postgres psql -U green -d green_saas << SQL
INSERT INTO tenants (id, name, plan) VALUES ('$TENANT_ID', '$TENANT_NAME', '$PLAN');
SELECT create_tenant_schema('${TENANT_ID//-/_}');
SQL
echo "✅ Tenant $TENANT_NAME créé"
