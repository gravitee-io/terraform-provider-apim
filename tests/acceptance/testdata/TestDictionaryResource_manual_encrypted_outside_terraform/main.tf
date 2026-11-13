provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "properties" {
  type = map(string)
}

variable "secrets" {
  type = map(string)
}

resource "apim_dictionary" "test" {
  hrid     = var.hrid
  name     = "[Terraform] Dictionary changed outside Terraform"
  deployed = false
  type     = "MANUAL"
  manual = {
    properties           = var.properties
    encrypted_properties = { for key, secret in var.secrets : key => { value = secret } }
  }
}
