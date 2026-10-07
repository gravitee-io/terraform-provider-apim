provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "default_member_roles" {
  type = map(string)
}

resource "apim_group" "test" {
  hrid                 = var.hrid
  name                 = var.hrid
  notify_members       = false
  default_member_roles = var.default_member_roles
}
