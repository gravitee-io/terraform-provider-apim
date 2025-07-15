resource "apim_shared_policy_group" "ratelimit-policy" {
  # should match the resource name
  hrid        = "ratelimit-policy"
  name        = "[Terraform] Rate limit shared policy"
  api_type    = "PROXY"
  description = "Single step rate limiting policy group"
  phase       = "REQUEST"
  steps = [
    {
      enabled     = true
      description = "Limit traffic to 10 request per minute"
      name        = "Rate Limit 10"
      policy      = "rate-limit"
      configuration = jsonencode({
        addHeaders = true
        async      = false
        rate = {
          key            = ""
          limit          = 10
          periodTime     = 1
          periodTimeUnit = "MINUTES"
          useKeyOnly     = false
        }
      })
    },
  ]
}