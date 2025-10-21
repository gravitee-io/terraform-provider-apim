---
page_title: "Gravitee resources"
---

# Gravitee resources specifics

Gravitee APIM Terraform provider brings a few novelties on how Gravitee resources are managed.

## Human-Readable Identifier (`hrid`) property

### Overview

Each Gravitee resource is identified by the unique value using `hrid` field.
This field can be used by other resources to manage dependencies between resources.

For example, an API that invokes a Shared Policy Group can refer to a particular one using its `hrid` value.
For consistency, we recommend that you use the same value for Terraform resource name and `hrid`.

We also recommend that you choose meaningful resource names / `hrid` values,
and only modify them when it is compulsory.

### Lifecycle

APIM provider will see a Terraform resource name or `hrid` update equally:
both changes will be considered as a resource swap. Terraform plan will prepare one destroy operation and one create operation.

~> If you modify the `hrid` of a resource (API and Application), analytics data are no longer accessible.
If you then reapply the original `hrid`, analytics are accessible again, but all analytics data collected in between will remain invisible.

That's the reason we encourage you to make Terraform resource name equal Gravitee `hrid` and change it only if it makes sense.

### Unicity

The following resources must be uniquely identified for a given environment:
* V4 API (proxy, message, Kafka native) `apim_v4api`
* Shared Policy Group `apim_shared_policy_group`
* Applications `apim_application`

-> Thus `hrid` must be unique for the environment.
They can have the same `hrid` in other environments; they will remain unique within an organization.

The following resources must be uniquely identified for a given `apim_v4api`:
* Subscriptions `apim_subscription`
* Nested objects of `apim_v4api` containing `hrid` (those do not trigger resource re-creation):
    * Pages
    * Plans

-> `hrid` must unique for the API.
Two `apim_subscription` with `hrid` set to "`foo`" can coexist as long as the attribute `api_hrid` is different.
Note that Terraform resource names need to be different.
We suggest that subscription have different `hrid` for different APIs to keep resource name and `hrid` in sync

## Plugin configuration

Gravitee includes a powerful plugin system that lets you extend its capabilities.
You can create plugins for your APIs using custom schema,
that is why from the perspective of the Management API, plugins are schema-less.
Thus, APIs and Shared Policy Groups have many properties where Terraform type is `any`.

This forces the user to use the `jsonencode(any)` function repeatedly in Terraform resources.
These schemaless properties that require this are all named `configuration`.

The configuration naming conventions used by Terraform and Gravitee differ.
Terraform uses conventional `snake_case`, whereas Gravitee API uses `camelCase`.

The following example (extract) illustrates this difference:
```hcl
  #...
  endpoints = [
    {
      name = "Default Kafka"
      type = "kafka"
      # configuration => Gravitee plugin
      configuration = jsonencode({
        # camelCase
        bootstrapServers = "localhost:8082"
      })
      # snake_case
      inherit_configuration = true
    }
  ]
  #...
```