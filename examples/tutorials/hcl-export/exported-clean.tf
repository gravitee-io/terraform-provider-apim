# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "{\n  \"organization_id\": \"DEFAULT\",\n  \"environment_id\": \"DEFAULT\",\n  \"hrid\": \"2f0d6943-1d80-4cfa-8d69-431d807cface\"\n}\n"
resource "apim_apiv4" "export" {
  analytics = {
    enabled = true
  }
  endpoint_groups = [
    {
      endpoints = [
        {
          configuration = jsonencode({
            "target" = "https://api.gravitee.io/echo"
          })
          inherit_configuration = true
          name                  = "Default HTTP proxy"
          secondary             = false
          services = {
          }
          type   = "http-proxy"
          weight = 1
        },
      ]
      load_balancer = {
        type = "ROUND_ROBIN"
      }
      name = "Default HTTP proxy group"
      services = {
      }
      shared_configuration = jsonencode({
        "headers" = [
          {
            "name"  = "Foo",
            "value" = "Bar"
          }
        ],
        "http" = {
          "connectTimeout"           = 3000,
          "followRedirects"          = false,
          "idleTimeout"              = 0,
          "keepAlive"                = true,
          "keepAliveTimeout"         = 30000,
          "maxConcurrentConnections" = 20,
          "pipelining"               = false,
          "propagateClientHost"      = true,
          "readTimeout"              = 10000,
          "useCompression"           = true,
          "version"                  = "HTTP_1_1"
        },
        "proxy" = {
          "enabled"        = false,
          "useSystemProxy" = false
        },
        "ssl" = {
          "hostnameVerifier" = true,
          "keyStore" = {
            "type" = ""
          },
          "trustAll" = false,
          "trustStore" = {
            "type" = ""
          }
        }
      })
      type = "http-proxy"
    },
  ]
  environment_id = "DEFAULT"
  flow_execution = {
    match_required = false
    mode           = "DEFAULT"
  }
  flows = [
    {
      enabled = true
      name    = "Default"
      request = [
        {
          configuration = jsonencode({
            "hrid" = "common-oas-policy"
          })
          enabled = true
          name    = "Common OAS policy"
          policy  = "shared-policy-group-policy"
        },
      ]
      response = [
        {
          configuration = jsonencode({
            "addHeaders" = [
              {
                "name"  = "Warning",
                "value" = "None"
              }
            ],
            "scope" = "REQUEST"
          })
          enabled = true
          name    = "Transform Headers"
          policy  = "transform-headers"
        },
      ]
      selectors = [
        {
          http = {
            path          = "/"
            path_operator = "EQUALS"
            type          = "HTTP"
          }
        },
      ]
    },
  ]
  hrid            = "automation-exportable"
  lifecycle_state = "UNPUBLISHED"
  listeners = [
    {
      http = {
        entrypoints = [
          {
            configuration = "{}"
            qos           = "AUTO"
            type          = "http-proxy"
          },
        ]
        paths = [
          {
            override_access = false
            path            = "/terraform/exported-example/"
          },
        ]
        type = "HTTP"
      }
    },
  ]
  metadata = [
    {
      default_value = "support@change.me"
      format        = "MAIL"
      hidden        = false
      key           = "email-support"
      name          = "email-support"
      value         = "$${(api.primaryOwner.email)!''}"
    },
  ]
  name            = "Automation exportable"
  organization_id = "DEFAULT"
  pages = [
    {
      homepage   = false
      hrid       = "legal"
      name       = "Legal"
      published  = true
      type       = "FOLDER"
      visibility = "PUBLIC"
    },
    {
      content     = "Anyone is free to use this example API"
      homepage    = false
      hrid        = "general-conditions"
      name        = "General Conditions"
      parent_hrid = "legal"
      published   = true
      type        = "MARKDOWN"
      visibility  = "PUBLIC"
    },
    {
      homepage   = false
      hrid       = "aside"
      name       = "Aside"
      published  = true
      type       = "SYSTEM_FOLDER"
      visibility = "PUBLIC"
    },
  ]
  plans = [
    {
      description             = "Default unsecured plan"
      general_conditions_hrid = "general-conditions"
      hrid                    = "default-keyless-unsecured"
      mode                    = "STANDARD"
      name                    = "Default Keyless (UNSECURED)"
      security = {
        configuration = "{}"
        type          = "KEY_LESS"
      }
      status     = "PUBLISHED"
      type       = "API"
      validation = "MANUAL"
    },
  ]
  state      = "STOPPED"
  type       = "PROXY"
  version    = "1"
  visibility = "PRIVATE"
}