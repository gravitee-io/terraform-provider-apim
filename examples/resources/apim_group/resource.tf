resource "apim_group" "my_group" {
  environment_id     = "a44e0d1b-9fa9-4d64-8b76-3634623a2e27"
  hrid               = "demo_api"
  hrid_contains_uuid = false
  members = [
    {
      roles = {
        key = "value"
      }
      source    = "gravitee"
      source_id = "john.doe@example.com"
    }
  ]
  name            = "developers"
  notify_members  = true
  organization_id = "dedd0e0f-b3e9-4d2f-89cd-b2a9de7cb145"
}