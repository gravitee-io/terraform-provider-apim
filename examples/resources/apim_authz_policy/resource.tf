resource "apim_authz_policy" "github_create_issue" {
  hrid        = "engineering-create-issue"
  entity_id   = "policy.engineering.create-issue"
  name        = "Engineering creates GitHub issues"
  policy_type = "MCP"
  policy_text = <<-EOT
    permit(
      principal in Group::"engineering",
      action == Action::"github.create_issue",
      resource == MCPTool::"github.create_issue"
    );
  EOT
  # An MCP proxy asks the default engine only.
  targets = ["default"]
}
