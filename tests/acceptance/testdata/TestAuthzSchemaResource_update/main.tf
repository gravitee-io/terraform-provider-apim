provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "name" {
  type = string
}

variable "schema" {
  type = string
}

resource "apim_authz_schema" "test" {
  hrid   = var.hrid
  name   = var.name
  schema = var.schema
}
