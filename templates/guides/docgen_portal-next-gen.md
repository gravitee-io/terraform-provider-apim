---
page_title: "Next-gen Portal"
subcategory: "Portal"
---

# End-to-end Next-gen Portal

This example ties together the next-gen portal resources:

1. A portal with navigation paths (`apim_portal`) and a listing that publishes APIs (`apim_portal_listing`)
2. Portal and API documentation pages (`apim_documentation_portal`, `apim_documentation_api`)
3. A V4 API with its own internal portal navigation tree (`apim_apiv4`)

~> **WARNING:** Multiple portals are not supported make sure you create one portal per environment.

## Portal and listing

The portal defines the top-level navigation tree for the developer portal. The listing publishes APIs into that tree at a chosen location.

```terraform
resource "apim_portal" "developer-portal" {
  hrid = "developer-portal"
  name = "[Terraform] Next-gen Developer Portal"
  navigation = [
    {
      path         = "/examples/docs"
      display_name = "Documentation"
    },
    {
      path         = "/examples"
      display_name = "Examples"
    },
    {
      path         = "/examples/apis"
      display_name = "APIs"
    },
    {
      path         = "/examples/archives"
      display_name = "Old stuffs"
    },
  ]
}

resource "apim_portal_listing" "public-apis" {
  portal_hrid = apim_portal.developer-portal.hrid
  hrid        = "public-apis"
  apis = [
    {
      api_hrid = apim_apiv4.pets.hrid
      location = "/examples/apis"
      order    = 1
    }
  ]
}

```

## Documentations

Portal documentation pages are bound to the portal and appear under portal navigation paths. API documentation pages are bound to an API and appear under that API's internal `portal_navigation` tree.

```terraform
resource "apim_documentation_portal" "getting-started-docs" {
  portal_hrid = apim_portal.developer-portal.hrid
  hrid        = "getting-started-docs"
  name        = "Getting Started"
  type        = "GRAVITEE_MARKDOWN"
  location    = "/examples/docs"
  order       = 1
  content     = <<-EOT
  # Getting Started

  Portal-bound documentation page managed with `apim_documentation_portal`.

  Browse the **APIs** section to discover published APIs.
  EOT
}

resource "apim_documentation_portal" "disclaimer" {
  portal_hrid = apim_portal.developer-portal.hrid
  hrid        = "disclaimer-docs"
  name        = "Disclaimer"
  type        = "GRAVITEE_MARKDOWN"
  location    = "/examples/docs"
  order       = 2
  content     = <<-EOT
  # Disclaimer

  Read the [documentation](https://documentation.gravitee.io/apim/developer-portal/new-developer-portal/customize-the-navigation/portal-automation) first!

  _Portal-bound documentation page managed with `apim_documentation_portal`._

  EOT
}

resource "apim_documentation_api" "pets-api-docs" {
  api_hrid = apim_apiv4.pets.hrid
  hrid     = "pets-api-docs"
  name     = "Markdown"
  type     = "GRAVITEE_MARKDOWN"
  location = "/md"
  order    = 2
  content  = <<-EOT
  # Pets API How-To

  Just call it!

  _API-bound documentation page managed with `apim_documentation_api`._
  EOT
}

resource "apim_documentation_api" "pets-api-oas" {
  api_hrid = apim_apiv4.pets.hrid
  hrid     = "pets-api-oas"
  name     = "Open API Spec"
  type     = "OPENAPI"
  location = "/oas"
  order    = 3
  content  = <<-EOT
openapi: 3.1.0
info:
  title: Pets API
  description: The famous Pets API
  version: 1.0.0
servers:
  - url: "https://petstore.swagger.io/v2"
tags:
  - name: Pets
paths:
  /pet:
    put:
      tags:
        - Pets
      summary: Update an existing pet.
      description: Update an existing pet by Id.
      operationId: updatePet
  EOT
}

```

## API

The V4 API declares its own internal documentation navigation with `portal_navigation`. API documentation pages (see above) are placed under those paths.

```terraform
resource "apim_apiv4" "pets" {
  hrid            = "pets"
  name            = "[Terraform] Pets API"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
  portal_navigation = [
    {
      path         = "/md"
      display_name = "How to"
    },
    {
      path         = "/oas"
      display_name = "Pets OAS"
    }
  ]
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
            path = "/pets/"
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
          configuration = jsonencode({
            target = "https://api.gravitee.io/echo"
          })
        }
      ]
    }
  ]
  analytics = {
    enabled = false
  }
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
