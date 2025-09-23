variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

resource "apim_apiv4" "test" {
  environment_id  = var.environment_id
  hrid            = var.hrid
  lifecycle_state = "UNPUBLISHED"
  name            = "minimal"
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
      services = {}
      endpoints = [
        {
          name = "Default HTTP proxy"
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
