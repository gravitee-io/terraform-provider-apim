provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "target_pdp_id" {
  type = string
}

variable "tag" {
  type = string
}

resource "apim_authz_pdp" "test" {
  hrid          = var.hrid
  name          = "Acceptance ${var.hrid}"
  target_pdp_id = var.target_pdp_id
  tag           = var.tag
}
