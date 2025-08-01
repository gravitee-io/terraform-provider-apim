resource "apim_apiv4" "simple-api" {
  # should match the resource name
  hrid            = "simple-api2"
  name            = "[Terraform] Simple PROXY API"
  description     = <<-EOT
    A simple API that routes traffic to gravitee echo API with an extra header.
    It is published to the API portal as public API and
    deployed to the Gateway
  EOT
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
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
            path = "/simple-api/"
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
          name   = "Default HTTP proxy"
          type   = "http-proxy"
          weight = 1
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
  flows = [
    {
      enabled = true
      selectors = [
        {
          http = {
            type         = "HTTP"
            path         = "/"
            pathOperator = "STARTS_WITH"
            methods = []
          }
        }
      ]
      request = [
        {
          enabled = true
          name      = "Add 1 header"
          policy = "transform-headers"
          # Configuration is JSON as the schema depends on the policy used
          configuration = jsonencode({
            scope = "REQUEST"
            addHeaders = [
              {
                name  = "X-Hello"
                value = "World!"
              }
            ]
          })
        }
      ]
    }
  ]
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
