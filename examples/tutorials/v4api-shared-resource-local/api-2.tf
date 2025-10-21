resource "apim_apiv4" "simple-api-shared-resource-2" {
  # should match the resource name
  hrid            = "simple-api-shared-resource-2"
  name            = "[Terraform] Simple API with shared resource [2/2]"
  description     = "A simple API that routes traffic to gravitee whattimeisit API. Using basic auth configured in a shared resource"
  version         = "1"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PRIVATE"
  lifecycle_state = "UNPUBLISHED"
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
            path = "/simple-api-shared-resource-2/"
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
          name = "Default HTTP proxy"
          type = "http-proxy"
          # Configuration is JSON as it is owned by the "http-proxy" endpoint plugin
          configuration = jsonencode({
            target = "https://api.gravitee.io/whattimeisit"
          })
        }
      ]
    }
  ]
  flows = [
    {
      enabled = true
      selectors = [
        {
          http = {
            type         = "HTTP"
            path         = "/"
            pathOperator = "STARTS_WITH"
          }
        }
      ]
      request = [
        {
          # Authentication policy
          name    = "Basic Authentication",
          enabled = true,
          policy  = "policy-basic-authentication",
          # Configuration is JSON as is depends on the
          configuration = jsonencode({
            authenticationProviders = [
              "In memory users"
            ]
            realm = "gravitee.io"
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
      # Where configuraiton file is included in the API resource
      configuration = data.local_file.api-resource-basic-auth.content
    }
  ]
  analytics = {
    enabled = false
  }
  plans = [
    {
      hrid        = "keyless"
      name        = "Key Less"
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
