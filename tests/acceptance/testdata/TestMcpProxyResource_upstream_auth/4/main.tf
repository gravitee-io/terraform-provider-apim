provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

resource "apim_mcp_proxy" "test" {
  hrid         = var.hrid
  entity_id    = "mcp-proxy.${var.hrid}"
  name         = "Acceptance ${var.hrid}"
  context_path = "/mcp/${var.hrid}"
  mode         = "PROXY"
  proxy = {
    server_url = "https://mcp.example.com/mcp"
  }
  plans = [
    {
      name = "open"
      security = {
        key_less = {}
      }
    }
  ]
}
