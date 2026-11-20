---
page_title: "Protect an API with Gravitee AM access tokens"
subcategory: "Tutorials"
---

# Protect an API with Gravitee AM access tokens

This tutorial has a backend service call an API on the APIM gateway with an access token from Gravitee Access
Management (AM). The service gets the token from the AM gateway (OAuth 2.0 client credentials), and the APIM
gateway checks it on every call. One Terraform configuration holds both sides: APIM with this provider, AM with
the `gravitee-io/am` provider.

~> All cited files in this tutorial are located in [examples/tutorials/v4api-am_oauth2](https://github.com/gravitee-io/terraform-provider-apim/tree/main/examples/tutorials/v4api-am_oauth2)

## Before you start

* The APIM provider reaches the APIM Automation API with `APIM_SERVER_URL` and `APIM_SA_TOKEN`; the AM provider
  reaches the AM Automation API with `AM_SERVER_URL` and a service-account token in `AM_SA_TOKEN`.
* Terraform 1.11 or later: the AM client secret values are write-only attributes.
* The APIM gateway reaches the AM gateway.

Both providers, in `provider.tf`:

```terraform
terraform {
  # Write-only attributes, such as the value_wo of am_client_secret, need Terraform 1.11 or later.
  required_version = ">= 1.11"
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
    am = {
      source = "gravitee-io/am"
    }
  }
}

# Set with APIM_SERVER_URL and APIM_SA_TOKEN.
provider "apim" {}

# Set with AM_SERVER_URL and AM_SA_TOKEN.
provider "am" {}

```

## The AM side

The AM domain issuing the tokens, the consumer's application and the application the APIM gateway checks the
tokens with, in `am.tf`:

```terraform
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

```

## The APIM side

The API, its OAuth2 plan and the consumer's subscription, in `apim.tf`. The API resource of type
`oauth2-am-resource` points at the AM domain, and the plan uses it to introspect each token:

```terraform
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

```

Apply with `TF_VAR_billing_client_secret` and `TF_VAR_apim_gateway_client_secret` set. The billing service then
gets a token from the AM gateway (curl prompts for its client secret) and calls the API with it:

```shell
TOKEN=$(curl -s -u acme-billing -d grant_type=client_credentials <AM gateway URL>/acme-invoices/oauth/token | jq -r .access_token)
curl -H "Authorization: Bearer $TOKEN" <APIM gateway URL>/invoices/
```

The API answers `401` without a token or with one AM did not issue.

## Check the tokens without calling AM: a JWT plan

A JWT plan checks the signature of the token with the public keys the AM domain publishes, its JWKS, instead of
calling AM for each request. Use `local.jwt_plan` in the `plans` of the API, and subscribe the application to the
`jwt` plan:

```terraform
# A JWT plan checks the token signature with the public keys the AM domain publishes (its JWKS), without a call
# to AM per request. Add it to the plans of apim_apiv4.invoices.
locals {
  jwt_plan = {
    hrid       = "jwt"
    name       = "JWT"
    type       = "API"
    mode       = "STANDARD"
    validation = "AUTO"
    status     = "PUBLISHED"
    security = {
      type = "JWT"
      configuration = jsonencode({
        signature         = "RSA_RS256"
        publicKeyResolver = "JWKS_URL"
        resolverParameter = "${local.am_gateway_url}${am_domain.invoices.path}/oidc/.well-known/jwks.json"
      })
    }
  }
}

```

## How it works

* APIM finds the subscription from the client ID in the token: the `client_id` of the APIM application is the
  one of the AM application, read from `am_application.billing`.
* `securityDomain` is the path of the AM domain without its leading slash; `version = "V3_X"` fits AM 3.x and
  later.
* The OAuth2 plan asks AM about every token: disabling the AM application (`enabled = false`) refuses its tokens
  at once. The JWT plan only checks the signature: the tokens AM already issued keep passing until they expire.
* Deleting an AM client secret does not revoke the tokens issued with it, on either plan.
* The API resource configuration is not write-only: the client secret of the APIM gateway is stored in the
  Terraform state through `apim_apiv4`. Keep the state where secrets are safe.
* To rotate the secret of the APIM gateway without an outage, apply three times: add an `am_client_secret` with a
  new key and value, then give the new value to the API resource, then delete the old secret.
