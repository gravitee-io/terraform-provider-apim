---
page_title: "Simple HTTP Proxy"
subcategory: "V4 API"
---

# Create a V4 HTTP proxy API with CORS support and custom response templates

The following example configures the V4 HTTP proxy API with CORS.
If the pre-flight phase fails, it will return a 412 status code whatever all content type with bespoke body.

```HCL
resource "apim_apiv4" "simple-api-response-template" {
  # should match the resource name
  hrid            = "simple-api-response-template"
  name            = "[Terraform] Simple PROXY API Response template"
  description     = "A simple API with CORS and a custom response template for pre-flight errors"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
  listeners = [
    {
      http = {
        type = "HTTP"
        cors = {
          enabled           = true
          allow_credentials = true
          allow_headers     = ["x-gravitee-test"]
          allow_methods     = ["POST", "GET"]
          allow_origin      = ["https://mydomain.com"]
          run_policies      = true
        }
        entrypoints = [
          {
            type = "http-proxy"
          }
        ]
        paths = [
          {
            path = "/simple-api-response-template/"
          }
        ]
      }
    }
  ]
  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      load_balancer = {
        type = "ROUND_ROBIN"
      }
      endpoints = [
        {
          name                  = "Default HTTP proxy"
          type                  = "http-proxy"
          weight                = 1
          inherit_configuration = false
          # Configuration is JSON as is depends on the type schema
          configuration = jsonencode({
            target = "https://api.gravitee.io/echo"
          })
        }
      ]
    }
  ]
  flow_execution = {
    mode           = "DEFAULT"
    match_required = false
  }
  flows = []
  response_templates = {
    CORS_PREFLIGHT_FAILED : {
      "*/*" : {
        status = 412
        headers = {
          X-Error = "Cors preflight"
        }
        body                        = "Custom CORS error message"
        propagate_error_key_to_logs = false
      }
    }
  }
  analytics = {
    enabled = false
  }
  # known limitation: will dispear in next version
  definition_context = {}
  plans = {
    # known limitation, key have to match name to avoid terraform plan to remain inconsistent
    KeyLess = {
      name        = "KeyLess"
      type        = "API"
      mode        = "STANDARD"
      validation  = "AUTO"
      status      = "PUBLISHED"
      description = "This plan does not require any authentication"
      security = {
        type = "KEY_LESS"
      }
    }
  }
}

```