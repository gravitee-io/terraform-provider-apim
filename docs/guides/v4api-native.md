---
page_title: "Kafka native"
subcategory: "V4 API"
---

# Create a v4 Native Kafka proxy API that assigns a custom attribute

The following example configures the v4 Native Kafka proxy API resource to create a Native Kafka proxyAPI with a Keyless plan.
This resource uses the Gravitee
[Assign Attributes policy](https://documentation.gravitee.io/apim/create-and-configure-apis/apply-policies/policy-reference/assign-attributes)
to assign a custom static attribute.

```HCL
resource "apim_apiv4" "kafka_native_api" {
  # should match the resource name
  hrid            = "kafka_native_api"
  name            = "[Terraform] Kafka Native proxy API"
  description     = "Connect to a local kafka cluster with a simple flow"
  version         = "1,0"
  type            = "NATIVE"
  state           = "STOPPED"
  visibility      = "PRIVATE"
  lifecycle_state = "UNPUBLISHED"
  listeners = [
    {
      kafka = {
        type = "KAFKA"
        entrypoints = [
          {
            type = "native-kafka"
          },
        ]
        host = "kafka.local"
        port = 9092
      }
    }
  ]
  endpoint_groups = [
    {
      name = "Default Native endpoint group"
      type = "native-kafka"
      endpoints = [
        {
          configuration = jsonencode({
            bootstrapServers = "kafka.local:9001"
          })
          inherit_configuration = true
          name                  = "Default Native proxy"
          secondary             = false
          type                  = "native-kafka"
          weight                = 1
        },
      ]
      shared_configuration = jsonencode({
        security = {
          protocol = "PLAINTEXT"
        }
      })
    },
  ]
  flows = [
    {
      name = "default"
      enabled : true,
      interact = [
        {
          enabled = true
          name    = "Assign custom static attribute as an example"
          policy  = "policy-assign-attributes"
          configuration = jsonencode({
            attributes = [
              {
                name  = "my.attribute"
                value = "example value"
              }
            ]
          })
        }
      ]
    },
  ]
  # known limitation, some default value is returned by default which appears to be remove during plan
  metadata=[{}]
  definition_context = {}
  plans = {
    # known limitation, key have to match name to avoid terraform plan to remain inconsistent
    KeyLess = {
      name       = "KeyLess"
      type       = "API"
      mode       = "STANDARD"
      validation = "AUTO"
      status     = "PUBLISHED"
      security = {
        type = "KEY_LESS"
      }
    }
  }
}
```