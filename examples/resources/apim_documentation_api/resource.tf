resource "apim_apiv4" "pets-api" {
  hrid = "pets-api"
  // other properties ommited for simplicity
}

resource "apim_documentation_api" "api-docs" {
  api_hrid = apim_apiv4.pets-api.hrid
  hrid     = "api-docs"
  name     = "API Documentation"
  type     = "GRAVITEE_MARKDOWN"
  content  = <<-EOT
  # Pets API

  This documentation page is attached directly to the API via `apim_documentation_api`.
  EOT
}
