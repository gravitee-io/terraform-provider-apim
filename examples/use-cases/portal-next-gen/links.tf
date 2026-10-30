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
