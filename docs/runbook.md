# Runbook Opérationnel

## Incidents courants

### Agent silencieux (alerte AgentSilent)
```bash
kubectl logs -n green-system -l app=green-agent --tail=50
kubectl describe pod -n green-system -l app=green-agent
# Vérifier la connectivité NATS
kubectl exec -n green-system -l app=green-agent -- nc -zv nats.green-system 4222
```

### API haute latence
```bash
# Vérifier TimescaleDB
kubectl exec -it -n green-system postgres-0 -- psql -U green -c "SELECT * FROM pg_stat_activity WHERE state = 'active';"
# Vérifier Redis
kubectl exec -it -n green-system redis-0 -- redis-cli info stats
```

### Worker NATS lag élevé
```bash
nats consumer info METRICS backend-worker
# Scale up workers si nécessaire
kubectl scale deployment green-worker --replicas=5 -n green-system
```

## Procédures de maintenance

### Rotation certificats mTLS
```bash
./scripts/rotate-certs.sh tenant-macif cluster-prod-k8s1
```

### Backup PostgreSQL
```bash
./scripts/backup-db.sh
```

### Ajouter un tenant
```bash
./scripts/create-tenant.sh tenant-nouvelle-entreprise "Nouvelle Entreprise" enterprise
```
