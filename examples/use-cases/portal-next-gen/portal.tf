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
