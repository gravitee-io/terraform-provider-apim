---
page_title: "Subscription to an API with a JWT Plan"
subcategory: ""
---

# Subscription to an API with a JWT Plan

This example demonstrates how to create an API with a JWT plan, set up an application, and configure a subscription that grants the application access to the API.

```terraform
resource "apim_apiv4" "api_app_subscription-api" {
  # should match the resource name
  hrid            = "simple-api-subscribed"
  name            = "[Terraform] Simple PROXY API With JWT Plan"
  description     = "A simple API with JWT meant to subscribed using Terraform"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
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
            path = "/simple-api-subscribed/"
          }
        ]
        type = "HTTP"
      }
    }
  ]
  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      endpoints = [
        {
          name = "Default HTTP proxy"
          type = "http-proxy"
          configuration = jsonencode({
            target = "https://api.gravitee.io/echo"
          })
        }
      ]
      services = {}
    }
  ]
  analytics = {
    enabled = true
  }
  plans = [
    {
      hrid        = "jwt"
      name        = "Jwt Plan"
      type        = "API"
      mode        = "STANDARD"
      validation  = "AUTO"
      status      = "PUBLISHED"
      description = "A JWT plan"
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


resource "apim_application" "api_app_subscription-app" {
  # should match the resource name
  hrid        = "simple-app-subscribed"
  name        = "[Terraform] Application with client_id"
  description = "Subscription tests application"
  settings = {
    app = {
      #
      client_id = "guest"
    }
  }
}

resource "apim_subscription" "api_app_subscription" {
  # should match the resource name
  hrid             = "simple-subscription"
  api_hrid         = apim_apiv4.api_app_subscription-api.hrid
  plan_hrid        = apim_apiv4.api_app_subscription-api.plans[0].hrid
  application_hrid = apim_application.api_app_subscription-app.hrid
  ending_at        = "2042-12-25T10:12:28+02:00"
}
```
