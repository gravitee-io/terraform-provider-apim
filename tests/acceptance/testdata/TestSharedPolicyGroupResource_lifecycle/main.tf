variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

resource "apim_shared_policy_group" "test" {
  api_type        = "PROXY"
  environment_id  = var.environment_id
  hrid            = var.hrid
  name            = "terraform_example"
  organization_id = var.organization_id
  phase           = "REQUEST"
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
      }),
    },
  ]
}
