resource "apim_apiv4" "simple-api-subscribed-apikey" {
  # should match the resource name
  hrid            = "simple-api-subscribed-apikey"
  name            = "[Terraform] Simple PROXY API With API KEY Plan"
  description     = "A simple API with API KEY meant to subscribed using Terraform"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
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
            path = "/simple-api-subscribed-apikey"
          }
        ]
        type = "HTTP"
      }
    }
  ]
  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      endpoints = [
        {
          name = "Default HTTP proxy"
          type = "http-proxy"
          configuration = jsonencode({
            target = "https://api.gravitee.io/echo"
          })
        }
      ]
      services = {}
    }
  ]
  analytics = {
    enabled = true
  }
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


resource "apim_application" "simple-app-subscribed-apikey" {
  # should match the resource name
  hrid        = "simple-app-subscribed-apikey"
  name        = "[Terraform] Simple Application"
  description = "Subscription tests application"
}

resource "apim_subscription" "simple-subscription-apikey" {
  # should match the resource name
  hrid             = "simple-subscription-apikey"
  api_hrid         = apim_apiv4.simple-api-subscribed-apikey.hrid
  plan_hrid        = apim_apiv4.simple-api-subscribed-apikey.plans[0].hrid
  application_hrid = apim_application.simple-app-subscribed-apikey.hrid
  api_keys = [
    {
      key       = "custom-user-typed-api-key-0123456789",
      expire_at = "2042-12-31T00:00:00Z"
    }
  ]
}
