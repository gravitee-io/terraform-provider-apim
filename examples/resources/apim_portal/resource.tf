resource "apim_portal" "developer-portal" {
  hrid = "developer-portal"
  name = "[Terraform] Developer Portal"
  # will displayed in that order
  navigation = [
    {
      path         = "/apis"
      display_name = "APIs"
    },
    {
      path         = "/docs"
      display_name = "Documentation"
    },
    {
      path         = "/docs/getting-started"
      display_name = "Getting Started"
    }
  ]
}
