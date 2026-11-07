resource "apim_authz_entity" "travel" {
  hrid        = "travel-group"
  entity_id   = "travel"
  kind        = "PRINCIPAL"
  entity_type = "Group"
  targets     = [apim_authz_pdp.orders_eu.routing_scope]
}

resource "apim_authz_entity" "agents" {
  hrid        = "agents-group"
  entity_id   = "agents"
  kind        = "PRINCIPAL"
  entity_type = "Group"
  attributes = jsonencode({
    department = "travel"
  })
  # Members of the travel group. The platform reports parents and targets
  # sorted: another order plans no change.
  parents = [
    {
      entity_type = apim_authz_entity.travel.entity_type
      entity_id   = apim_authz_entity.travel.entity_id
    }
  ]
  targets = [apim_authz_pdp.orders_eu.routing_scope]
}
