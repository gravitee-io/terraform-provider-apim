variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

variable "general_conditions_hrid" {
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

  pages = [
    {
      hrid      = "general_conditions"
      content   = "General conditions"
      name      = "Plan general conditions"
      type      = "MARKDOWN"
      published = true
    },
    {
      hrid      = "general_conditions_v2"
      content   = "General conditions"
      name      = "Plan general conditions v2"
      type      = "MARKDOWN"
      published = true
    },
    {
      hrid      = "unpublished_general_conditions"
      content   = "General conditions"
      name      = "Deprecated Plan general conditions"
      type      = "MARKDOWN"
      published = false
    },
  ]
  plans = [
    {
      hrid        = "Keyless"
      description = "No security"
      mode        = "STANDARD"
      name        = "No security"
      status      = "PUBLISHED"
      type        = "API"
      validation  = "AUTO"
      security = {
        type = "KEY_LESS"
      }
      general_conditions_hrid = var.general_conditions_hrid
    }
  ]
}
