variable "hrid" {
  type = string
}

variable "client_certificates" {
  type = list(object({
    name      = string
    content   = string
    starts_at = string
    ends_at   = string
  }))
}

resource "apim_application" "test" {
  hrid        = var.hrid
  name        = "terraform cert rotation test"
  description = "Acceptance test for certificate rotation"
  settings = {
    app = {
      type = "test"
    }
    tls = {
      client_certificates = var.client_certificates
    }
  }
}
