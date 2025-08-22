resource "apim_subscription" "my_subscription" {
  api_hrid         = "demo-api"
  application_hrid = "demo-app"
  ending_at        = "2040-12-25T09:12:28Z"
  environment_id   = "...my_environment_id..."
  hrid             = "my_demo_api"
  organization_id  = "...my_organization_id..."
  plan_hrid        = "demo-plan"
}