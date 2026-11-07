provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "entity_id" {
  type = string
}

variable "policy_type" {
  type = string
}

resource "apim_authz_policy" "test" {
  hrid        = var.hrid
  entity_id   = var.entity_id
  name        = "Acceptance ${var.hrid}"
  policy_type = var.policy_type
  policy_text = <<-EOT
    permit(
      principal in Group::"engineering",
      action == Action::"github.create_issue",
      resource == MCPTool::"github.create_issue"
    );
  EOT
}
