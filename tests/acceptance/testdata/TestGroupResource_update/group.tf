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

variable "members" {
  type = list(
    object({
      source    = string
      source_id = string
      roles     = map(string)
    })
  )
}

resource "apim_group" "test" {
  hrid    = var.hrid
  name    = var.name
  members = var.members
}
