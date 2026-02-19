
resource "apim_apiv4" "shared_endpoint_config" {
  hrid            = "shared_endpoint_config"
  lifecycle_state = "UNPUBLISHED"
  name            = "[Terraform] Proxy with shared endpoint configuration"
  state           = "STOPPED"
  type            = "PROXY"
  version         = "1"
  visibility      = "PRIVATE"
  listeners = [
    {
      http = {
        type = "HTTP"
        entrypoints = [
          {
            type          = "http-proxy"
          }
        ]
        paths = [
          {
            path           = "/shared-endpoint-config/"
          }
        ]
      }
    }
  ]
  endpoint_groups = [
    {
      name = "Default HTTP proxy group with shared configuration"
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
            target               = "/"
          })
        }
      }
      shared_configuration = jsonencode({
        proxy = {
          useSystemProxy = true
          enabled        = true
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
            target = "https://api.gravitee.io/echo"
          })
        }
      ]
    }
  ]
  flow_execution = {
    mode           = "DEFAULT"
    match_required = false
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
