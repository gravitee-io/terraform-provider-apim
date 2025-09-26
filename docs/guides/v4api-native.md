---
page_title: "Kafka native"
subcategory: "V4 API"
---

# Create a v4 Native Kafka proxy API that assigns a custom attribute

The following example configures the v4 Native Kafka proxy API resource to create a Native Kafka proxyAPI with a Keyless plan.
This resource uses the Gravitee
[Assign Attributes policy](https://documentation.gravitee.io/apim/create-and-configure-apis/apply-policies/policy-reference/assign-attributes)
to assign a custom static attribute.

```terraform
resource "apim_apiv4" "kafka_native" {
  # should match the resource name
  hrid            = "kafka_native"
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
            scope = "REQUEST"
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
  plans = [
    {
      hrid       = "KeyLess"
      name       = "KeyLess"
      type       = "API"
      mode       = "STANDARD"
      validation = "AUTO"
      status     = "PUBLISHED"
      security = {
        type = "KEY_LESS"
      }
    }
  ]
}
```