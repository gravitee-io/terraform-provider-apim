# The APIM side: an API whose OAuth2 plan has the APIM gateway introspect each access token at AM, and the
# consumer's APIM application, which carries the client ID of its AM application.

locals {
  # The AM gateway URL, as the APIM gateway reaches it.
  am_gateway_url = "https://am-gateway.example.com"
}

resource "apim_apiv4" "invoices" {
  hrid            = "invoices"
  name            = "Invoices"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
  listeners = [{
    http = {
      type        = "HTTP"
      entrypoints = [{ type = "http-proxy" }]
      paths       = [{ path = "/invoices/" }]
    }
  }]
  endpoint_groups = [{
    name = "default-group"
    type = "http-proxy"
    endpoints = [{
      name          = "default"
      type          = "http-proxy"
      configuration = jsonencode({ target = "https://api.gravitee.io/echo" })
    }]
  }]

  # The API resource the OAuth2 plan introspects the tokens with: the AM domain and the client above.
  resources = [{
    name    = "acme-invoices"
    type    = "oauth2-am-resource"
    enabled = true
    configuration = jsonencode({
      serverURL      = local.am_gateway_url
      securityDomain = trimprefix(am_domain.invoices.path, "/")
      clientId       = am_application.apim_gateway.client_id
      clientSecret   = var.apim_gateway_client_secret
      version        = "V3_X"
      userClaim      = "sub"
    })
  }]

  plans = [{
    hrid       = "oauth2"
    name       = "OAuth2"
    type       = "API"
    mode       = "STANDARD"
    validation = "AUTO"
    status     = "PUBLISHED"
    security = {
      type          = "OAUTH2"
      configuration = jsonencode({ oauthResource = "acme-invoices" })
    }
  }]
}

# The consumer in APIM: an application with the client ID of its AM application, subscribed to the OAuth2 plan.
resource "apim_application" "billing" {
  hrid        = "billing"
  name        = "Billing"
  description = "The billing service, which calls the invoices API"
  settings = {
    app = {
      type      = "backend"
      client_id = am_application.billing.client_id
    }
  }
}

resource "apim_subscription" "billing" {
  hrid             = "billing-invoices"
  api_hrid         = apim_apiv4.invoices.hrid
  plan_hrid        = "oauth2"
  application_hrid = apim_application.billing.hrid
}
