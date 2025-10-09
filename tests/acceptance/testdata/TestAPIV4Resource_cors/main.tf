variable "environment_id" {
  type    = string
  default = "DEFAULT"
}

variable "hrid" {
  type    = string
  default = "coucou"
}

variable "organization_id" {
  type    = string
  default = "DEFAULT"
}

variable "allow_headers" {
  type    = list(string)
  default = []
}

variable "allow_methods" {
  type    = list(string)
  default = []
}

variable "allow_origin" {
  type    = list(string)
  default = []
}

variable "expose_headers" {
  type    = list(string)
  default = []
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
  labels = [
    "test",
    "all props"
  ]
  listeners = [
    {
      http = {
        cors = {
          allow_credentials = true
          allow_headers     = var.allow_headers
          allow_methods     = var.allow_methods
          allow_origin      = var.allow_origin
          enabled           = true
          expose_headers    = var.expose_headers
          max_age           = 100
          run_policies      = true
        }
        entrypoints = [
          {
            configuration = jsonencode({})
            qos           = "AUTO"
            type          = "http-proxy"
          }
        ]
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
  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      load_balancer = {
        type = "ROUND_ROBIN"
      }
      services = {
        "health_check" : {
          enabled = true
          type    = "http-health-check"
          configuration = jsonencode({
            schedule = "*/10 * * * * *"
            headers = [
              {
                name  = "X-HC-Data"
                value = "check"
              }
            ]
            overrideEndpointPath = true
            method               = "GET"
            failureThreshold     = 3
            assertion            = "{#response.status < 300}"
            successThreshold     = 3
            target               = "/test"
          })
        }
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
          secondary = false
          services = {
            "health_check" : {
              enabled               = true
              type                  = "http-health-check"
              overrideConfiguration = true
              configuration = jsonencode({
                schedule = "*/10 * * * * *"
                headers = [
                  {
                    name  = "X-HC-Data"
                    value = "check override"
                  }
                ]
                overrideEndpointPath = true
                method               = "GET"
                failureThreshold     = 3
                assertion            = "{#response.status < 300}"
                successThreshold     = 3
                target               = "/test"
              })
            }
          }
          tenants : ["foo"]
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
            username = "admin"
            password = "admin"
            roles    = []
          },
          {
            username = "user"
            password = "password"
            roles    = []
          }
        ]
      })
    }
  ]
  response_templates = {
    INVALID_HTTP_METHOD = {
      "*/*" : {
        status = 400
        headers = {
          X-Error    = "invalid method"
          X-Key      = "INVALID_HTTP_METHOD"
          X-Encoding = "all"
        }
        body                        = "http method override denied"
        propagate_error_key_to_logs = false
      },
    }
    CORS_PREFLIGHT_FAILED : {
      "*/*" : {
        status = 412
        headers = {
          X-Error   = "Cors preflight"
          X-Details = "Cors preflight"
        }
        body                        = "Custom CORS error message"
        propagate_error_key_to_logs = false
      },
      "application/json" : {
        status = 412
        headers = {
          X-Details = "Cors preflight, JSON"
          X-Error   = "Cors preflight"
        }
        body                        = "Custom CORS error message"
        propagate_error_key_to_logs = false
      }
    }
    NO_ENDPOINT_FOUND : {
      "application/json" : {
        status = 404
        headers = {
          X-Encoding = "JSON"
          X-Key      = "NO_ENDPOINT_FOUND"
          X-Error    = "endpoint not found"
        }
        body                        = "Custom endpoint not found error message"
        propagate_error_key_to_logs = false
      },
      "*/*" : {
        status = 404
        headers = {
          X-Error    = "endpoint not found"
          X-Key      = "NO_ENDPOINT_FOUND"
          X-Encoding = "all"
        }
        body                        = "Custom endpoint not found error message"
        propagate_error_key_to_logs = false
      },
    }

  }
  plans = [
    {
      hrid        = "Keyless"
      description = "No security"
      mode        = "STANDARD"
      name        = "Keyless"
      status      = "PUBLISHED"
      type        = "API"
      validation  = "AUTO"
      security = {
        type = "KEY_LESS"
      }
    }
  ]
}
