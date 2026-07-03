resource "apim_portal" "developer-portal" {
  hrid = "developer-portal"
  name = "[Terraform] Developer Portal"
  navigation = [
    {
      path         = "/apis"
      display_name = "APIs"
      order        = 1
    }
  ]
}

resource "apim_apiv4" "pets" {
  // other properties ommited for simplicity
  hrid = "pets"
  # rest is ommited for simplicity
}

resource "apim_portal_listing" "public-apis" {

  portal_hrid = apim_portal.developer-portal.hrid
  hrid        = "public-apis"
  apis = [
    {
      api_hrid = apim_apiv4.pets.hrid
      location = "/apis"
      order    = 1
    }
  ]
}
