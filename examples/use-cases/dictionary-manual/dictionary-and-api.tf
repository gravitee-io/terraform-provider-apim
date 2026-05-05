resource "apim_dictionary" "manual" {
  # should match the resource name
  hrid     = "manual"
  name     = "[Terraform] Manual dictionary"
  deployed = true
  type     = "MANUAL"
  manual = {
    type = "MANUAL"
    properties = {
      env = "test"
    }
  }
}

resource "apim_apiv4" "proxy-with-dictionary" {
  # should match the resource name
  hrid            = "proxy-with-dictionary"
  name            = "[Terraform] Simple PROXY API with dict ref"
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
            path = "/proxy-with-dictionary/"
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
            type          = "HTTP"
            path          = "/"
            path_operator = "STARTS_WITH"
            methods       = []
          }
        }
      ]
      request = [
        {
          enabled = true
          name    = "Add 1 header"
          policy  = "transform-headers"
          configuration = jsonencode({
            scope = "REQUEST"
            addHeaders = [
              {
                name = "X-Env"
                # HRID of the dictionary is used as the key
                value = "{#dictionaries['manual']['env']}"
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
