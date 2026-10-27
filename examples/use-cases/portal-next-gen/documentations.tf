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
