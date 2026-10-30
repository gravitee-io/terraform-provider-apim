variable "hrid" {
  type = string
}

variable "omit_events" {
  type = bool
}

variable "empty_events" {
  type = bool
}

resource "apim_apiv4" "test" {
  hrid            = var.hrid
  name            = "Connection events defaults"
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
    }
  ]

  analytics = merge(
    { enabled = true, reporter_metrics_enabled = false },
    var.omit_events ? {} : { connection_events = var.empty_events ? [] : ["DISCONNECTED"] }
  )

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
      bootstrap_port     = 9002
      broker_range_start = 9003
      broker_range_end   = 9010
    }
  ]
}
