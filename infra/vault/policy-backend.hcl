# Politique Vault — Backend API + Worker
# Principe least privilege : accès uniquement aux secrets nécessaires

path "k8s-green/backend/database" {
  capabilities = ["read"]
}

path "k8s-green/backend/redis" {
  capabilities = ["read"]
}

path "k8s-green/backend/nats" {
  capabilities = ["read"]
}

path "k8s-green/backend/jwt" {
  capabilities = ["read"]
}

# Rotation automatique des credentials DB
path "database/creds/green-backend" {
  capabilities = ["read"]
}

# Interdire tout le reste explicitement
path "k8s-green/agent/*" {
  capabilities = ["deny"]
}
