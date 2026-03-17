variable "organization_id" {
  type = string
}

variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "client_certificate" {
  type = string
}

resource "apim_application" "test" {
  hrid            = var.hrid
  environment_id  = var.environment_id
  organization_id = var.organization_id
  name            = "terraform example"
  description     = "Conformance tests application minimal atttributes"
  settings = {
    app = {
      client_id = "client-id-${var.hrid}"
      type      = "backend to backend"
    }
    tls = {
      client_certificate = var.client_certificate
    }
  }
}
