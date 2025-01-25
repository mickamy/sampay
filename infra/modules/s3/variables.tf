variable "env" {
  description = "Environment name (e.g., stg, prod)"
  type        = string
}

variable "geo_locations" {
  description = "Geo locations to allow access"
  type = list(string)
}
