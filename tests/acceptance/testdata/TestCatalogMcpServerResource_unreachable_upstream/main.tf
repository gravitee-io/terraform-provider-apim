provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "lookup" {
  type        = bool
  description = "Read the server instead of declaring it"
}

resource "apim_catalog_mcp_server" "test" {
  count = var.lookup ? 0 : 1

  hrid      = var.hrid
  entity_id = "mcp-server.${var.hrid}"
  connection = {
    endpoint = "http://unreachable.invalid/mcp"
    auth = {
      none = {}
    }
  }
}

data "apim_catalog_mcp_server" "test" {
  count = var.lookup ? 1 : 0

  hrid = var.hrid
}
