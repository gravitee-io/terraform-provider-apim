variable "hrid" {
  type = string
}

variable "events" {
  type = list(string)
}

variable "notification_groups" {
  type = list(string)
}

provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

resource "apim_group" "notif" {
  hrid = "${var.hrid}-notif-group"
  name = "Notification Group"
}

resource "apim_apiv4" "test" {
  depends_on      = [apim_group.notif]
  hrid            = var.hrid
  name            = "Test Console Notification"
  version         = "1"
  type            = "PROXY"
  state           = "STOPPED"
  visibility      = "PRIVATE"
  lifecycle_state = "UNPUBLISHED"

  groups = [
    apim_group.notif.hrid
  ]

  console_notification = {
    groups = var.notification_groups
    events = var.events
  }

  listeners = [
    {
      http = {
        entrypoints = [
          {
            type = "http-proxy"
          }
        ]
        paths = [
          {
            path = "/${var.hrid}/"
          }
        ]
        type = "HTTP"
      }
    }
  ]

  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      load_balancer = {
        type = "ROUND_ROBIN"
      }
      endpoints = [
        {
          name = "default"
          type = "http-proxy"
          configuration = jsonencode({
            target = "https://example.com"
          })
        }
      ]
    }
  ]

  plans = [
    {
      hrid        = "Keyless"
      description = "No sec"
      mode        = "STANDARD"
      name        = "No security"
      status      = "PUBLISHED"
      type        = "API"
      validation  = "AUTO"
      security = {
        type = "KEY_LESS"
      }
    }
  ]
}
