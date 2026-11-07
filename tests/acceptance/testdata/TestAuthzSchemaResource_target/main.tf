provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "target" {
  type = number
}

resource "apim_authz_pdp" "first" {
  hrid          = "${var.hrid}-first"
  name          = "Acceptance ${var.hrid} first"
  target_pdp_id = "${var.hrid}-first"
}

resource "apim_authz_pdp" "second" {
  hrid          = "${var.hrid}-second"
  name          = "Acceptance ${var.hrid} second"
  target_pdp_id = "${var.hrid}-second"
}

resource "apim_authz_schema" "test" {
  hrid   = var.hrid
  name   = "Acceptance ${var.hrid}"
  target = [apim_authz_pdp.first.routing_scope, apim_authz_pdp.second.routing_scope][var.target]
  schema = <<-EOT
    entity User in [Group];
    entity Group;
  EOT
}
