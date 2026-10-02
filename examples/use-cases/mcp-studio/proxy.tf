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
