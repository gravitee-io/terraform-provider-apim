---
page_title: "API with Console Notifications"
subcategory: "V4 API"
---

# Setting up Console Notifications on an API

This example demonstrates how to configure console notifications on a V4 API using the `console_notification` attribute.

## Notification groups must belong to the API

The groups listed in `console_notification.groups` must be a subset of the groups assigned to the API via its `groups` attribute.
This is because APIM uses the API's group membership to resolve which users should receive notifications.

In the example below, `apim_group.api-developers` is referenced in both `groups` and `console_notification.groups`,
ensuring that members of this group will receive notifications when the API is started or stopped.

## Terraform dependency management

Because `apim_apiv4.with_notif` references `apim_group.api-developers.hrid` in both `groups` and `console_notification.groups`,
Terraform automatically infers a dependency between the two resources.
This guarantees the group is fully created before the API, without requiring an explicit `depends_on` block.

```terraform
resource "apim_apiv4" "with_notif" {
  # should match the resource name
  hrid            = "with_notif"
  name            = "[Terraform] API With a console notifications"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
  groups = [
    apim_group.api-developers.hrid
  ]
  console_notification = {
    groups = [
      apim_group.api-developers.hrid
    ]
    events = [
      "API_STARTED",
      "API_STOPPED"
    ]
  }
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
            path = "/api-with-notifications"
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
