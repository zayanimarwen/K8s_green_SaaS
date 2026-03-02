variable "cluster_name" {}
variable "region"       {}
variable "db_name"      { default = "green_saas" }
variable "plan"         { default = "essential" }

resource "ovh_cloud_project_database" "postgres" {
  engine     = "postgresql"
  version    = "16"
  plan       = var.plan
  nodes      = [{ region = var.region }]
  flavor     = "db1-7"
  description = "${var.cluster_name}-postgres"
}

# Extension TimescaleDB activée via migration SQL
output "connection_url" {
  value = "postgres://${ovh_cloud_project_database.postgres.user}:PASSWORD@${ovh_cloud_project_database.postgres.endpoints[0].uri}/${var.db_name}"
  sensitive = true
}
