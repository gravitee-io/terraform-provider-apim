variable "apis" {
  type = list(object({
    api      = string
    location = string
    order    = number
  }))
}

variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

variable "portal_hrid" {
  type = string
}

resource "apim_portal" "primary" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = var.portal_hrid
  name            = "Primary Portal"
  navigation = [
    {
      path = "/apis-${var.portal_hrid}"
    },
    {
      path = "/apis-${var.portal_hrid}/featured-${var.portal_hrid}"
    }
  ]
}

resource "apim_apiv4" "first" {
  environment_id  = var.environment_id
  hrid            = "api-a-${var.hrid}"
  lifecycle_state = "UNPUBLISHED"
  name            = "terraform_example_a"
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
            target = "https://example.com/a"
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
            path           = "/api-a/${var.hrid}/"
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

resource "apim_apiv4" "second" {
  environment_id  = var.environment_id
  hrid            = "api-b-${var.hrid}"
  lifecycle_state = "UNPUBLISHED"
  name            = "terraform_example_b"
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
            target = "https://example.com/b"
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
            path           = "/api-b/${var.hrid}/"
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
  portal_hrid     = apim_portal.primary.hrid
  hrid            = var.hrid
  apis = [for entry in var.apis : {
    # looks odd, but we do this so that we have stable test where the API will be created first without "depends on"
    api_hrid = entry.api == "first" ? apim_apiv4.first.hrid : apim_apiv4.second.hrid
    location = entry.location
    order    = entry.order
  }]
}
