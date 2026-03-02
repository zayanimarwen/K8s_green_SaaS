module "saas" {
  source       = "../../"
  environment  = "prod"
  cluster_name = "green-prod"
  node_count   = 5
  node_flavor  = "b3-16"
  db_plan      = "business"
  region       = "GRA7"
}
