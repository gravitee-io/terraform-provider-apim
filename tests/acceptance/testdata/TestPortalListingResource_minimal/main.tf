variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

resource "apim_portal" "test" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "portal-${var.hrid}"
  name            = "Test Portal"
  navigation = [
    {
      path = "/apis-${var.hrid}"
    }
  ]
}

resource "apim_apiv4" "test" {
  environment_id  = var.environment_id
  hrid            = "api-${var.hrid}"
  lifecycle_state = "UNPUBLISHED"
  name            = "terraform_example"
  organization_id = var.organization_id
  state           = "STARTED"
  type            = "PROXY"
  version         = "1"
  visibility      = "PRIVATE"
  analytics = {
    enabled = true
  }
  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      endpoints = [
        {
          name = "Default HTTP proxy"
          type = "http-proxy"
          configuration = jsonencode({
            target = "https://example.com"
          })
        }
      ]
      services = {}
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
            path           = "/api/${var.hrid}/"
            overrideAccess = false
          }
        ]
        type = "HTTP"
      }
    }
  ]
  plans = [
    {
      hrid       = "default"
      name       = "Default Plan"
      type       = "API"
      mode       = "STANDARD"
      validation = "AUTO"
      status     = "PUBLISHED"
      security = {
        type = "KEY_LESS"
      }
    }
  ]
}

resource "apim_portal_listing" "test" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  portal_hrid     = apim_portal.test.hrid
  hrid            = var.hrid
  apis = [
    {
      api_hrid = apim_apiv4.test.hrid
      location = "/apis-${var.hrid}"
    }
  ]
}
