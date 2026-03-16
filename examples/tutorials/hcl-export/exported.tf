# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "{\n  \"organization_id\": \"DEFAULT\",\n  \"environment_id\": \"DEFAULT\",\n  \"hrid\": \"2f0d6943-1d80-4cfa-8d69-431d807cface\"\n}\n"
resource "apim_apiv4" "export" {
  analytics = {
    enabled  = true
    logging  = null
    sampling = null
    tracing  = null
  }
  categories  = []
  description = null
  endpoint_groups = [
    {
      endpoints = [
        {
          configuration         = "{\"target\":\"https://api.gravitee.io/echo\"}"
          inherit_configuration = true
          name                  = "Default HTTP proxy"
          secondary             = false
          services = {
            health_check = null
          }
          shared_configuration_override = null
          tenants                       = []
          type                          = "http-proxy"
          weight                        = 1
        },
      ]
      load_balancer = {
        type = "ROUND_ROBIN"
      }
      name = "Default HTTP proxy group"
      services = {
        discovery    = null
        health_check = null
      }
      shared_configuration = "{\"headers\":[{\"name\":\"Foo\",\"value\":\"Bar\"}],\"http\":{\"connectTimeout\":3000,\"followRedirects\":false,\"idleTimeout\":0,\"keepAlive\":true,\"keepAliveTimeout\":30000,\"maxConcurrentConnections\":20,\"pipelining\":false,\"propagateClientHost\":true,\"readTimeout\":10000,\"useCompression\":true,\"version\":\"HTTP_1_1\"},\"proxy\":{\"enabled\":false,\"useSystemProxy\":false},\"ssl\":{\"hostnameVerifier\":true,\"keyStore\":{\"type\":\"\"},\"trustAll\":false,\"trustStore\":{\"type\":\"\"}}}"
      type                 = "http-proxy"
    },
  ]
  environment_id = "DEFAULT"
  failover       = null
  flow_execution = {
    match_required = false
    mode           = "DEFAULT"
  }
  flows = [
    {
      enabled = true
      entrypoint_connect = [
      ]
      interact = [
      ]
      name = "Default"
      publish = [
      ]
      request = [
        {
          condition         = null
          configuration     = "{\"hrid\":\"common-oas-policy\"}"
          description       = null
          enabled           = true
          message_condition = null
          name              = "Common OAS policy"
          policy            = "shared-policy-group-policy"
        },
      ]
      response = [
        {
          condition         = null
          configuration     = "{\"addHeaders\":[{\"name\":\"Warning\",\"value\":\"None\"}],\"removeHeaders\":[],\"scope\":\"REQUEST\",\"whitelistHeaders\":[]}"
          description       = null
          enabled           = true
          message_condition = null
          name              = "Transform Headers"
          policy            = "transform-headers"
        },
      ]
      selectors = [
        {
          channel   = null
          condition = null
          http = {
            methods       = []
            path          = "/"
            path_operator = "EQUALS"
            type          = "HTTP"
          }
          mcp = null
        },
      ]
      subscribe = [
      ]
      tags = []
    },
  ]
  groups          = []
  hrid            = "automation-exportable"
  labels          = []
  lifecycle_state = "UNPUBLISHED"
  listeners = [
    {
      http = {
        cors = null
        entrypoints = [
          {
            configuration = "{}"
            dlq           = null
            qos           = "AUTO"
            type          = "http-proxy"
          },
        ]
        paths = [
          {
            host            = null
            override_access = false
            path            = "/terraform/exported-example/"
          },
        ]
        servers = []
        type    = "HTTP"
      }
      kafka        = null
      subscription = null
      tcp          = null
    },
  ]
  members = [
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
  notify_members  = null
  organization_id = "DEFAULT"
  pages = [
    {
      configuration = null
      content       = null
      homepage      = false
      hrid          = "legal"
      name          = "Legal"
      parent_hrid   = null
      published     = true
      source        = null
      type          = "FOLDER"
      visibility    = "PUBLIC"
    },
    {
      configuration = null
      content       = "Anyone is free to use this example API"
      homepage      = false
      hrid          = "general-conditions"
      name          = "General Conditions"
      parent_hrid   = "legal"
      published     = true
      source        = null
      type          = "MARKDOWN"
      visibility    = "PUBLIC"
    },
    {
      configuration = null
      content       = null
      homepage      = false
      hrid          = "aside"
      name          = "Aside"
      parent_hrid   = null
      published     = true
      source        = null
      type          = "SYSTEM_FOLDER"
      visibility    = "PUBLIC"
    },
  ]
  plans = [
    {
      characteristics = []
      description     = "Default unsecured plan"
      excluded_groups = []
      flows = [
      ]
      general_conditions_hrid = "general-conditions"
      hrid                    = "default-keyless-unsecured"
      mode                    = "STANDARD"
      name                    = "Default Keyless (UNSECURED)"
      security = {
        configuration = "{}"
        type          = "KEY_LESS"
      }
      selection_rule = null
      status         = "PUBLISHED"
      tags           = []
      type           = "API"
      validation     = "MANUAL"
    },
  ]
  properties = [
  ]
  resources = [
  ]
  response_templates = null
  services           = null
  state              = "STOPPED"
  tags               = []
  type               = "PROXY"
  version            = "1"
  visibility         = "PRIVATE"
}
