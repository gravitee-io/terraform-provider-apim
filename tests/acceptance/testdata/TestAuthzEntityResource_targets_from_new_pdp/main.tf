provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

resource "apim_authz_pdp" "target" {
  hrid          = "${var.hrid}-pdp"
  name          = "Acceptance ${var.hrid} pdp"
  target_pdp_id = "${var.hrid}-pdp"
}

resource "apim_authz_entity" "test" {
  hrid        = var.hrid
  entity_id   = var.hrid
  kind        = "PRINCIPAL"
  entity_type = "Group"
  targets     = [apim_authz_pdp.target.routing_scope]
}
