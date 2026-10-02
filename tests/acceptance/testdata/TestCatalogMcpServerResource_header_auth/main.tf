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

variable "value" {
  type      = string
  sensitive = true
}

resource "apim_catalog_mcp_server" "test" {
  hrid      = var.hrid
  entity_id = "mcp-server.${var.hrid}"
  connection = {
    endpoint = var.endpoint
    auth = {
      header = {
        name  = "Authorization"
        value = var.value
      }
    }
  }
}
