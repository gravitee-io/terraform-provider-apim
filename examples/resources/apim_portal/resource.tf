resource "apim_portal" "my_portal" {
  environment_id = "a44e0d1b-9fa9-4d64-8b76-3634623a2e27"
  hrid           = "demo_api"
  name           = "Default Portal"
  navigation = [
    {
      display_name = "Alpha"
      order        = 1
      path         = "/projects/alpha"
    }
  ]
  organization_id = "dedd0e0f-b3e9-4d2f-89cd-b2a9de7cb145"
}