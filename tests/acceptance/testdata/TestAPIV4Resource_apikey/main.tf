variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "name" {
  type = string
}

variable "organization_id" {
  type = string
}

resource "apim_apiv4" "test" {
  environment_id  = var.environment_id
  hrid            = var.hrid
  lifecycle_state = "UNPUBLISHED"
  name            = var.name
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
      endpoints = [
        {
          name   = "Default HTTP proxy"
          type   = "http-proxy"
          weight = 1
          configuration = jsonencode({
            target = "https://example.com"
          })
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
      hrid        = "apikey"
      description = "API Key"
      mode        = "STANDARD"
      name        = "Api Key"
      status      = "PUBLISHED"
      type        = "API"
      validation  = "AUTO"
      security = {
        type = "API_KEY"
      }
    }
  ]
}
