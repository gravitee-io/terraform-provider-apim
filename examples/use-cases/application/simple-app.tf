resource "apim_application" "simple-app" {
  hrid        = "simple-app"
  name        = "[Terraform] Simple Application"
  description = "Demonstrate applications can be created with Terraform"
  settings = {
    app = {
      client_id = "admin"
    }
  }
}