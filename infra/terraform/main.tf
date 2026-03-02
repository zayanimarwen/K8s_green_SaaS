terraform {
  required_version = ">= 1.8"
  required_providers {
    ovh   = { source = "ovh/ovh",         version = "~> 0.43" }
    helm  = { source = "hashicorp/helm",  version = "~> 2.14" }
    kubernetes = { source = "hashicorp/kubernetes", version = "~> 2.31" }
  }
  backend "s3" {
    bucket = "k8s-green-tfstate"
    key    = "prod/terraform.tfstate"
    region = "eu-west-3"
  }
}

provider "ovh" {
  endpoint = var.ovh_endpoint
}

module "k8s_cluster" {
  source       = "./modules/k8s-cluster"
  cluster_name = var.cluster_name
  region       = var.region
  node_count   = var.node_count
  node_flavor  = var.node_flavor
}

module "postgres" {
  source       = "./modules/postgres"
  cluster_name = var.cluster_name
  region       = var.region
  db_name      = "green_saas"
  plan         = var.db_plan
}

module "redis" {
  source       = "./modules/redis"
  cluster_name = var.cluster_name
  region       = var.region
  plan         = var.redis_plan
}

module "nats" {
  source       = "./modules/nats"
  namespace    = "green-system"
  cluster_name = var.cluster_name
}
