---
page_title: "API with a Group"
subcategory: "Group"
---

# Binding a Group to an API

This example demonstrates how to associate a group with a V4 API using the `groups` attribute.

By referencing `apim_group.api-developers.hrid` in the API's `groups` list,
Terraform automatically infers a dependency between the two resources.
This ensures the group is created before the API, without requiring an explicit `depends_on` block.

When the resources are destroyed, Terraform reverses the order: the API is removed first, then the group.

```terraform
resource "apim_apiv4" "with_group" {
  # should match the resource name
  hrid            = "with_group"
  name            = "[Terraform] API With a group"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
  groups = [
    apim_group.api-developers.hrid
  ]
  listeners = [
    {
      http = {
        type = "HTTP"
        entrypoints = [
          {
            type = "http-proxy"
          }
        ]
        paths = [
          {
            path = "/api-wth-group"
          }
        ]
      }
    }
  ]
  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      load_balancer = {
        type = "ROUND_ROBIN"
      }
      endpoints = [
        {
          name                  = "Default HTTP proxy"
          type                  = "http-proxy"
          weight                = 1
          inherit_configuration = false
          # Configuration is JSON as is depends on the type schema
          configuration = jsonencode({
            target = "https://api.gravitee.io/echo"
          })
        }
      ]
    }
  ]
  analytics = {
    enabled = false
  }
  plans = [
    {
      hrid        = "KeyLess"
      name        = "No security"
      type        = "API"
      mode        = "STANDARD"
      validation  = "AUTO"
      status      = "PUBLISHED"
      description = "This plan does not require any authentication"
      security = {
        type = "KEY_LESS"
      }
    }
  ]
}

resource "apim_group" "api-developers" {
  hrid = "api-developers"
  name = "Api Developers"
  members = [
    {
      source    = "memory"
      source_id = "api1"
      roles = {
        API         = "OWNER"
        APPLICATION = "USER"
        INTEGRATION = "USER"
      }
    }
  ]
}


```
