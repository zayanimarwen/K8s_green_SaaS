# Modèle de sécurité

## Threat Model

| Menace | Mitigation |
|--------|------------|
| Agent compromis chez client | mTLS cert unique par cluster, signature HMAC payload |
| Exfiltration données inter-tenant | Schema PostgreSQL séparé, middleware tenant 100% routes |
| Token JWT volé | Durée de vie 5min, refresh silencieux OIDC |
| Injection NATS (usurpation tenant) | ACL mTLS + sujet `metrics.{tenant_id}.*` vérifié |
| Accès secret en clair | Vault + External Secrets Operator, zero secrets dans Git |
| Brute force API | Rate limit Redis sliding window 1000 req/min/tenant |

## Politique secrets
- **Zéro secret dans Git** (détection pre-commit git-secrets)
- Vault comme source unique de vérité
- External Secrets Operator sync Vault → K8s Secrets
- Rotation automatique certificats mTLS (24h via ExternalSecret)
- Rotation credentials DB (Vault dynamic secrets)

## Audit
- Chaque requête API loguée dans `audit_logs` (append-only)
- Rétention 3 ans (conformité RGPD + SI assurance)
- Logs Loki avec label `tenant_id` pour isolation
