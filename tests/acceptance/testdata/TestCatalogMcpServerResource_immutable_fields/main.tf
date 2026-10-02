provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "endpoint" {
  type = string
}

variable "entity_id" {
  type = string
}

resource "apim_catalog_mcp_server" "test" {
  hrid      = var.hrid
  entity_id = var.entity_id
  connection = {
    endpoint = var.endpoint
    auth = {
      none = {}
    }
  }
}
