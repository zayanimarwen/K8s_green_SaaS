# 🌿 K8s Green SaaS — Enterprise Platform

Plateforme SaaS multi-tenant pour l'optimisation des clusters Kubernetes.
Détecte le surprovisionnement, calcule l'impact carbone, recommande des actions.

## Architecture

| Composant | Techno | Rôle |
|-----------|--------|------|
| `agent/`    | Go 1.22 + client-go | Collecte K8s + Prometheus, in-cluster |
| `backend/`  | Go + Gin + TimescaleDB | API REST + Worker NATS |
| `frontend/` | React 19 + Tailwind  | Dashboard SaaS |
| `infra/`    | Terraform + Helm + ArgoCD | IaC + GitOps |

## Démarrage rapide

```bash
make setup    # Bootstrap env dev (Go, Node, Docker)
make dev      # Lance la stack complète
# → Dashboard: http://localhost:3000
# → API: http://localhost:9000
# → Keycloak: http://localhost:8080 (admin/admin123)
# → NATS: http://localhost:8222
```

## Clients cibles

- MACIF, MAIF, MAAF — assurance mutuelle
- Isolation tenant garantie (schema PostgreSQL + ACL NATS + RBAC Keycloak)
- Compatible PingFederate via SAML 2.0 bridge Keycloak

## Docs

- [Architecture](docs/architecture.md)
- [Onboarding client](docs/onboarding-client.md)
- [API Reference](docs/api-reference.md)
- [Multi-tenancy](docs/multi-tenancy.md)
- [Sécurité](docs/security.md)
- [Runbook](docs/runbook.md)
