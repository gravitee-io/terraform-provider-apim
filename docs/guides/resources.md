---
page_title: "Gravitee resources"
---

# Gravitee resources specifics

There are a few things to know before working with Gravitee resources.

## Lifecycle

Each resource is identified by the unique value assigned to its `hrid` field. This field can be referenced by other resources to enable object dependency. For example, an API that invokes a Shared Policy Group can refer to a particular Shared Policy Group resource using its `hrid` value.

-> For consistency, we recommend that you use the same value for a resource's name and `hrid`. We also recommend that you choose meaningful resource names / `hrid` values, and only modify them when it is compulsory.

If the name or `hrid` of a resource is modified, Terraform does not recognize it as an existing entity. Instead, Terraform considers the resource to be a new entity. To avoid dangling objects that were created by Terraform but are no longer managed by it, Gravitee deletes a resource whose `hrid` was modified, and then creates a new resource with the modified `hrid` value.&#x20;

-> If you modify the `hrid` of an API resource, its analytics data is no longer accessible. If you then reapply the original `hrid`, analytics are accessible in the state prior to the `hrid` change.

## Plugin configuration

Gravitee includes a powerful plugin system that lets you extend its capabilities.&#x20;

You can create plugins for your APIs using custom schema, because from the perspective of the Management API, plugins are schema-less. This forces you to use the `jsonencode(any)` function repeatedly in Terraform resources. The schema-less properties that require this are all named `configuration`.

The `configuration` naming conventions used by Terraform and Gravitee differ. Terraform uses conventional `snake_case`, whereas Gravitee API uses `lowerCamelCase`. The following example illustrates this difference:

```hcl
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
```