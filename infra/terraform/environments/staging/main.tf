module "saas" {
  source       = "../../"
  environment  = "staging"
  cluster_name = "green-staging"
  node_count   = 3
  node_flavor  = "b3-8"
  region       = "GRA7"
}
