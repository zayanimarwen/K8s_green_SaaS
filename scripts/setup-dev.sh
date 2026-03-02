#!/bin/bash
set -e

echo "🌿 Bootstrap environnement de développement K8s Green SaaS"

# Vérifier les dépendances
command -v docker >/dev/null || (echo "Docker requis" && exit 1)
command -v go >/dev/null || (echo "Go 1.22+ requis" && exit 1)
command -v node >/dev/null || (echo "Node.js 20+ requis" && exit 1)

echo "📦 Installation des dépendances..."
cd agent && go mod download && cd ..
cd backend && go mod download && cd ..
cd frontend && npm install && cd ..

echo "🔑 Création .env de développement..."
if [ ! -f .env ]; then
cat > .env << 'ENV'
DATABASE_URL=postgres://green:green_dev_password@localhost:5432/green_saas?sslmode=disable
REDIS_URL=redis://localhost:6379
NATS_URL=nats://localhost:4222
KEYCLOAK_URL=http://localhost:8080
KEYCLOAK_REALM=green-saas
KEYCLOAK_CLIENT_ID=green-backend
KEYCLOAK_CLIENT_SECRET=dev-secret-backend
JWT_ISSUER=http://localhost:8080/realms/green-saas
GREEN_SIGNING_KEY=dev-signing-key-minimum-32-chars-long
ENV
fi

echo "🐳 Démarrage stack Docker..."
docker compose up -d postgres redis nats

echo "⏳ Attente PostgreSQL..."
until docker compose exec postgres pg_isready -U green -q; do sleep 1; done

echo "🗄️ Application migrations..."
cd backend && go run cmd/migrate/main.go up 2>/dev/null || echo "Migration via SQL direct"
for f in migrations/*.sql; do
  docker compose exec -T postgres psql -U green -d green_saas < "../$f" 2>/dev/null || true
done
cd ..

echo "✅ Environnement prêt!"
echo "   make dev    → lance la stack complète"
echo "   make test   → lance les tests"
