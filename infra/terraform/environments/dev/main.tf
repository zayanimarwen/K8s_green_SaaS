module "saas" {
  source       = "../../"
  environment  = "dev"
  cluster_name = "green-dev"
  node_count   = 2
  node_flavor  = "b3-4"
  db_plan      = "essential"
  region       = "GRA7"
}
