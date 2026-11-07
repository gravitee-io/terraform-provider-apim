resource "apim_authz_schema" "orders" {
  hrid   = "orders-schema"
  name   = "Orders schema"
  target = apim_authz_pdp.orders_eu.routing_scope
  schema = <<-EOT
    entity User in [Group];
    entity Group in [Group] { department?: String };
  EOT
}
