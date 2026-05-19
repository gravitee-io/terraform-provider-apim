resource "apim_document" "my_document" {
  api_hrid        = "my_demo_api"
  app_hrid        = "simple_demo_app"
  content         = "Some content here"
  environment_id  = "a44e0d1b-9fa9-4d64-8b76-3634623a2e27"
  hrid            = "my_mock_document"
  name            = "My Mock Document"
  organization_id = "dedd0e0f-b3e9-4d2f-89cd-b2a9de7cb145"
}