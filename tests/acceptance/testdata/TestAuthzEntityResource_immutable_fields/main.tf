provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "entity_id" {
  type = string
}

variable "kind" {
  type = string
}

variable "entity_type" {
  type = string
}

resource "apim_authz_entity" "test" {
  hrid        = var.hrid
  entity_id   = var.entity_id
  kind        = var.kind
  entity_type = var.entity_type
}
