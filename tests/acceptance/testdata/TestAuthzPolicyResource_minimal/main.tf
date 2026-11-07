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
  policy_text = <<-EOT
    permit(
      principal in Group::"engineering",
      action == Action::"github.create_issue",
      resource == MCPTool::"github.create_issue"
    );
  EOT
}
