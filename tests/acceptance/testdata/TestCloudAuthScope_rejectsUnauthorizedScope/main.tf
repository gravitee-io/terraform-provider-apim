provider "apim" {
  cloud_auth      = var.cloud_token
  server_url      = "http://127.0.0.1:9/automation"
  organization_id = var.provider_organization_id
  environment_id  = var.provider_environment_id
}

variable "cloud_token" {
  type      = string
  sensitive = true
}

variable "provider_organization_id" {
  type    = string
  default = "DEFAULT"
}

variable "provider_environment_id" {
  type    = string
  default = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "resource_organization_id" {
  type    = string
  default = null
}

variable "resource_environment_id" {
  type    = string
  default = null
}

resource "apim_group" "test" {
  hrid            = var.hrid
  name            = "Test"
  organization_id = var.resource_organization_id
  environment_id  = var.resource_environment_id
}
