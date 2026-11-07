provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "attributes" {
  type = string
}

variable "with_parent" {
  type = bool
}

resource "apim_authz_entity" "parent" {
  hrid        = "${var.hrid}-parent"
  entity_id   = "${var.hrid}-parent"
  kind        = "PRINCIPAL"
  entity_type = "Group"
}

resource "apim_authz_entity" "test" {
  hrid        = var.hrid
  entity_id   = var.hrid
  kind        = "PRINCIPAL"
  entity_type = "User"
  attributes  = var.attributes
  parents = var.with_parent ? [
    {
      entity_type = apim_authz_entity.parent.entity_type
      entity_id   = apim_authz_entity.parent.entity_id
    }
  ] : []
}
