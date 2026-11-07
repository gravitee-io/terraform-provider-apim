resource "apim_authz_entity" "agents" {
  hrid        = "travel-agents"
  entity_id   = "agents"
  kind        = "PRINCIPAL"
  entity_type = "Group"
  attributes = jsonencode({
    department = "travel"
  })
  targets = [apim_authz_pdp.travel_eu.routing_scope]
}

resource "apim_authz_entity" "alice" {
  hrid        = "travel-alice"
  entity_id   = "alice"
  kind        = "PRINCIPAL"
  entity_type = "User"
  parents = [
    {
      entity_type = apim_authz_entity.agents.entity_type
      entity_id   = apim_authz_entity.agents.entity_id
    }
  ]
  targets = [apim_authz_pdp.travel_eu.routing_scope]
}

resource "apim_authz_entity" "grand_budapest" {
  hrid        = "travel-grand-budapest"
  entity_id   = "grand-budapest"
  kind        = "RESOURCE"
  entity_type = "Hotel"
  targets     = [apim_authz_pdp.travel_eu.routing_scope]
}
