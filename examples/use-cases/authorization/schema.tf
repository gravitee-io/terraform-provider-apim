# One schema per PDP gateway, named by its routing scope: travel@eu.
resource "apim_authz_schema" "travel" {
  hrid   = "travel-schema"
  name   = "Travel schema"
  target = apim_authz_pdp.travel_eu.routing_scope
  schema = <<-EOT
    entity User in [Group];
    entity Group { department?: String };
    entity Hotel;
    action read appliesTo { principal: [User], resource: [Hotel] };
  EOT
}
