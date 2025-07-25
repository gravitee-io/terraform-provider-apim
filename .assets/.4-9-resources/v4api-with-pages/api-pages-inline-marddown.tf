resource "apim_apiv4" "api-with-pages-inline" {
  # should match the resource name
  hrid            = "api-with-pages-inline"
  name            = "Terraform: Simple PROXY API, Page with inlined Markdown"
  description     = "Simple proxy API containing an inline documentation written in Markdown"
  version         = "1"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PRIVATE"
  lifecycle_state = "PUBLISHED"
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
            path = "/api-with-pages-inline"
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
  pages = {
    markdown = {
      content = <<-EOT
          Hello world!
          --
          This is markdown.
          EOT
      name    = "hello-markdown"
      parent  = "markdowns-folder"
      type    = "MARKDOWN"
    }
    markdowns-folder = {
      name = "markdowns"
      type = "FOLDER"
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
