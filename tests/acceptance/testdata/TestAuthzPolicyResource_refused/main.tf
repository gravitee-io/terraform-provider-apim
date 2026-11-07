provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

resource "apim_authz_policy" "test" {
  hrid        = var.hrid
  entity_id   = "policy.${var.hrid}"
  name        = "Acceptance ${var.hrid}"
  policy_type = "MCP"
  policy_text = "permit(principal == ?principal, action, resource);"
}
