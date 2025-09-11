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
  description = "Conformance tests application to be updated"
  settings = {
    app = {
      client_id = "client-id-${var.hrid}"
    }
  }
  notify_members = true
  metadata = [
    {
      key          = "usage"
      name         = "Usage"
      format       = "STRING"
      defaultValue = "Any"
    }
  ]
}
