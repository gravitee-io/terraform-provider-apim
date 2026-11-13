variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

# Each plugin configuration declares a few fields only: APIM stores and returns
# every other field with its default.
resource "apim_apiv4" "test" {
  environment_id  = var.environment_id
  hrid            = var.hrid
  lifecycle_state = "UNPUBLISHED"
  name            = "minimal plugin configuration"
  organization_id = var.organization_id
  state           = "STOPPED"
  type            = "PROXY"
  version         = "1"
  visibility      = "PRIVATE"

  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      load_balancer = {
        type = "ROUND_ROBIN"
      }
      # readTimeout removed: APIM keeps 5000.
      shared_configuration = jsonencode({
        http = {}
      })
      services = {}
      endpoints = [
        {
          name = "default"
          type = "http-proxy"
          configuration = jsonencode({
            target = "https://example.com"
          })
          services = {}
        }
      ]
    }
  ]

  listeners = [
    {
      http = {
        entrypoints = [
          {
            type = "http-proxy"
          }
        ]
        paths = [
          {
            path = "/${var.hrid}/"
          }
        ]
        type = "HTTP"
      }
    }
  ]

  flows = [
    {
      enabled = true
      request = [
        {
          name    = "Quota"
          enabled = true
          policy  = "quota"
          configuration = jsonencode({
            quota = {
              limit          = 1000
              periodTime     = 1
              periodTimeUnit = "HOURS"
            }
          })
        }
      ]
    }
  ]

  plans = [
    {
      hrid        = "Keyless"
      description = "No sec"
      mode        = "STANDARD"
      name        = "No security"
      status      = "PUBLISHED"
      type        = "API"
      validation  = "AUTO"
      security = {
        type = "KEY_LESS"
      }
    }
  ]
}

data "apim_apiv4" "test" {
  environment_id  = var.environment_id
  hrid            = apim_apiv4.test.hrid
  organization_id = var.organization_id
}
