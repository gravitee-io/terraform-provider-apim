resource "apim_apiv4" "api-with-pages-fetcher" {
  # should match the resource name
  hrid            = "api-with-pages-fetcher"
  name            = "Terraform: Simple PROXY API, Page from HTTP fetcher"
  description     = "Simple proxy API containing PetStore Swagger API spec fetched from swagger website"
  version         = "1"
  type            = "PROXY"
  state           = "STOPPED"
  visibility      = "PRIVATE"
  lifecycle_state = "UNPUBLISHED"
  listeners = [
    {
      http = {
        type = "HTTP"
        entrypoints = [
          {
            type = "http-proxy"
          }
        ]
        paths = [
          {
            path = "/api-with-pages-fetcher/"
          }
        ]
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
      endpoints = [
        {
          name                  = "Default HTTP proxy"
          type                  = "http-proxy"
          weight                = 1
          inherit_configuration = false
          # Configuration is JSON as is depends on the type schema
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
  flows = []
  analytics = {
    enabled = false
  }
  pages = [
    {
      hrid  = "docs-folder"
      name  = "Specifications"
      type  = "FOLDER"
      order = 0
    },
    {
      hrid        = "swagger"
      name        = "Pet Store"
      parent_hrid = "docs-folder"
      source = {
        configuration = jsonencode({
          fetchCron      = "*/10 * * * * *"
          url            = "https://petstore.swagger.io/v2/swagger.json"
          autoFetch      = false
          useSystemProxy = false
        })
        type = "http-fetcher"
      }
      type  = "SWAGGER"
      order = 1
    }
  ]

  plans = [
    {
      hrid        = "KeyLess"
      name        = "No security"
      type        = "API"
      mode        = "STANDARD"
      validation  = "AUTO"
      status      = "PUBLISHED"
      description = "This plan does not require any authentication"
      security = {
        type = "KEY_LESS"
      }
    }
  ]

}
