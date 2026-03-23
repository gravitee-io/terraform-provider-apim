provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

resource "apim_application" "test" {
  hrid        = var.hrid
  name        = "terraform example"
  description = "Test both app and oauth settings set"
  settings = {
    app = {
      client_id = "client-id"
    }
    oauth = {
      application_type = "web"
    }
  }
}
