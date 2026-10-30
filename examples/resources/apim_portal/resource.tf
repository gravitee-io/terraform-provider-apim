resource "apim_portal" "developer-portal" {
  hrid = "developer-portal"
  name = "[Terraform] Developer Portal"
  # Entries appear in list order. Use top_navbar = [] to explicitly clear them.
  structure = {
    top_navbar = [
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
}
