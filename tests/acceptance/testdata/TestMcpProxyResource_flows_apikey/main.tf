provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "limit" {
  type = number
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
          # The policy's defaults are written out: the platform reports them.
          configuration = jsonencode({
            errorStrategy = "FALLBACK_PASS_TROUGH"
            rate = {
              limit          = var.limit
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
      name = "keys"
      security = {
        api_key = {
          source = "HEADER"
        }
      }
      flows = [
        {
          name = "Plan quota"
          request = [
            {
              name   = "Plan quota"
              policy = "rate-limit"
              configuration = jsonencode({
                errorStrategy = "FALLBACK_PASS_TROUGH"
                rate = {
                  limit          = var.limit * 10
                  periodTime     = 1
                  periodTimeUnit = "HOURS"
                  key            = ""
                  useKeyOnly     = false
                }
              })
            }
          ]
        }
      ]
    }
  ]
}
