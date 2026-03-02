variable "ovh_endpoint"  { default = "ovh-eu" }
variable "cluster_name"  { default = "green-saas" }
variable "region"        { default = "GRA7" }
variable "node_count"    { default = 3 }
variable "node_flavor"   { default = "b3-8" }
variable "db_plan"       { default = "essential" }
variable "redis_plan"    { default = "essential-2" }
variable "environment"   {
  default = "prod"
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "environment doit être dev, staging ou prod."
  }
}
