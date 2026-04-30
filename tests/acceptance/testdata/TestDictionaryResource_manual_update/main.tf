provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "name" {
  type = string
}
variable "deployed" {
  type = bool
}

resource "apim_dictionary" "test" {
  hrid        = var.hrid
  name        = var.name
  description = "Test example"
  deployed    = var.deployed
  type        = "MANUAL"
  manual = {
    type = "MANUAL"
    properties = {
      "env" : "test"
    }
  }
}
