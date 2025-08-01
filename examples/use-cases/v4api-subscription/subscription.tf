resource "apim_subscription" "simple-subscription" {
  # should match the resource name
  hrid             = "simple-subscription"
  api_hrid         = apim_apiv4.simple-api-subscribed.hrid
  plan_hrid        = apim_apiv4.simple-api-subscribed.plans[0].hrid
  application_hrid = apim_application.simple-app-subscribed.hrid
  ending_at        = "2040-12-25T10:12:28+01:00"
}
