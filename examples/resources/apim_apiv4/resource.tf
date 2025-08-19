resource "apim_apiv4" "quick-start-api" {
  # should match the resource name
  hrid            = "quick-start-api"
  name            = "[Terraform] Quick Start PROXY API"
  description     = "A simple API that routes traffic to gravitee echo API"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"   # API will be deployed
  lifecycle_state = "PUBLISHED" # Will be published in Portal
  visibility      = "PUBLIC"    # Will be public in the Portal
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
            path = "/quick-start-api/"
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
          # Configuration is JSON as endpoint can be custom plugins
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
  analytics = {
    enabled = false
  }
  # known limitation: will be fixed in future releases
  definition_context = {}
  plans = {
    # known limitation, key should equal name for clean terraform plans
    # will be fixed in future release
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
