provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

resource "apim_apiv4" "test" {
  hrid    = var.hrid
  name    = "[Terraform] Simple Webhook API"
  version = "1.0"
  # API type
  type            = "MESSAGE"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
  listeners = [
    {
      # Specific listener
      subscription = {
        type = "SUBSCRIPTION"
        entrypoints = [
          {
            name = "webhook"
            type = "webhook"
          }
        ]
      }
    }
  ]
  endpoint_groups = [
    {
      name = "Mock endpoint group"
      type = "mock"
      endpoints = [
        {
          name                  = "Mock endpoint"
          type                  = "mock"
          inherit_configuration = false
          # Emits identical messages every 5 seconds, 10 times
          configuration = jsonencode({
            messageInterval = 5
            messageContent  = "Message from Gravitee mock endpoint"
            messageCount    = 10
          })
        }
      ]
    }
  ]
  flow_execution = {
    mode           = "DEFAULT"
    match_required = false
  }
  analytics = {
    enabled = false
  }
  plans = [
    {
      hrid = "push"
      name = "Push plan"
      type = "API"
      # For subscription listeners
      mode        = "PUSH"
      validation  = "AUTO"
      status      = "PUBLISHED"
      description = "This plan does not require any authentication"
    }
  ]
}
