variable "github_mcp_token" {
  type      = string
  sensitive = true
}

# PROXY mode: one upstream MCP server fronted whole.
resource "apim_mcp_proxy" "github" {
  hrid         = "github-mcp"
  entity_id    = "mcp-proxy.github"
  name         = "GitHub MCP"
  context_path = "/mcp/github"
  mode         = "PROXY"
  proxy = {
    server_url = "https://api.githubcopilot.com/mcp/"
    upstream_auth = {
      bearer = {
        token = var.github_mcp_token
      }
    }
  }
  flows = [
    {
      name = "Rate limit tool calls"
      selectors = [
        {
          mcp = {
            methods = ["tools/call"]
          }
        }
      ]
      request = [
        {
          name   = "Quota"
          policy = "rate-limit"
          # The platform reports a configuration with the policy's defaults filled in.
          # Writing them keeps the plan empty.
          configuration = jsonencode({
            errorStrategy = "FALLBACK_PASS_TROUGH"
            rate = {
              limit          = 100
              periodTime     = 1
              periodTimeUnit = "MINUTES"
              key            = ""
              useKeyOnly     = false
            }
          })
        }
      ]
    }
  ]
  plans = [
    {
      name = "Default"
      security = {
        api_key = {
          source = "HEADER"
        }
      }
    }
  ]
}

resource "apim_catalog_mcp_server" "github" {
  # properties omitted for simplicity
  hrid = "github"
}

# STUDIO mode: tools picked from catalog MCP servers.
resource "apim_mcp_proxy" "dev-tools" {
  hrid         = "dev-tools"
  entity_id    = "mcp-proxy.dev-tools"
  name         = "Developer Tools"
  context_path = "/mcp/dev-tools"
  mode         = "STUDIO"
  studio = {
    # Sorted by server, then by tool.
    tools = [
      {
        server = apim_catalog_mcp_server.github.hrid
        tool   = "create_issue"
      },
      {
        server = apim_catalog_mcp_server.github.hrid
        tool   = "list_pull_requests"
        alias  = "list_prs"
      }
    ]
    # One entry per server a tool comes from.
    upstream_auth = [
      {
        server = apim_catalog_mcp_server.github.hrid
        auth = {
          bearer = {
            token = var.github_mcp_token
          }
        }
      }
    ]
    enable_fga = true
  }
  plans = [
    {
      name = "Default"
      security = {
        api_key = {
          source = "HEADER"
        }
      }
    }
  ]
}
