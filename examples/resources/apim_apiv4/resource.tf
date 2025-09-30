resource "apim_apiv4" "example" {
  # should match the resource name
  hrid            = "example"
  name            = "[Terraform] Example API"
  description     = <<-EOT
    A example API that routes traffic to gravitee echo API with an extra header.
    It uses basic authentication with johndoe/unknown as credentials.
    It is published to the API portal as public API and deployed to the Gateway.
    Some analytics are configured and documentation page is configured.
  EOT
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
  labels = [
    "example",
    "proxy",
    "terraform"
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
            path = "/example/"
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
  properties = [
    {
      key   = "hello",
      value = "World!",
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
            methods      = []
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
                name  = "X-Hello"
                value = "{#api.properties['hello']}"
              }
            ]
          })
        }
      ]
    }
  ]
  resources = [
    {
      enabled = true
      name    = "In memory users"
      type    = "auth-provider-inline-resource"
      configuration = jsonencode({
        users = [
          {
            username = "johndoe"
            password = "unknown"
            roles    = []
          }
        ]
      })
    }
  ]
  plans = [
    {
      hrid        = "keyless"
      name        = "No security"
      type        = "API"
      mode        = "STANDARD"
      validation  = "AUTO"
      status      = "PUBLISHED"
      description = "This plan does not require any authentication"
      security = {
        type = "KEY_LESS"
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
                methods      = []
              }
            }
          ]
          request = [
            {
              # Authentication policy
              "name" : "Basic Authentication",
              "enabled" : true,
              "policy" : "policy-basic-authentication",
              "configuration" : jsonencode({
                "authenticationProviders" = [
                  "In memory users"
                ]
                "realm" = "gravitee.io"
              })
            }
          ]
        }

      ]
    }
  ]
  analytics = {
    enabled = true
    logging = {
      condition = "{#request.headers['Accept'][0] == '*/*'}"
      content = {
        headers         = true
        messageHeaders  = false
        payload         = true
        messagePayload  = false
        messageMetadata = false
      }
      phase = {
        request  = true
        response = true
      }
      mode = {
        endpoint   = true
        entrypoint = true
      }
    }
    tracing = {
      enabled = true
      verbose = true
    }
  }
  pages = [
    {
      hrid     = "homepage"
      name     = "Home"
      content  = <<-EOT
          # Homepage
          Terraform example API document home page.
          From now on only your imagination is the limit.
          EOT
      homepage = true
      type     = "MARKDOWN"
      order    = 0
    }
  ]
}
