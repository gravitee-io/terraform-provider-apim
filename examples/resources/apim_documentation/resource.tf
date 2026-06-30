resource "apim_documentation" "my_documentation" {
  content         = "...my_content..."
  environment_id  = "a44e0d1b-9fa9-4d64-8b76-3634623a2e27"
  hrid            = "demo_api"
  location        = "/projects/alpha/docs"
  name            = "Getting Started"
  order           = 3
  organization_id = "dedd0e0f-b3e9-4d2f-89cd-b2a9de7cb145"
  portal_hrid     = "default-portal"
  type            = "ASYNCAPI"
}