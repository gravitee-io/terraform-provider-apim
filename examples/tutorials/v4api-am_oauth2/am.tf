# The AM side of an API protected on the APIM gateway: the domain issuing the access tokens, the application of
# the API consumer, and the application the APIM gateway checks the tokens with.

resource "am_domain" "invoices" {
  key           = "acme-invoices"
  name          = "ACME invoices"
  description   = "Access tokens for the invoices API published on the APIM gateway"
  path          = "/acme-invoices"
  data_plane_id = "default"
}

# Signs the access tokens. Any certificate works; system = true builds it from AM's domain defaults.
resource "am_certificate" "default" {
  domain_key = am_domain.invoices.key
  key        = "default"
  system     = true
}

# The API consumer: a backend service that gets access tokens on its own behalf (client credentials).
resource "am_application" "billing" {
  domain_key  = am_domain.invoices.key
  key         = "billing"
  name        = "Billing"
  type        = "service"
  client_id   = "acme-billing"
  certificate = am_certificate.default.key
}

resource "am_client_secret" "billing" {
  domain_key      = am_domain.invoices.key
  application_key = am_application.billing.key
  key             = "billing-2026-11"
  value_wo        = var.billing_client_secret
}

# The client the APIM gateway introspects the tokens with.
resource "am_application" "apim_gateway" {
  domain_key  = am_domain.invoices.key
  key         = "apim-gateway"
  name        = "APIM gateway"
  type        = "service"
  client_id   = "acme-apim-gateway"
  certificate = am_certificate.default.key
}

resource "am_client_secret" "apim_gateway" {
  domain_key      = am_domain.invoices.key
  application_key = am_application.apim_gateway.key
  key             = "apim-gateway-2026-11"
  value_wo        = var.apim_gateway_client_secret
}

# Read from your secret store and given to the billing service too, for example with TF_VAR_billing_client_secret.
variable "billing_client_secret" {
  type      = string
  sensitive = true
  ephemeral = true
}

# Also written in the API resource of the APIM API, which is not write-only: this variable cannot be ephemeral,
# and the value is stored in the state through apim_apiv4.
variable "apim_gateway_client_secret" {
  type      = string
  sensitive = true
}
