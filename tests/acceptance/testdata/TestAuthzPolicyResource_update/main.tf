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

variable "policy_text" {
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

locals {
  targets = [apim_authz_pdp.a.routing_scope, apim_authz_pdp.b.routing_scope]
}

resource "apim_authz_policy" "test" {
  hrid        = var.hrid
  entity_id   = "policy.${var.hrid}"
  name        = var.name
  policy_type = "CUSTOM"
  policy_text = var.policy_text
  targets     = var.reversed ? reverse(local.targets) : local.targets
}
