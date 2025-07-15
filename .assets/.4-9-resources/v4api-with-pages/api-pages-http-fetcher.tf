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
            path = "/api-with-pages-fetcher"
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
          name   = "Default HTTP proxy"
          type   = "http-proxy"
          weight = 1
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
  pages = {
    docs-folder = {
      name = "specifications"
      type = "FOLDER"
    }
    swagger = {
      name   = "pet-store"
      parent = "docs-folder"
      source = {
        configuration = {
          fetchCron = "*/10 * * * * *"
          url       = "https://petstore.swagger.io/v2/swagger.json"
        }
        type = "http-fetcher"
      }
      type = "SWAGGER"
    }
  }
  # known limitation: will dispear in next version
  definition_context = {}
  plans = {
    # known limitation, key have to match name to avoid terraform plan to remain inconsistent
    KeyLess = {
      name        = "KeyLess"
      type        = "API"
      mode        = "STANDARD"
      validation  = "AUTO"
      status      = "PUBLISHED"
      description = "This plan does not require any authentication"
      security = {
        type = "KEY_LESS"
      }
    }
  }
}
