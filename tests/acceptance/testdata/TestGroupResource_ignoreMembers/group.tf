provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "ignore_members" {
  type = bool
}

variable "members" {
  type = list(object({
    source    = string
    source_id = string
    roles     = map(string)
  }))
  default = null
}

resource "apim_group" "test" {
  hrid           = var.hrid
  name           = var.hrid
  notify_members = false
  ignore_members = var.ignore_members
  members        = var.members
}
