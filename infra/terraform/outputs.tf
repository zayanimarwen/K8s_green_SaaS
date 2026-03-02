output "kubeconfig"      { value = module.k8s_cluster.kubeconfig;  sensitive = true }
output "postgres_url"    { value = module.postgres.connection_url;  sensitive = true }
output "redis_url"       { value = module.redis.connection_url;     sensitive = true }
output "nats_url"        { value = module.nats.nats_url }
output "cluster_id"      { value = module.k8s_cluster.cluster_id }
