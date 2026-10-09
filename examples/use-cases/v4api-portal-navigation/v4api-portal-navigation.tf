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
