---
page_title: "Application with a Group"
subcategory: "Group"
---

# Binding a Group to an Application

This example demonstrates how to associate a group with an application using the `groups` attribute.

By referencing `apim_group.developers.hrid` in the application's `groups` list,
Terraform automatically infers a dependency between the two resources.
This ensures the group is created before the application, without requiring an explicit `depends_on` block.

When the resources are destroyed, Terraform reverses the order: the application is removed first, then the group.

```terraform
resource "apim_application" "with-group" {
  # should match the resource name
  hrid        = "with-group"
  name        = "[Terraform] Application bound to a group"
  description = "Demonstrate an application with a group"
  settings = {
    app = {
      type = "test"
    }
  }
  metadata = [{
    name   = "hello"
    format = "STRING"
  }]
  groups = [
    apim_group.developers.hrid
  ]
}

resource "apim_group" "developers" {
  hrid = "app-developers"
  name = "App Developers"
  members = [
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