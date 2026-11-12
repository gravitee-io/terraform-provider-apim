provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "reversed" {
  type = bool
}

resource "apim_authz_pdp" "a" {
  hrid          = "${var.hrid}-a"
  name          = "Acceptance ${var.hrid} a"
  target_pdp_id = "${var.hrid}-a"
}

resource "apim_authz_pdp" "b" {
  hrid          = "${var.hrid}-b"
  name          = "Acceptance ${var.hrid} b"
  target_pdp_id = "${var.hrid}-b"
}

resource "apim_authz_entity" "parent_1" {
  hrid        = "${var.hrid}-p1"
  entity_id   = "${var.hrid}-p1"
  kind        = "PRINCIPAL"
  entity_type = "Group"
}

resource "apim_authz_entity" "parent_2" {
  hrid        = "${var.hrid}-p2"
  entity_id   = "${var.hrid}-p2"
  kind        = "PRINCIPAL"
  entity_type = "Group"
}

locals {
  targets = [apim_authz_pdp.a.routing_scope, apim_authz_pdp.b.routing_scope]
  parents = [
    { entity_type = "Group", entity_id = apim_authz_entity.parent_1.entity_id },
    { entity_type = "Group", entity_id = apim_authz_entity.parent_2.entity_id },
  ]
  # Indexes, not `var.reversed ? reverse(local.parents) : local.parents`:
  # Terraform 1.9 validates a conditional on a variable whose branches have
  # the same length as objects whose keys are null, and the provider refuses
  # a parent without entity_type or entity_id. Terraform 1.5 does not.
  order = var.reversed ? [1, 0] : [0, 1]
}

resource "apim_authz_entity" "test" {
  hrid        = var.hrid
  entity_id   = var.hrid
  kind        = "PRINCIPAL"
  entity_type = "User"
  targets     = [for i in local.order : local.targets[i]]
  parents     = [for i in local.order : local.parents[i]]
}
