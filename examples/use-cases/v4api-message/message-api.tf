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