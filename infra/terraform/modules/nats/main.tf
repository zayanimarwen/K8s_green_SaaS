variable "namespace"    { default = "green-system" }
variable "cluster_name" {}

# NATS JetStream déployé via Helm dans le cluster K8s
resource "helm_release" "nats" {
  name       = "nats"
  repository = "https://nats-io.github.io/k8s/helm/charts/"
  chart      = "nats"
  version    = "1.2.0"
  namespace  = var.namespace
  create_namespace = true

  set { name = "nats.jetstream.enabled";            value = "true" }
  set { name = "nats.jetstream.memStorage.enabled"; value = "true" }
  set { name = "nats.jetstream.fileStorage.enabled"; value = "true" }
  set { name = "cluster.enabled";                   value = "true" }
  set { name = "cluster.replicas";                  value = "3" }
}

output "nats_url" { value = "nats://nats.${var.namespace}:4222" }
