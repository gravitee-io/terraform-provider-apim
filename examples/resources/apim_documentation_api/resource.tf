resource "apim_documentation_api" "my_documentationapi" {
  api_hrid        = "my_demo_api"
  content         = "...my_content..."
  environment_id  = "a44e0d1b-9fa9-4d64-8b76-3634623a2e27"
  hrid            = "demo_api"
  location        = "/reference"
  name            = "API Reference"
  order           = 2
  organization_id = "dedd0e0f-b3e9-4d2f-89cd-b2a9de7cb145"
  type            = "ASYNCAPI"
}