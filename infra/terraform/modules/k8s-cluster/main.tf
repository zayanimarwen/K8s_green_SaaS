variable "cluster_name" {}
variable "region"       {}
variable "node_count"   { default = 3 }
variable "node_flavor"  { default = "b3-8" }

resource "ovh_cloud_project_kube" "cluster" {
  name    = var.cluster_name
  region  = var.region
  version = "1.30"
}

resource "ovh_cloud_project_kube_nodepool" "workers" {
  kube_id      = ovh_cloud_project_kube.cluster.id
  name         = "workers"
  flavor_name  = var.node_flavor
  desired_nodes = var.node_count
  min_nodes    = 2
  max_nodes    = var.node_count + 5
  autoscale    = true
}

output "cluster_id"  { value = ovh_cloud_project_kube.cluster.id }
output "kubeconfig"  { value = ovh_cloud_project_kube.cluster.kubeconfig; sensitive = true }
