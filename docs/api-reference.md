# API Reference — K8s Green SaaS v1

Base URL: `https://api.green-optimizer.io/v1`

## Authentification
Bearer token JWT (Keycloak). Header: `Authorization: Bearer <token>`
Tenant: `X-Tenant-ID: tenant-macif`

## Endpoints

### GET /clusters
Liste les clusters du tenant courant.

### GET /clusters/:id/score
Retourne le Green Score courant.

```json
{
  "score": 78.5, "grade": "B+",
  "breakdown": {
    "cpu_efficiency": 72.0, "mem_efficiency": 68.0,
    "node_packing": 85.0, "hpa_coverage": 60.0, "limit_compliance": 90.0
  }
}
```

### GET /clusters/:id/waste
Rapport de surprovisionnement détaillé.

### GET /clusters/:id/carbon
Impact carbone estimé (CO₂, kWh, équivalences).

### GET /clusters/:id/savings
Économies potentielles annuelles.

### GET /clusters/:id/recommendations
Liste des recommandations priorisées (HIGH/MEDIUM/LOW).

### POST /clusters/:id/simulate
What-if analysis. Body: `{deployment, namespace, new_cpu_req_m, new_mem_req_mi}`

### GET /ws/clusters/:id/live
WebSocket pour le score live (push toutes les 30s).
