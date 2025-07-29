variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

resource "apim_application" "test" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = var.hrid
  name            = "terraform example"
  description     = "Conformance tests application"
}
