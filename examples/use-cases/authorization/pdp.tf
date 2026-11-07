# The engine of target_pdp_id "travel", run on the gateways tagged "eu", with its
# AuthZEN endpoint. target_pdp_id and tag are fixed at creation.
resource "apim_authz_pdp" "travel_eu" {
  hrid          = "travel-eu"
  name          = "Travel PDP (EU)"
  target_pdp_id = "travel"
  tag           = "eu"
  endpoint = {
    path = "/authz/travel-eu"
  }
}
