provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

resource "apim_group" "test" {
  hrid = var.hrid
  name = "Test"
}
