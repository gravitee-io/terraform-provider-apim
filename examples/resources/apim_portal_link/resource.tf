resource "apim_portal" "developer-portal" {
  hrid      = "developer-portal"
  name      = "Developer Portal"
  structure = { top_navbar = [{ path = "/resources" }] }
}

resource "apim_portal_link" "external-docs" {
  portal_hrid = apim_portal.developer-portal.hrid
  hrid        = "external-docs"
  name        = "External documentation"
  href        = "https://example.com/docs"
  location    = "/resources"
  order       = 1
  visibility  = "PUBLIC"
}
