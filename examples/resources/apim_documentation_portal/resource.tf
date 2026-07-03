resource "apim_portal" "developer-portal" {
  hrid = "developer-portal"
  name = "[Terraform] Developer Portal"
  navigation = [
    {
      path         = "/docs/getting-started"
      display_name = "Getting Started"
      order        = 1
    }
  ]
}

resource "apim_documentation_portal" "getting-started" {
  portal_hrid = apim_portal.developer-portal.hrid
  hrid        = "getting-started"
  name        = "Getting Started"
  type        = "GRAVITEE_MARKDOWN"
  location    = "/docs/getting-started"
  order       = 1
  content     = <<-EOT
  # Getting Started

  Welcome to the next-gen developer portal.

  This page demonstrates multiline Markdown content managed with Terraform.
  EOT
}
