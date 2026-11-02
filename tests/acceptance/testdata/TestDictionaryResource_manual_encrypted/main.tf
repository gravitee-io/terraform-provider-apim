provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "properties" {
  type    = map(string)
  default = null
}

variable "secret" {
  type      = string
  sensitive = true
}

variable "with_secret" {
  type = bool
}

resource "apim_dictionary" "test" {
  hrid     = var.hrid
  name     = "[Terraform] Encrypted dictionary"
  deployed = false
  type     = "MANUAL"
  manual = {
    properties = var.properties
    encrypted_properties = var.with_secret ? {
      "api-key" = { value = var.secret }
    } : {}
  }
}
