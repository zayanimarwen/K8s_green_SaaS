# Architecture K8s Green SaaS

## Vue d'ensemble

```
[Cluster Client]                    [SaaS Platform]
  Agent Go                           Backend API (Go + Gin)
    ↓ collect K8s API                Backend Worker (NATS consumer)
    ↓ collect Prometheus             TimescaleDB (métriques)
    ↓ sign HMAC-SHA256               PostgreSQL (tenants/clusters)
    ↓                                Redis (cache + rate limit)
    └──── NATS JetStream ───────────→ Worker → Analyze → Score → Push WS
           (mTLS cert client)                         ↓
                                     Dashboard React 19
```

## Décisions architecturales

### Pourquoi deux binaires backend ?
`cmd/api` et `cmd/worker` sont des binaires séparés pour permettre un scaling indépendant.
50 clusters = beaucoup plus de messages NATS que d'appels API HTTP.

### Pourquoi TimescaleDB ?
- Compression automatique (>80% space saving sur métriques)
- Continuous aggregates natifs (daily_scores sans cron)
- Compatible PostgreSQL (pas de nouveau paradigme à apprendre)
- Rétention automatique avec `add_retention_policy`

### Isolation multi-tenant
- PostgreSQL: schema par tenant (`tenant_macif.clusters`)
- Redis: keyspace préfixé `tenant:{id}:...`
- NATS: subject `metrics.{tenant_id}.{cluster_id}` + ACL mTLS
- API: middleware tenant 100% routes, claim JWT extrait côté backend
