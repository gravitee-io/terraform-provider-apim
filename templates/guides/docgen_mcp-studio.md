---
page_title: "MCP proxy and studio"
subcategory: "MCP"
---

# MCP servers, proxies and studios

This example declares what the AI Management module manages for MCP:

1. An upstream MCP server registered in the AI Catalog (`apim_catalog_mcp_server`)
2. A proxy fronting that server whole (`apim_mcp_proxy`, `mode = "PROXY"`)
3. A studio exposing tools picked from the catalog (`apim_mcp_proxy`, `mode = "STUDIO"`)

It requires APIM 4.13 with the AI Management module and a license that includes it.

-> `apim_mcp_proxy` is the resource the Console's MCP wizard creates. `apim_apiv4` with `type = "MCP_PROXY"` declares
the API underneath by hand, with its listeners and endpoint groups, and knows nothing of the catalog, of tool selection
or of tool-level authorization.

## Catalog MCP server

The platform discovers the server when the resource is applied. An upstream it cannot reach fails the apply and
registers nothing. The discovered tools, prompts and resources are read-only attributes: each carries the `entity_id`
an authorization policy names, such as `mcp-tool.example-deepwiki.ask_wiki_question`.

```terraform
# An upstream MCP server registered in the AI Catalog. The platform discovers it
# when the resource is applied: its tools, prompts and resources come back in the
# state, each with the entity id an authorization policy names.
resource "apim_catalog_mcp_server" "deepwiki" {
  hrid        = "example-deepwiki"
  entity_id   = "mcp-server.example-deepwiki"
  description = "DeepWiki documentation tools for coding agents"
  connection = {
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

```

~> The platform discovers a server again only when its `connection` changes. A tool added upstream afterwards does not
reach the state on a refresh.

## Proxy

`mode` and `entity_id` cannot be changed on an existing proxy. Every other attribute converges on apply.

```terraform
# PROXY mode: the upstream server fronted whole, behind a keyless plan. The
# platform does not contact the upstream when the proxy is applied; the gateway
# reaches it at runtime.
resource "apim_mcp_proxy" "deepwiki" {
  hrid         = "example-deepwiki-proxy"
  entity_id    = "mcp-proxy.example-deepwiki"
  name         = "[Terraform] DeepWiki"
  context_path = "/example/mcp/deepwiki"
  mode         = "PROXY"
  proxy = {
    server_url = "https://mcp.deepwiki.com/mcp"
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
          # The platform reports a configuration with the policy's defaults filled
          # in. Writing them keeps the plan empty.
          configuration = jsonencode({
            errorStrategy = "FALLBACK_PASS_TROUGH"
            rate = {
              limit          = 10
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
      name = "Open"
      security = {
        key_less = {}
      }
    }
  ]
}

```

Plans converge by name: a plan added to the list is created and published, a plan removed from it is closed, which
ends its subscriptions, and the `security` of an existing plan cannot be changed.

## Studio

```terraform
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

```

## Order of the lists

The platform reports `plans` and `identity_providers` sorted by name, `studio.tools` sorted by server then tool, and
`studio.upstream_auth` sorted by server. Declare them in that order: another order is refused when planning, with the
expected order in the message.

## Credentials

The platform stores a credential as written and never returns it, so an imported resource needs its credentials declared
before the next apply. On a proxy or a studio, `upstream_auth` also takes a `secret://` URI pointing at a secret provider
configured on the platform: the gateway resolves it at runtime. The credential of a catalog server is the one the
platform itself presents to discover the server.
