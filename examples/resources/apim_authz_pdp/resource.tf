resource "apim_authz_pdp" "orders_eu" {
  hrid          = "orders-eu"
  name          = "Orders PDP (EU)"
  target_pdp_id = "orders"
  tag           = "eu"
  endpoint = {
    path = "/authz/orders-eu"
  }
}

# Schemas, entities and policies name the PDP gateway by its routing scope,
# here orders@eu.
output "orders_eu_routing_scope" {
  value = apim_authz_pdp.orders_eu.routing_scope
}
