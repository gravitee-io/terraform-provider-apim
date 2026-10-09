---
page_title: "API Portal Navigation"
subcategory: "V4 API"
---

# API Internal Portal Navigation

This example demonstrates how to declare an API's internal documentation navigation tree for the next-gen portal using the `portal_navigation` attribute on `apim_apiv4`.

Paths are ordered — the order in the list is preserved.
Intermediate folders are implicitly created when child paths are declared.

```terraform
resource "apim_apiv4" "with-portal-navigation" {
  # should match the resource name
  hrid            = "with-portal-navigation"
  name            = "[Terraform] API with portal navigation"
  description     = "V4 API declaring an internal documentation tree for the next-gen portal"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
  portal_navigation = [
    {
      path         = "/api-docs"
      display_name = "API Docs"
      order        = 1
    },
    {
      path         = "/api-docs/reference"
      display_name = "Reference"
      order        = 1
    }
  ]
  listeners = [
    {
      http = {
        type = "HTTP"
        entrypoints = [
          {
            type = "http-proxy"
          }
        ]
        paths = [
          {
            path = "/with-portal-navigation/"
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
          configuration = jsonencode({
            target = "https://api.gravitee.io/echo"
          })
        }
      ]
    }
  ]
  analytics = {
    enabled = false
  }
  plans = [
    {
      hrid        = "KeyLess"
      name        = "No security"
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
