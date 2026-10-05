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
  hrid      = "server-${var.hrid}"
  entity_id = "mcp-server.${var.hrid}"
  server_connection = {
    endpoint = var.endpoint
    auth = {
      none = {}
    }
  }
}

locals {
  # Two of the tools the platform discovered, in the order it reports a studio's tools.
  tools = slice(sort(apim_catalog_mcp_server.test.tools[*].name), 0, 2)
}

variable "enable_fga" {
  type = bool
}

resource "apim_mcp_proxy" "test" {
  hrid         = var.hrid
  entity_id    = "mcp-proxy.${var.hrid}"
  name         = "Acceptance ${var.hrid}"
  context_path = "/mcp/${var.hrid}"
  mode         = "STUDIO"
  studio = {
    tools = [
      for tool in local.tools : {
        server = apim_catalog_mcp_server.test.hrid
        tool   = tool
      }
    ]
    upstream_auth = [
      {
        server = apim_catalog_mcp_server.test.hrid
        auth = {
          none = {}
        }
      }
    ]
    enable_fga = var.enable_fga
  }
  plans = [
    {
      name = "keys"
      security = {
        api_key = {
          source = "HEADER"
        }
      }
    }
  ]
}
