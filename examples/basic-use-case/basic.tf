provider "apim" {

}

resource "apim_shared_policy_group" "add_context_path_header" {
  api_type = "PROXY"
  hrid     = "add_context_path_header"
  name     = "Example SPG for Terraform"
  phase    = "REQUEST"
  steps = [
    {
      configuration = jsonencode({
        scope = "REQUEST"
        addHeaders = [
          {
            name  = "X-Context-Path"
            value = "{#request.contextPath}"
          }
        ],
        removeHeaders = ["User-Agent"]
      })
      enabled = true
      name    = "add context path header, remove user agent"
      policy  = "transform-headers"
    }
  ]
}

resource "apim_apiv4" "simple_api" {
  name            = "Simple Terraform API"
  hrid            = "simple_api"
  state           = "STARTED"
  type            = "PROXY"
  version         = "1"
  visibility      = "PRIVATE"
  lifecycle_state = "PUBLISHED"
  definition_context = {}
  plans = {
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
  listeners = [
    {
      http = {
        entrypoints = [
          {
            qos  = "AUTO"
            type = "http-proxy"
          }
        ]
        paths = [
          {
            path = "/httpbin-v4-spg/"
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
            target = "http://httpbin-1.default.svc:8080"
          })
          secondary = false
        }
      ]
    }
  ]
  # warning:
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
          name    = "Example SPG for Terraform"
          policy  = "shared-policy-group-policy",
          configuration = jsonencode({
            hrid = apim_shared_policy_group.add_context_path_header.hrid
          })
        }
      ]
    }
  ]
  analytics = {
    enabled = false
  }
}
