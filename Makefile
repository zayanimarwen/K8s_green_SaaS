.PHONY: setup dev dev-all build test lint clean migrate seed help

help: ## Afficher l'aide
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Bootstrap environnement dev complet
	@bash scripts/setup-dev.sh

dev: ## Lancer API + Worker + Frontend (sans Keycloak)
	docker compose up api worker frontend postgres redis nats

dev-all: ## Lancer toute la stack y compris Keycloak
	docker compose --profile auth up

dev-backend: ## Lancer uniquement le backend
	docker compose up api worker postgres redis nats

dev-frontend: ## Lancer uniquement le frontend
	docker compose up frontend

migrate: ## Appliquer les migrations SQL
	docker compose exec api go run ./cmd/migrate/main.go

seed: ## Seed données de test
	@bash scripts/seed-db.sh

build: ## Builder tous les binaires Go
	@echo "Build agent..."
	cd agent && CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/agent ./cmd/agent
	@echo "Build backend API..."
	cd backend && CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/api ./cmd/api
	@echo "Build backend Worker..."
	cd backend && CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/worker ./cmd/worker
	@echo "Build frontend..."
	cd frontend && npm run build

build-docker: ## Builder les images Docker
	docker build -t ghcr.io/k8s-green/agent:latest ./agent
	docker build -t ghcr.io/k8s-green/backend:latest ./backend
	docker build -t ghcr.io/k8s-green/frontend:latest ./frontend

test: ## Lancer tous les tests
	cd agent && go test ./... -v -race -timeout 60s
	cd backend && go test ./... -v -race -timeout 60s

lint: ## Linter Go et TypeScript
	cd agent && golangci-lint run
	cd backend && golangci-lint run
	cd frontend && npm run lint

clean: ## Nettoyer les artefacts
	rm -rf agent/bin backend/bin backend/tmp
	docker compose down -v

create-tenant: ## Créer un tenant (usage: make create-tenant ID=tenant-x NAME="Nom" PLAN=enterprise)
	@bash scripts/create-tenant.sh $(ID) "$(NAME)" $(PLAN)

rotate-certs: ## Rotation certificats mTLS (usage: make rotate-certs TENANT=tenant-x CLUSTER=prod-k8s1)
	@bash scripts/rotate-certs.sh $(TENANT) $(CLUSTER)

helm-lint: ## Vérifier les charts Helm
	helm lint infra/helm/k8s-green-agent
	helm lint infra/helm/k8s-green-saas

tf-init: ## Initialiser Terraform (env dev)
	cd infra/terraform/environments/dev && terraform init

tf-plan: ## Plan Terraform (env dev)
	cd infra/terraform/environments/dev && terraform plan

tf-apply: ## Apply Terraform (env dev)
	cd infra/terraform/environments/dev && terraform apply

logs: ## Voir les logs de la stack dev
	docker compose logs -f api worker

ps: ## État de la stack
	docker compose ps
