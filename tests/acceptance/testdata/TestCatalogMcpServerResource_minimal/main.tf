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

resource "apim_catalog_mcp_server" "test" {
  hrid      = var.hrid
  entity_id = "mcp-server.${var.hrid}"
  server_connection = {
    endpoint = var.endpoint
    auth = {
      none = {}
    }
  }
}
