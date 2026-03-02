# Modèle d'isolation multi-tenant

## Principe
Chaque client (MACIF, MAIF, MAAF) est un **tenant isolé**.
L'isolation est garantie à **chaque couche** de la stack.

## Couches d'isolation

| Couche | Mécanisme | Garantie |
|--------|-----------|----------|
| PostgreSQL | Schema séparé par tenant | Requêtes impossible entre tenants côté DB |
| Redis | Keyspace préfixé `tenant:{id}:` | Clés jamais partagées |
| NATS | Subject `metrics.{tenant_id}.*` + ACL mTLS | Agent ne peut publier que pour son tenant |
| API | Middleware tenant — 100% des routes | Extraction JWT claim, injection context Go |
| Logs (Loki) | Label `tenant_id` chaque ligne | Queries filtrées automatiquement |
| Keycloak | Groupes par tenant, RBAC granulaire | Utilisateurs ne voient que leur tenant |

## Onboarding nouveau tenant

```bash
./scripts/create-tenant.sh tenant-maif "MAIF" enterprise
```

Ce script:
1. Insère le tenant dans la table `tenants`
2. Crée le schema `tenant_maif` dans PostgreSQL
3. Configure les ACL NATS pour le préfixe `metrics.tenant-maif.*`
4. Crée le groupe Keycloak `tenant-maif`
5. Configure les politiques Vault pour le tenant
