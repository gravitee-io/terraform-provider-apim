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

variable "mode" {
  type = string
}

resource "apim_mcp_proxy" "test" {
  hrid         = var.hrid
  entity_id    = var.entity_id
  name         = "Acceptance ${var.hrid}"
  context_path = "/mcp/${var.hrid}"
  mode         = var.mode
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
