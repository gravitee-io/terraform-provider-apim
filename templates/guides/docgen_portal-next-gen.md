---
page_title: "Next-gen Portal"
subcategory: "Portal"
---

# End-to-end Next-gen Portal

This example ties together the next-gen portal resources:

1. A portal with navigation paths (`apim_portal`) and a listing that publishes APIs (`apim_portal_listing`)
2. A homepage, and portal and API documentation pages (`apim_documentation_portal`, `apim_documentation_api`)
3. A V4 API with its own internal portal navigation tree (`apim_apiv4`)
4. A theme the portal uses (`apim_portal_theme`)
5. Links to external sites, in the portal navigation and in the API's documentation (`apim_portal_link`, `apim_api_link`)

-> **NOTE:** These resources require APIM 4.13 or later.

~> **WARNING:** Multiple portals are not supported make sure you create one portal per environment.

## Portal and listing

The portal defines the top-level navigation tree for the developer portal. The listing publishes APIs into that tree at a chosen location.

```terraform
resource "apim_portal" "developer-portal" {
  hrid = "developer-portal"
  name = "[Terraform] Next-gen Developer Portal"
  # Enables this theme and disables the other themes of the environment
  active_theme_hrid = apim_portal_theme.brand.hrid
  structure = { top_navbar = [
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
  ] }
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

A portal page with `area = "HOMEPAGE"` is the portal's homepage instead: it has no location, and it replaces the homepage created in the Console. A portal has one homepage, so a second one with another `hrid` fails. Destroying it sends consumers to the catalog.

```terraform
resource "apim_documentation_portal" "homepage" {
  portal_hrid = apim_portal.developer-portal.hrid
  hrid        = "homepage"
  name        = "Home"
  type        = "GRAVITEE_MARKDOWN"
  # Replaces the Console or seeded homepage. A homepage has no location.
  area    = "HOMEPAGE"
  content = <<-EOT
  # Welcome

  Homepage of the developer portal, managed with `apim_documentation_portal`.
  EOT
}

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

## Theme

A theme holds the colors, font, custom CSS, and images of the portal. Applying a theme stores it without activating it: the portal's `active_theme_hrid` selects it, enables it, and disables the other themes of the environment. Removing `active_theme_hrid` re-enables another theme of the environment, or the default theme when there is none.

```terraform
resource "apim_portal_theme" "brand" {
  hrid = "brand-theme"
  name = "[Terraform] Brand theme"
  definition = {
    custom_css = "body { line-height: 1.5; }"
    font       = { font_family = "Arial, sans-serif" }
    color = {
      primary    = "#146b50"
      secondary  = "#344256"
      tertiary   = "#7a8899"
      error      = "#b3261e"
      background = { page = "#ffffff", card = "#f2f4f6" }
    }
  }
}

```

## Links

A link opens a URL from the navigation. A portal link sits in the portal navigation at its `location`. An API link belongs to the API: its `location` is a path of the API's `portal_navigation`, and it appears under the API wherever the API is listed.

```terraform
resource "apim_portal_link" "gravitee-docs" {
  portal_hrid = apim_portal.developer-portal.hrid
  hrid        = "gravitee-docs"
  name        = "Gravitee documentation"
  href        = "https://documentation.gravitee.io"
  location    = "/examples/docs"
  order       = 3
}

resource "apim_api_link" "pets-status" {
  api_hrid = apim_apiv4.pets.hrid
  hrid     = "pets-status"
  name     = "Status page"
  href     = "https://status.gravitee.io"
  location = "/md"
  order    = 3
}

```
