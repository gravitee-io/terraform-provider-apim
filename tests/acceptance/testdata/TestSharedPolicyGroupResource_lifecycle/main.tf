provider "apim" {
  organization_id = "DEFAULT"
}

variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

resource "apim_shared_policy_group" "test" {
  api_type        = "PROXY"
  environment_id  = var.environment_id
  hrid            = var.hrid
  name            = "terraform_example"
  organization_id = var.organization_id
  phase           = "REQUEST"
}
