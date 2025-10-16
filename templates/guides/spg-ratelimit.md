---
page_title: "Rate limiting example"
subcategory: "Shared Policy Group"
---

# Create a Shared Policy Group for rate limiting

The following example configures the Shared Policy Group resource to use the Gravitee
[Rate Limit policy](https://documentation.gravitee.io/apim/create-and-configure-apis/apply-policies/policy-reference/rate-limit).
It applies rate limiting on the Request phase to limit traffic to 10 requests per minute.

```terraform
resource "apim_shared_policy_group" "simple" {
  # should match the resource name
  hrid        = "simple"
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
          key            = "rl"
          limit          = 10
          periodTime     = 1
          periodTimeUnit = "MINUTES"
          useKeyOnly     = false
        },
        errorStrategy = "FALLBACK_PASS_TROUGH"
      })
    },
  ]
}
```