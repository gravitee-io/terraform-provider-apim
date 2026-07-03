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

variable "navigation" {
  type = list(object({
    path         = string
    display_name = string
  }))
}

resource "apim_portal" "test" {
  hrid       = var.hrid
  name       = var.name
  navigation = var.navigation
}
