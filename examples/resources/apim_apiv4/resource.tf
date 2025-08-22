resource "apim_apiv4" "my_apiv4" {
  analytics = {
    enabled = false
    logging = {
      condition = "...my_condition..."
      content = {
        headers          = false
        message_headers  = true
        message_metadata = false
        message_payload  = true
        payload          = true
      }
      message_condition = "...my_message_condition..."
      mode = {
        endpoint   = false
        entrypoint = true
      }
      phase = {
        request  = true
        response = false
      }
    }
    sampling = {
      type  = "TEMPORAL"
      value = "...my_value..."
    }
    tracing = {
      enabled = true
      verbose = true
    }
  }
  categories = [
    "health",
    "media",
  ]
  description = "I can use many characters to describe this API."
  endpoint_groups = [
    {
      endpoints = [
        {
          configuration         = "{ \"see\": \"documentation\" }"
          inherit_configuration = false
          name                  = "default-endpoint"
          secondary             = false
          services = {
            health_check = {
              configuration          = "{ \"see\": \"documentation\" }"
              enabled                = false
              override_configuration = true
              type                   = "...my_type..."
            }
          }
          shared_configuration_override = "{ \"see\": \"documentation\" }"
          tenants = [
            "..."
          ]
          type   = "mock"
          weight = 9
        }
      ]
      load_balancer = {
        type = "ROUND_ROBIN"
      }
      name = "default-endpoint-group"
      services = {
        discovery = {
          configuration          = "{ \"see\": \"documentation\" }"
          enabled                = true
          override_configuration = true
          type                   = "...my_type..."
        }
        health_check = {
          configuration          = "{ \"see\": \"documentation\" }"
          enabled                = true
          override_configuration = true
          type                   = "...my_type..."
        }
      }
      shared_configuration = "{ \"see\": \"documentation\" }"
      type                 = "default"
    }
  ]
  environment_id = "...my_environment_id..."
  failover = {
    enabled             = true
    max_failures        = 6
    max_retries         = 8
    open_state_duration = 510
    per_subscription    = true
    slow_call_duration  = 54
  }
  flow_execution = {
    match_required = true
    mode           = "BEST_MATCH"
  }
  flows = [
    {
      connect = [
        {
          condition         = "...my_condition..."
          configuration     = "{ \"see\": \"documentation\" }"
          description       = "...my_description..."
          enabled           = false
          message_condition = "...my_message_condition..."
          name              = "...my_name..."
          policy            = "...my_policy..."
        }
      ]
      enabled = true
      id      = "4e6abbd2-c0c6-462d-be9e-6371209af34b"
      interact = [
        {
          condition         = "...my_condition..."
          configuration     = "{ \"see\": \"documentation\" }"
          description       = "...my_description..."
          enabled           = true
          message_condition = "...my_message_condition..."
          name              = "...my_name..."
          policy            = "...my_policy..."
        }
      ]
      name = "My Flow"
      publish = [
        {
          condition         = "...my_condition..."
          configuration     = "{ \"see\": \"documentation\" }"
          description       = "...my_description..."
          enabled           = false
          message_condition = "...my_message_condition..."
          name              = "...my_name..."
          policy            = "...my_policy..."
        }
      ]
      request = [
        {
          condition         = "...my_condition..."
          configuration     = "{ \"see\": \"documentation\" }"
          description       = "...my_description..."
          enabled           = true
          message_condition = "...my_message_condition..."
          name              = "...my_name..."
          policy            = "...my_policy..."
        }
      ]
      response = [
        {
          condition         = "...my_condition..."
          configuration     = "{ \"see\": \"documentation\" }"
          description       = "...my_description..."
          enabled           = true
          message_condition = "...my_message_condition..."
          name              = "...my_name..."
          policy            = "...my_policy..."
        }
      ]
      selectors = [
        {
          channel = {
            channel          = "/my/channel"
            channel_operator = "EQUALS"
            entrypoints = [
              "..."
            ]
            operations = [
              "SUBSCRIBE"
            ]
            type = "HTTP"
          }
        }
      ]
      subscribe = [
        {
          condition         = "...my_condition..."
          configuration     = "{ \"see\": \"documentation\" }"
          description       = "...my_description..."
          enabled           = true
          message_condition = "...my_message_condition..."
          name              = "...my_name..."
          policy            = "...my_policy..."
        }
      ]
      tags = [
        "tag1",
        "tag2",
      ]
    }
  ]
  groups = [
    "..."
  ]
  hrid = "my_demo_api"
  labels = [
    "..."
  ]
  lifecycle_state = "CREATED"
  listeners = [
    {
      http = {
        cors = {
          allow_credentials = false
          allow_headers = [
            "..."
          ]
          allow_methods = [
            "..."
          ]
          allow_origin = [
            "..."
          ]
          enabled = false
          expose_headers = [
            "..."
          ]
          max_age      = 3
          run_policies = false
        }
        entrypoints = [
          {
            configuration = "{ \"see\": \"documentation\" }"
            dlq = {
              endpoint = "...my_endpoint..."
            }
            qos  = "NONE"
            type = "http-get"
          }
        ]
        path_mappings = [
          "..."
        ]
        paths = [
          {
            host            = "...my_host..."
            override_access = false
            path            = "...my_path..."
          }
        ]
        servers = [
          "..."
        ]
        type = "HTTP"
      }
    }
  ]
  members = [
    {
      role      = "REVIEWER"
      source    = "system"
      source_id = "john.doe@example.com"
    }
  ]
  metadata = [
    {
      default_value = "...my_default_value..."
      format        = "STRING"
      key           = "...my_key..."
      name          = "...my_name..."
      value         = "...my_value..."
    }
  ]
  name            = "My Api"
  notify_members  = false
  organization_id = "...my_organization_id..."
  pages = [
    {
      access_controls = [
        {
          reference_id   = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
          reference_type = "GROUP"
        }
      ]
      attached_media = [
        {
          attached_at = "2018-01-01T00:00:00Z"
          hash        = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
          name        = "My Media"
        }
      ]
      configuration = {
        key = "value"
      }
      content = "My Page content"
      content_revision = {
        id       = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
        revision = 1
      }
      content_type             = "application/json"
      excluded_access_controls = true
      general_conditions       = true
      hidden                   = true
      homepage                 = true
      hrid                     = "my_demo_api"
      last_contributor         = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
      metadata = {
        key = "value"
      }
      name        = "My Page"
      order       = 1
      parent_id   = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
      parent_path = "/parent"
      published   = true
      source = {
        configuration = "{ \"see\": \"documentation\" }"
        type          = "http-fetcher"
      }
      translations = [
        {
          access_controls = [
            {
              reference_id   = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
              reference_type = "GROUP"
            }
          ]
          attached_media = [
            {
              attached_at = "2018-01-01T00:00:00Z"
              hash        = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
              name        = "My Media"
            }
          ]
          configuration = {
            key = "value"
          }
          content = "My Page content"
          content_revision = {
            id       = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
            revision = 1
          }
          content_type             = "application/json"
          excluded_access_controls = false
          general_conditions       = false
          hidden                   = true
          homepage                 = true
          hrid                     = "my_demo_api"
          last_contributor         = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
          metadata = {
            key = "value"
          }
          name        = "My Page"
          order       = 1
          parent_id   = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
          parent_path = "/parent"
          published   = true
          source = {
            configuration = "{ \"see\": \"documentation\" }"
            type          = "http-fetcher"
          }
          type       = "MARKDOWN"
          updated_at = "2018-01-01T00:00:00Z"
          visibility = "PUBLIC"
        }
      ]
      type       = "MARKDOWN"
      updated_at = "2018-01-01T00:00:00Z"
      visibility = "PUBLIC"
    }
  ]
  plans = [
    {
      characteristics = [
        "..."
      ]
      description = "...my_description..."
      excluded_groups = [
        "..."
      ]
      flows = [
        {
          connect = [
            {
              condition         = "...my_condition..."
              configuration     = "{ \"see\": \"documentation\" }"
              description       = "...my_description..."
              enabled           = false
              message_condition = "...my_message_condition..."
              name              = "...my_name..."
              policy            = "...my_policy..."
            }
          ]
          enabled = false
          id      = "4e6abbd2-c0c6-462d-be9e-6371209af34b"
          interact = [
            {
              condition         = "...my_condition..."
              configuration     = "{ \"see\": \"documentation\" }"
              description       = "...my_description..."
              enabled           = true
              message_condition = "...my_message_condition..."
              name              = "...my_name..."
              policy            = "...my_policy..."
            }
          ]
          name = "My Flow"
          publish = [
            {
              condition         = "...my_condition..."
              configuration     = "{ \"see\": \"documentation\" }"
              description       = "...my_description..."
              enabled           = false
              message_condition = "...my_message_condition..."
              name              = "...my_name..."
              policy            = "...my_policy..."
            }
          ]
          request = [
            {
              condition         = "...my_condition..."
              configuration     = "{ \"see\": \"documentation\" }"
              description       = "...my_description..."
              enabled           = true
              message_condition = "...my_message_condition..."
              name              = "...my_name..."
              policy            = "...my_policy..."
            }
          ]
          response = [
            {
              condition         = "...my_condition..."
              configuration     = "{ \"see\": \"documentation\" }"
              description       = "...my_description..."
              enabled           = false
              message_condition = "...my_message_condition..."
              name              = "...my_name..."
              policy            = "...my_policy..."
            }
          ]
          selectors = [
            {
              http = {
                methods = [
                  "GET"
                ]
                path          = "/my/path"
                path_operator = "EQUALS"
                type          = "HTTP"
              }
            }
          ]
          subscribe = [
            {
              condition         = "...my_condition..."
              configuration     = "{ \"see\": \"documentation\" }"
              description       = "...my_description..."
              enabled           = false
              message_condition = "...my_message_condition..."
              name              = "...my_name..."
              policy            = "...my_policy..."
            }
          ]
          tags = [
            "tag1",
            "tag2",
          ]
        }
      ]
      general_conditions = "...my_general_conditions..."
      hrid               = "my_demo_api"
      mode               = "STANDARD"
      name               = "...my_name..."
      security = {
        configuration = "{ \"see\": \"documentation\" }"
        type          = "KEY_LESS"
      }
      selection_rule = "...my_selection_rule..."
      status         = "STAGING"
      tags = [
        "..."
      ]
      type       = "API"
      validation = "AUTO"
    }
  ]
  primary_owner = {
    display_name = "John Doe"
    email        = "john.doe@example.com"
    id           = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
    type         = "USER"
  }
  properties = [
    {
      dynamic     = true
      encryptable = true
      encrypted   = true
      key         = "...my_key..."
      value       = "...my_value..."
    }
  ]
  resources = [
    {
      configuration = "{ \"see\": \"documentation\" }"
      enabled       = true
      name          = "...my_name..."
      type          = "...my_type..."
    }
  ]
  response_templates = {
    key = {
      # ...
    }
  }
  services = {
    dynamic_property = {
      configuration          = "{ \"see\": \"documentation\" }"
      enabled                = true
      override_configuration = false
      type                   = "...my_type..."
    }
  }
  state = "STARTED"
  tags = [
    "public",
  ]
  type       = "MESSAGE"
  version    = "1.0.0"
  visibility = "PUBLIC"
}