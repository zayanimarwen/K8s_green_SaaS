# Politique Vault — Agent (in-cluster client)
# Accès UNIQUEMENT aux certificats mTLS de l'agent

path "k8s-green/agent/mtls" {
  capabilities = ["read"]
}

path "k8s-green/agent/config" {
  capabilities = ["read"]
}

# Tout le reste interdit
path "*" {
  capabilities = ["deny"]
}
