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
