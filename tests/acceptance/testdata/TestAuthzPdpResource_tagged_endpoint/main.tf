provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "endpoint_path" {
  type    = string
  default = null
}

resource "apim_authz_pdp" "test" {
  hrid          = var.hrid
  name          = "Acceptance ${var.hrid}"
  target_pdp_id = var.hrid
  tag           = "eu"
  endpoint      = var.endpoint_path == null ? null : { path = var.endpoint_path }
}
