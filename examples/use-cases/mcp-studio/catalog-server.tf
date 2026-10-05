# An upstream MCP server registered in the AI Catalog. The platform discovers it
# when the resource is applied: its tools, prompts and resources come back in the
# state, each with the entity id an authorization policy names.
resource "apim_catalog_mcp_server" "deepwiki" {
  hrid        = "example-deepwiki"
  entity_id   = "mcp-server.example-deepwiki"
  description = "DeepWiki documentation tools for coding agents"
  server_connection = {
    endpoint = "https://mcp.deepwiki.com/mcp"
    # DeepWiki is public. A server that needs a credential declares
    # header = { name, value } or oauth2 = { client_id, client_secret, token_url }.
    auth = {
      none = {}
    }
  }
}

output "deepwiki_tools" {
  description = "What the platform discovered on the server"
  value       = apim_catalog_mcp_server.deepwiki.tools
}
