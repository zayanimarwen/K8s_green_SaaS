variable "namespace"    { default = "green-system" }

resource "helm_release" "keycloak" {
  name       = "keycloak"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "keycloak"
  version    = "21.0.0"
  namespace  = var.namespace

  set { name = "auth.adminUser";     value = "admin" }
  set { name = "replicaCount";       value = "2" }
  set { name = "postgresql.enabled"; value = "false" }
  set { name = "externalDatabase.host";     value = var.postgres_host }
  set { name = "extraEnvVars[0].name";  value = "KC_HOSTNAME" }
  set { name = "extraEnvVars[0].value"; value = var.hostname }
}

variable "postgres_host" { default = "localhost" }
variable "hostname"       { default = "keycloak.example.com" }
