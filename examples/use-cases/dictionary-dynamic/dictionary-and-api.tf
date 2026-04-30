resource "apim_dictionary" "dynamic" {
  hrid        = "dynamic"
  name        = "[Terraform] Dynamic dictionary"
  description = "Expose all headers of Gravitee echo API as properties"
  deployed    = true
  type        = "DYNAMIC"
  dynamic = {
    provider = {
      http = {
        type   = "HTTP"
        url    = "https://api.gravitee.io/echo"
        method = "GET"
        # This header will returned and then used
        # as a property in the API policy
        headers = [
          {
            name  = "X-Test-Specific"
            value = "ABCDEF"
          }
        ]
        specification = <<-EOT
        [
          {
            "operation": "shift",
            "spec": {
              "headers": {
                "*": {
                  "$": "[#2].key",
                  "@": "[#2].value"
                }
              }
            }
          }
        ]
        EOT
      }
    }
    trigger = {
      rate = 5
      unit = "SECONDS"
    }
  }
}



resource "apim_apiv4" "proxy-with-dyn-dictionary" {
  # should match the resource name
  hrid            = "proxy-with-dyn-dictionary"
  name            = "[Terraform] Simple PROXY API with dynamic dict ref"
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
            path = "/proxy-with-dyn-dictionary/"
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
          # Configuration is JSON as the schema depends on the policy used
          configuration = jsonencode({
            scope = "REQUEST"
            addHeaders = [
              {
                name = "X-Env"
                # HRID of the dictionary is used as the key
                value = "{#dictionaries['dynamic']['X-Test-Specific']}"
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
