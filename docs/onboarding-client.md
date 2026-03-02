# Guide d'installation agent — Client MACIF

## Prérequis
- Cluster Kubernetes >= 1.26
- Prometheus Operator installé (ou Prometheus accessible)
- Helm >= 3.14
- Accès sortant vers `nats.green-optimizer.io:4222`

## Étape 1 — Créer le namespace

```bash
kubectl create namespace green-system
```

## Étape 2 — Récupérer les credentials (fournis par K8s Green)
- `CLUSTER_ID` : identifiant unique de votre cluster
- `TENANT_ID` : votre ID tenant (`tenant-macif`)
- `NATS_URL` : URL du broker NATS SaaS
- `SIGNING_KEY` : clé HMAC-SHA256 (min 32 chars)

## Étape 3 — Installer l'agent via Helm

```bash
helm repo add k8s-green https://charts.k8s-green.io
helm repo update

helm install green-agent k8s-green/k8s-green-agent \
  --namespace green-system \
  --set config.clusterId=macif-prod-k8s1 \
  --set config.tenantId=tenant-macif \
  --set-string "externalSecrets.enabled=false" \
  --set-string "config.prometheusUrl=http://prometheus-operated.monitoring:9090"
```

## Étape 4 — Vérification

```bash
kubectl logs -n green-system -l app=green-agent -f
# Expected: "Collecte OK" toutes les 30s
```

## Accès dashboard
https://green.macif.fr — SSO via PingFederate MACIF
