resource "apim_application" "simple" {
  # should match the resource name
  hrid        = "simple"
  name        = "[Terraform] Simple Application"
  description = "Demonstrate applications can be created with Terraform"
  settings = {
    app = {
      client_id = "admin"
    }
  }
  metadata = [{
    name   = "hello"
    format = "STRING"
  }]
}