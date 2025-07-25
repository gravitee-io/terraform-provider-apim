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
  name            = "terraform_example"
  organization_id = var.organization_id
  state           = "STOPPED"
  type            = "PROXY"
  version         = "1"
  visibility      = "PRIVATE"

  analytics = {
    enabled = true
  }

  definition_context = {}

  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      load_balancer = {
        type = "ROUND_ROBIN"
      }
      shared_configuration = jsonencode({
        proxy = {
          useSystemProxy = false
          enabled        = false
        }
        http = {
          keepAliveTimeout         = 30000
          keepAlive                = true
          followRedirects          = false
          readTimeout              = 10000
          idleTimeout              = 60000
          connectTimeout           = 3000
          useCompression           = false
          maxConcurrentConnections = 20
          version                  = "HTTP_1_1"
          pipelining               = false
        }
        ssl = {
          keyStore = {
            type = ""
          }
          hostnameVerifier = false
          trustStore = {
            type = ""
          }
          trustAll = false
        }
      })
      endpoints = [
        {
          name                  = "Default HTTP proxy"
          type                  = "http-proxy"
          weight                = 1
          inherit_configuration = true
          configuration = jsonencode({
            target = "https://example.com"
          })
          services  = {}
          secondary = false
        }
      ]
      services = {}
    }
  ]

  listeners = [
    {
      http = {
        cors = {
          allow_credentials = false
          enabled           = false
          run_policies      = false
        }
        entrypoints = [
          {
            configuration = jsonencode({})
            qos           = "AUTO"
            type          = "http-proxy"
          }
        ]
        path_mappings = []
        paths = [
          {
            path           = "/${var.hrid}/"
            overrideAccess = false
          }
        ]
        type = "HTTP"
      }
    }
  ]

  plans = {
    Keyless = {
      description = "No sec"
      mode        = "STANDARD"
      name        = "Keyless"
      status      = "PUBLISHED"
      type        = "API"
      validation  = "AUTO"

      security = {
        type = "KEY_LESS"
      }
    }
  }
}
