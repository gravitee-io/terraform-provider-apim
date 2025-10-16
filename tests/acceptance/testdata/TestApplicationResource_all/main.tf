variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

variable "environment_id" {
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
  description     = "Conformance tests application with all attributes"
  domain          = "examples.com"
  notify_members  = true
  members = [
    {
      role      = "USER"
      source    = "memory"
      source_id = "api1"
    }
  ]
  background  = "https://upload.wikimedia.org/wikipedia/commons/d/df/Green_Red_Gradient_Background.png"
  picture_url = "https://upload.wikimedia.org/wikipedia/fr/0/09/Logo_App_Store_d%27Apple.png"
  settings = {
    app = {
      client_id = "client-id-${var.hrid}"
      type      = "backend to backend"
    }
    tls = {
      client_certificate = var.client_certificate
    }
  }
  metadata = [
    {
      key          = "usage"
      name         = "Usage"
      format       = "STRING"
      defaultValue = "Any"
    }
  ]
}
