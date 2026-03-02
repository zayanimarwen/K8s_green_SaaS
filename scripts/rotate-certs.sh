#!/bin/bash
set -e
TENANT_ID=$1; CLUSTER_ID=$2
if [ -z "$TENANT_ID" ] || [ -z "$CLUSTER_ID" ]; then
  echo "Usage: $0 <tenant-id> <cluster-id>"
  exit 1
fi
echo "🔑 Rotation certificats mTLS: $TENANT_ID / $CLUSTER_ID"
# Générer nouveau cert via Vault
vault write pki/issue/green-agent \
  common_name="agent.$CLUSTER_ID.$TENANT_ID" \
  ttl=8760h \
  format=pem_bundle > /tmp/new-cert.json
# L'External Secrets Operator va synchroniser automatiquement dans 24h
# Pour forcer la synchronisation immédiate:
kubectl annotate externalsecret green-agent-mtls \
  -n green-system \
  force-sync=$(date +%s) \
  --overwrite
echo "✅ Rotation initiée — resync dans ~30s"
