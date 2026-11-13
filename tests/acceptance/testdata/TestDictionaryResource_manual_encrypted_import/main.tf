provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "secret" {
  type      = string
  sensitive = true
}

import {
  to = apim_dictionary.test
  id = jsonencode({
    environment_id  = "DEFAULT"
    hrid            = var.hrid
    organization_id = "DEFAULT"
  })
}

resource "apim_dictionary" "test" {
  hrid     = var.hrid
  name     = "[Terraform] Imported encrypted dictionary"
  deployed = false
  type     = "MANUAL"
  manual = {
    properties = {
      env = "test"
    }
    encrypted_properties = {
      "api-key" = { value = var.secret }
    }
  }
}
