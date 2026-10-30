resource "apim_api_link" "external-docs" {
  api_hrid   = "petstore-api"
  hrid       = "external-docs"
  name       = "API documentation"
  href       = "https://example.com/petstore"
  visibility = "PUBLIC"
}
