provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

resource "apim_dictionary" "test" {
  hrid     = var.hrid
  name     = "[Terraform] Manual dictionary"
  deployed = true
  type     = "MANUAL"
  manual = {
    type = "MANUAL"
    properties = {
      "env" : "test"
    }
  }
}
