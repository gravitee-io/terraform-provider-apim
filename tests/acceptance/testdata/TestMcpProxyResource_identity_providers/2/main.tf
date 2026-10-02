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
  identity_providers = [
    {
      oauth2_generic = {
        name                          = "keycloak"
        issuer_url                    = "https://keycloak.example.com/realms/gravitee/"
        introspection_endpoint        = "/protocol/openid-connect/token/introspect"
        introspection_endpoint_method = "POST"
        client_id                     = "gravitee"
        client_secret                 = "acceptance-test"
      }
    }
  ]
  plans = [
    {
      name = "tokens"
      security = {
        oauth2 = {
          provider = "keycloak"
        }
      }
    }
  ]
}
