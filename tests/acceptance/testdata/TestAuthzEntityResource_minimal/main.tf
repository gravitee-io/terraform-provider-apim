provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

resource "apim_authz_entity" "test" {
  hrid        = var.hrid
  entity_id   = var.hrid
  kind        = "PRINCIPAL"
  entity_type = "Group"
}
