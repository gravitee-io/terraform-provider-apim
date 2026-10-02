# STUDIO mode: a tool surface picked from catalog MCP servers. A tool is named by
# the hrid of its server and the name the server advertises.
resource "apim_mcp_proxy" "docs" {
  hrid         = "example-docs-studio"
  entity_id    = "mcp-proxy.example-docs-studio"
  name         = "[Terraform] Docs Studio"
  context_path = "/example/mcp/docs-studio"
  mode         = "STUDIO"
  studio = {
    # Sorted by server, then by tool: the order the platform reports.
    tools = [
      {
        server = apim_catalog_mcp_server.deepwiki.hrid
        tool   = "ask_wiki_question"
        alias  = "ask"
      },
      {
        server = apim_catalog_mcp_server.deepwiki.hrid
        tool   = "read_wiki_structure"
      }
    ]
    # One entry per server a tool comes from, none included.
    upstream_auth = [
      {
        server = apim_catalog_mcp_server.deepwiki.hrid
        auth = {
          none = {}
        }
      }
    ]
    # Runs the platform's authorization policy on tools/call. The policies it
    # evaluates are managed in the Authorization module.
    enable_fga = true
  }
  plans = [
    {
      name = "Keys"
      security = {
        api_key = {
          source = "HEADER"
        }
      }
    }
  ]
}
