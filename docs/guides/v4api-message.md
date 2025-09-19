---
page_title: "Kafka to Websocket"
subcategory: "V4 API"
---

# Create a v4 message API that fetches Kafka messages

The following example configures the v4 message API resource to create a v4 message API with a Keyless plan.
This API has a WebSocket entrypoint and Kafka endpoint.
It fetches messages from a Kafka cluster and publishes them to a client's WebSocket connection.

```HCL
resource "apim_apiv4" "message" {
  # should match the resource name
  hrid            = "message"
  name            = "[Terraform] Websocket to Kafka message API"
  description     = "Message API that publishes message fetch a Kafla cluster to a websocket."
  version         = "1,0"
  type            = "MESSAGE"
  state           = "STOPPED"
  visibility      = "PRIVATE"
  lifecycle_state = "UNPUBLISHED"
  listeners = [
    {
      http = {
        type = "HTTP"
        entrypoints = [
          {
            type = "websocket"
            configuration = jsonencode({
              publisher = {
                enabled = true
              }
              subscriber = {
                enabled = true
              }
            })
          }
        ]
        paths = [
          {
            path = "/message/"
          }
        ]
      }
    }
  ]
  endpoint_groups = [
    {
      name = "Default Kafka group"
      type = "kafka"
      endpoints = [
        {
          name = "Default Kafka"
          type = "kafka"
          configuration = jsonencode({
            bootstrapServers = "localhost:8082"
          })
          inherit_configuration = true
        }
      ]
      shared_configuration = jsonencode({
        consumer = {
          enabled               = true
          autoOffsetReset       = "latest"
          checkTopicExistence   = false
          encodeMessageId       = true
          removeConfluentHeader = false
          topics = [
            "test"
          ]
        }
        security = {
          protocol = "PLAINTEXT"
        }
      })
    }
  ]
  flow_execution = {
    match_required = false
    mode           = "DEFAULT"
  }
  flows = []
  analytics = {
    enabled = true
    sampling = {
      type  = "COUNT"
      value = 10
    }
  }
  plans = [
    {
      hrid        = "key-less"
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
```