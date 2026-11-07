# Deployed on every apply to the PDP gateways its targets name.
resource "apim_authz_policy" "agents_read_hotels" {
  hrid        = "travel-agents-read-hotels"
  entity_id   = "policy.travel.agents-read-hotels"
  name        = "Travel agents read hotels"
  policy_type = "CUSTOM"
  policy_text = <<-EOT
    permit(
      principal in Group::"agents",
      action == Action::"read",
      resource == Hotel::"grand-budapest"
    );
  EOT
  targets     = [apim_authz_pdp.travel_eu.routing_scope]
}
