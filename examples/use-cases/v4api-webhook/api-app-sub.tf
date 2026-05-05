resource "apim_apiv4" "webhook" {
  # should match the resource name
  hrid    = "proxy"
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

resource "apim_application" "webhook-app" {
  # should match the resource name
  hrid        = "webhook-app"
  name        = "[Terraform] Simple Application"
  description = "Simple to support Webhook subscription"
}

resource "apim_subscription" "push_subscription" {
  hrid             = "push_subscription"
  api_hrid         = apim_apiv4.webhook.hrid
  application_hrid = apim_application.webhook-app.hrid
  plan_hrid        = apim_apiv4.webhook.plans[0].hrid
  # Required for webhooks
  consumer_configuration = {
    # Matches the name of the entrypoint
    entrypoint_id = "webhook"
    # Specific to webhook entrypoint type
    entrypoint_configuration = jsonencode({
      # only this is mandatory (no underscores, this is JSON)
      callbackUrl = "https://webhook.site/bbd53b8c-e330-4881-b5ad-ddca91c52af1"
      # below is optional
      headers = [{
        name  = "X-Gravitee-Custom"
        value = "Hello"
      }],
      auth = {
        type = "basic"
        basic = {
          username = "admin"
          password = "admin"
        }
      }
      ssl = {
        // defaults
        hostnameVerifier = true
        trustAll         = false
      }
      retry = {
        // defaults
        retryOnFail         = true
        retryStrategy       = "LINEAR"
        maxAttempts         = 3
        initialDelaySeconds = 10
        maxDelaySeconds     = 30
        retryOption         = "Retry On Fail"
      }
    })
  }
}
