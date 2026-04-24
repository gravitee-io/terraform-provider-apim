variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

resource "apim_apiv4" "test" {
  environment_id  = var.environment_id
  hrid            = "api-${var.hrid}"
  lifecycle_state = "UNPUBLISHED"
  name            = "terraform_example"
  organization_id = var.organization_id
  state           = "STARTED"
  type            = "PROXY"
  version         = "1"
  visibility      = "PRIVATE"
  analytics = {
    enabled = true
  }
  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      endpoints = [
        {
          name = "Default HTTP proxy"
          type = "http-proxy"
          configuration = jsonencode({
            target = "https://example.com"
          })
        }
      ]
      services = {}
    }
  ]
  listeners = [
    {
      http = {
        entrypoints = [
          {
            type = "http-proxy"
          }
        ]
        paths = [
          {
            path           = "/api/${var.hrid}/"
            overrideAccess = false
          }
        ]
        type = "HTTP"
      }
    }
  ]
  plans = [
    {
      hrid        = "apikey"
      description = "API Key"
      mode        = "STANDARD"
      name        = "Api Key"
      status      = "PUBLISHED"
      type        = "API"
      validation  = "AUTO"
      security = {
        type = "API_KEY"
      }
    }
  ]
}

resource "apim_application" "test" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "app-${var.hrid}"
  name            = "terraform example"
  description     = "Subscription tests application"
}

resource "apim_subscription" "test" {
  environment_id   = var.environment_id
  organization_id  = var.organization_id
  hrid             = "test"
  api_hrid         = apim_apiv4.test.hrid
  plan_hrid        = apim_apiv4.test.plans[0].hrid
  application_hrid = apim_application.test.hrid
  custom_api_key   = "foobar-api-key"
}
