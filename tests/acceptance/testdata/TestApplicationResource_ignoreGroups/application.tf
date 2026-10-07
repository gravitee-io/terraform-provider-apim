provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "ignore_groups" {
  type = bool
}

variable "declare_group" {
  type = bool
}

resource "apim_group" "test" {
  hrid = "${var.hrid}-group"
  name = "${var.hrid}-group"
}

resource "apim_application" "test" {
  hrid          = var.hrid
  name          = var.hrid
  description   = "ignore_groups acceptance test"
  ignore_groups = var.ignore_groups
  groups        = var.declare_group ? [apim_group.test.hrid] : null
}
