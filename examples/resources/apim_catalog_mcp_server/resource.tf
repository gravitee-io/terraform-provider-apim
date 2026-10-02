variable "github_mcp_authorization" {
  type        = string
  sensitive   = true
  description = "Full value of the Authorization header, such as \"Bearer <token>\""
}

resource "apim_catalog_mcp_server" "github" {
  hrid        = "github"
  entity_id   = "mcp-server.github"
  description = "GitHub tools for coding agents"
  connection = {
    endpoint = "https://api.githubcopilot.com/mcp/"
    auth = {
      header = {
        name  = "Authorization"
        value = var.github_mcp_authorization
      }
    }
  }
}

# The platform discovers the server when it is applied. Its tools come back in
# the state, each with the entity id an authorization policy names.
output "github_tools" {
  value = apim_catalog_mcp_server.github.tools
}
