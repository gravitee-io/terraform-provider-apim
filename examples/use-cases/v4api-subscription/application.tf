
resource "apim_application" "simple-app-subscribed" {
  # should match the resource name
  hrid        = "simple-app-subscribed"
  name        = "[Terraform] Application with client_id"
  description = "Subscription tests application"
  settings = {
    app = {
      client_id = "foo"
    }
  }
}