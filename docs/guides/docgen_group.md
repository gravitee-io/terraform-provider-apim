---
page_title: "Simple Group"
subcategory: "Group"
---

# Simple Group

This example demonstrates how to use manage groups of user.
It allows assigning roles (pre-defined or custom) to members (allowing them can do) on given scopes (API, APPLICATION...).

In this example, `api1` an "In Memory" user, has the `OWNER` role to manage API but only a simple `USER` role for Applications and Integration.

```terraform
resource "apim_group" "developers" {
  hrid = "developers"
  name = "Developers"
  members = [
    {
      source    = "memory"
      source_id = "api1"
      roles = {
        API         = "OWNER"
        APPLICATION = "USER"
        INTEGRATION = "USER"
      }
    },
    {
      source    = "memory"
      source_id = "application1"
      roles = {
        API         = "USER"
        APPLICATION = "OWNER"
        INTEGRATION = "USER"
      }
    }
  ]
}

```
