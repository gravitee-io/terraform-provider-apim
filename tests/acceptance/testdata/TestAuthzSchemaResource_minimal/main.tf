provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

resource "apim_authz_schema" "test" {
  hrid   = var.hrid
  name   = "Acceptance ${var.hrid}"
  schema = <<-EOT
    entity User in [Group];
    entity Group;
  EOT
}
