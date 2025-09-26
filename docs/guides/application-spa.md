---
page_title: "Simple page Application"
subcategory: "Application"
---

# Single page application

The following example creates a Single Page Application with OAuth.

```terraform
resource "apim_application" "spa" {
  # should match the resource name
  hrid        = "spa"
  name        = "[Terraform] Application for Single Page Application with OAuth"
  description = "Demonstrates applications for Oauth can be created with Terraform"
  domain      = "https://example.com"
  settings = {
    oauth = {
      application_type = "browser"
      redirect_uris = [
        "https://example.com/callback"
      ]
      grant_types = [
        "authorization_code"
      ]
    }
  }
}
```