provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "api_hrid" {
  type = string
}

variable "content" {
  type = string
}

variable "hrid" {
  type = string
}

variable "name" {
  type = string
}

resource "apim_apiv4" "api" {
  hrid            = var.api_hrid
  lifecycle_state = "UNPUBLISHED"
  name            = "Documentation API"
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
            path = "/${var.api_hrid}/"
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

resource "apim_documentation_api" "test" {
  api_hrid = apim_apiv4.api.hrid
  hrid     = var.hrid
  name     = var.name
  type     = "GRAVITEE_MARKDOWN"
  content  = var.content
}
