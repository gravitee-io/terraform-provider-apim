---
page_title: "API with pages"
subcategory: "V4 API"
---

# Create a simple V4 HTTP that contains documentation pages.

## Using a fetcher

The HTTP fetcher is used to get then and poll an Open API spec every ten minutes. The spec is added to a "Specifications" folder.

```terraform
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

```

## Using inline Markdown

The API contains a folder that contains a fake documentation page. It also has set up a home page with a simple content.

```terraform
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
            path = "/api-with-pages-inline/"
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
      hrid     = "homepage"
      content  = <<-EOT
          # Homepage
          Go to the "Markdowns" folder to find some content.
          EOT
      name     = "Home"
      homepage = true
      type     = "MARKDOWN"
    },
    {
      hrid        = "markdown"
      content     = <<-EOT
          Hello world!
          --
          This is markdown.
          EOT
      name        = "Hello Markdown"
      parent_hrid = "markdowns-folder"
      type        = "MARKDOWN"
      order       = 0
    },
    {
      hrid  = "markdowns-folder"
      name  = "Markdowns"
      type  = "FOLDER"
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

```
