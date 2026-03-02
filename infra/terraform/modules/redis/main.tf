variable "cluster_name" {}
variable "region"       {}
variable "plan"         { default = "essential-2" }

resource "ovh_cloud_project_database" "redis" {
  engine      = "redis"
  version     = "7.2"
  plan        = var.plan
  description = "${var.cluster_name}-redis"
  nodes       = [{ region = var.region }]
  flavor      = "db1-4"
}

output "connection_url" {
  value     = "redis://:PASSWORD@${ovh_cloud_project_database.redis.endpoints[0].uri}"
  sensitive = true
}
