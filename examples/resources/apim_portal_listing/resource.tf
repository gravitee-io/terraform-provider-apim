resource "apim_portal_listing" "my_portallisting" {
  apis = [
    {
      api_hrid = "pets-api"
      location = "/projects/alpha"
      order    = 1
    }
  ]
  environment_id  = "a44e0d1b-9fa9-4d64-8b76-3634623a2e27"
  hrid            = "demo_api"
  organization_id = "dedd0e0f-b3e9-4d2f-89cd-b2a9de7cb145"
  portal_hrid     = "default-portal"
}