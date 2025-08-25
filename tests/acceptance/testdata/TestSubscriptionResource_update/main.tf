provider "apim" {
  organization_id = "DEFAULT"
}

variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

variable "ending_at" {
  type = string
}

resource "apim_apiv4" "test" {
  environment_id  = var.environment_id
  hrid            = "api-${var.hrid}"
  lifecycle_state = "UNPUBLISHED"
  name            = "terraform_example"
  organization_id = var.organization_id
  state           = "STOPPED"
  type            = "PROXY"
  version         = "1"
  visibility      = "PRIVATE"
  analytics = {
    enabled = true
  }
  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      endpoints = [
        {
          name = "Default HTTP proxy"
          type = "http-proxy"
          configuration = jsonencode({
            target = "https://example.com"
          })
        }
      ]
      services = {}
    }
  ]
  listeners = [
    {
      http = {
        entrypoints = [
          {
            type = "http-proxy"
          }
        ]
        paths = [
          {
            path           = "/api/${var.hrid}/"
            overrideAccess = false
          }
        ]
        type = "HTTP"
      }
    }
  ]
  plans = [
    {
      hrid       = "jwt"
      name       = "Jwt Plan"
      type       = "API"
      mode       = "STANDARD"
      validation = "AUTO"
      status     = "PUBLISHED"
      description : "A JWT plan"
      security = {
        type = "JWT"
        configuration = jsonencode({
          signature         = "RSA_RS256"
          publicKeyResolver = "GIVEN_KEY"
          resolverParameter = <<-EOT
          -----BEGIN PUBLIC KEY-----
          MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAv8YGSPoQEl7lXnp8OHkb
          AOPYZ81rzXkmO83d0P8G78qWzi3gPnODm6Qxi2NbgcWXqQlZXxPkDTS3Xck1V3WY
          E9voqQE7UEwpFBolqtUHQqL4w2vr/eUtZv9t3DdtoCcIj4xLmJUw7PS7jAb9quq0
          XiVN692d6LI62T+9LyN+kcWHTpUyMBB8oxfQ9ekkGHskTc6LgYovKK+9lKoJv6gg
          0ge8YAFbpjJBZbVX3jV8qeszgw9Xdhs3w/S8QnvWa3Cv0+c47oxZjXwpAa8ARzfn
          D/5oK4CWRRy+t3QUndSR0cBR+bU0YFks3mmbl514/ywOXRf/sZmXaJkNejfNHQVa
          hJgj/Z3W3F8GKksuFF14+BK2KX30bsQL3e4SeN0Wv6DF1UloG0T396yDd/o7L3ZC
          DBlRB44OZ8sO3h8iSW7wVX0sGj/OKc4smo5dgP0r4+Fm2EVmVFU5YvEkFcy0Xoth
          QmLwq0lJc7BdRMpAfRZLbW5WSlb2jgvxA/VI/ScLTRWZI7DGbzHRBS6J8Rnt3Inq
          jo7mUV1juBs3RhpxdOmg1LpGLAtQdcSSnX3IyyEVbzTVb22Px0EGAlKzMs6bnTJf
          3TbZd/C0iqd6QOyaTh7D4Nr7ClfWAaYGZBA/FsHWA88fOsIQCtovWjp9A8i1+VQ5
          HEy1rpaHPGHt1DFt2hu+d3MCAwEAAQ==
          -----END PUBLIC KEY-----
          EOT
          userClaim         = "sub"
          clientIdClaim     = "client_id"
        })
      }
    }
  ]
}

resource "apim_application" "test" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "app-${var.hrid}"
  name            = "terraform example"
  description     = "Subscription tests application"
  settings = {
    app = {
      client_id = "foo-${var.hrid}"
    }
  }
}

resource "apim_subscription" "test" {
  environment_id   = var.environment_id
  organization_id  = var.organization_id
  hrid             = var.hrid
  api_hrid         = apim_apiv4.test.hrid
  plan_hrid        = apim_apiv4.test.plans[0].hrid
  application_hrid = apim_application.test.hrid
  ending_at        = var.ending_at
}
