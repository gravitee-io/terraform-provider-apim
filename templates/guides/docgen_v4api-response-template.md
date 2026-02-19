---
page_title: "HTTP Proxy with CORS and response templates"
subcategory: "V4 API"
---

# Create a V4 HTTP proxy API with CORS support and custom response templates

The following example configures the V4 HTTP proxy API with CORS.
If the pre-flight phase fails, it will return a 412 status code whatever all content type with bespoke body.

```terraform
resource "apim_apiv4" "v4api-proxy_response_templates" {
  # should match the resource name
  hrid            = "proxy_response_templates"
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
            path = "/proxy_response_templates/"
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
  plans = [
    {
      hrid        = "KeyLess"
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
  ]
}

```