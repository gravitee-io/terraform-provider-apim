---
page_title: "Export API to HCL"
subcategory: "Tutorials"
---
This tutorial demonstrates how to export an API to HCL configuration using the `terraform plan -generate-config-out` command.

~> This feature is experimental from the Terraform standpoint. It is available from the `0.5.0` release and works with Terraform from version 1.5.

It takes advantage of the READ operation during the plan to retrieve the API and then generates the HCL configuration.

Applications, Subscriptions and API are identified by `hrid`. An existing API created with console, however, only has an ID and has no `hrid`.

HRID are generated if no hrid is found from name where all that is not a letter or a number is replaced by a dash:
* API: from the API name
* Plan: from the plan name
* Page: from the page name
* Shared Policy Group: from the **step** name.

Coherence is kept between pages and parents as well as between page and plan general conditions.

Be aware that:
* If the HRID is the same as another API in your configuration, it will be updated and the other API will have its content replaced by the exported one.
* If page or plans do not have unique names, `hrid` won't be unique, which can lead to error when applying the API.
* Share Policy Group binding will only work if the Shared Policy Group was created with Terraform.
  Nevertheless, we create a HRID for SPG to avoid errors during export.

~> All cited files in that tutorial are located in [examples/tutorials/hcl-export](https://github.com/gravitee-io/terraform-provider-apim/tree/main/examples/tutorials/hcl-export)

## Provider configuration

The provider needs to express Gravitee API that `hrid`s need to be set, hence during the READ operation API UUID is used instead of the `hrid`.

```terraform
terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

provider "apim" {
  http_headers = {
    "X-Gravitee-Set-Hrid" = "true"
  }
}


```

In this example the provider configuration does not contain any other configuration.
It is assumed that APIM URL and credentials are already configured as environment variables.

## Import file

Terraform plan will do a READ operation to fulfill instruction in the import file that will eventually be written to a file.

```terraform
import {
  to = apim_apiv4.export
  id = <<-EOT
	    {
	      "organization_id": "DEFAULT",
	      "environment_id": "DEFAULT",
	      "hrid": "<<API ID>>"
	    }
	  EOT
}

```

You now need to replace `<<API ID>>` with the actual API ID. You can find it in the API URL of the Gravitee console.

The API URL is in the format `https://<<console URL>>/#!/<<environement>>/apis/30675173-378c-40c0-a751-73378c20c09c`.

~> You may need to adapt environment and organization as well depending on your setup.

## Run the plan to export the API as HCL

You can run the plan with the following command:

```bash
terraform plan -generate-config-out=exported.tf
```

The generated HCL configuration will be saved in the `exported.tf` file.

## Clean up (experimental)

In the tutorial directory you will a python script example to clean up the file after export.
* Uses `jsonencode()` instead of JSON strings
* Remove null values
* Remove empty arrays

```bash
python3 cleanup.py exported.tf exported-clean.tf
```

## Limitations

1. It uses a Terraform experimental feature (`-generate-config-out`).
2. It is limited to APIs
3. You can only export APIs one at a time.

## Example

### API in exported format (for reference)

This is what you would get from the standard export feature in the console.

```json
{
  "export": {
    "date": "2026-02-02T11:01:47.139884934Z",
    "apimVersion": "4.11.0-SNAPSHOT"
  },
  "api": {
    "definitionVersion": "V4",
    "type": "PROXY",
    "listeners": [
      {
        "type": "HTTP",
        "paths": [
          {
            "path": "/terraform/exported-example/",
            "overrideAccess": false
          }
        ],
        "pathMappings": [],
        "entrypoints": [
          {
            "type": "http-proxy",
            "qos": "AUTO",
            "configuration": {}
          }
        ],
        "servers": []
      }
    ],
    "endpointGroups": [
      {
        "name": "Default HTTP proxy group",
        "type": "http-proxy",
        "loadBalancer": {
          "type": "ROUND_ROBIN"
        },
        "sharedConfiguration": "{\"headers\":[{\"name\":\"Foo\",\"value\":\"Bar\"}],\"proxy\":{\"useSystemProxy\":false,\"enabled\":false},\"http\":{\"keepAliveTimeout\":30000,\"keepAlive\":true,\"propagateClientHost\":true,\"followRedirects\":false,\"readTimeout\":10000,\"idleTimeout\":0,\"connectTimeout\":3000,\"useCompression\":true,\"maxConcurrentConnections\":20,\"version\":\"HTTP_1_1\",\"pipelining\":false},\"ssl\":{\"keyStore\":{\"type\":\"\"},\"hostnameVerifier\":true,\"trustStore\":{\"type\":\"\"},\"trustAll\":false}}",
        "endpoints": [
          {
            "name": "Default HTTP proxy",
            "type": "http-proxy",
            "weight": 1,
            "inheritConfiguration": true,
            "configuration": {
              "target": "https://api.gravitee.io/echo"
            },
            "services": {},
            "secondary": false,
            "tenants": []
          }
        ],
        "services": {}
      }
    ],
    "analytics": {
      "enabled": true
    },
    "flowExecution": {
      "mode": "DEFAULT",
      "matchRequired": false
    },
    "flows": [
      {
        "name": "Default",
        "enabled": true,
        "selectors": [
          {
            "type": "HTTP",
            "path": "/",
            "pathOperator": "EQUALS",
            "methods": []
          }
        ],
        "request": [
          {
            "name": "Common OAS policy",
            "enabled": true,
            "policy": "shared-policy-group-policy",
            "configuration": {
              "sharedPolicyGroupId": "ffd6f930-6fbd-4f11-96f9-306fbd9f11e0"
            }
          }
        ],
        "response": [
          {
            "name": "Transform Headers",
            "enabled": true,
            "policy": "transform-headers",
            "configuration": {
              "whitelistHeaders": [],
              "addHeaders": [
                {
                  "name": "Warning",
                  "value": "None"
                }
              ],
              "scope": "REQUEST",
              "removeHeaders": []
            }
          }
        ],
        "subscribe": [],
        "publish": [],
        "connect": [],
        "interact": [],
        "tags": []
      }
    ],
    "name": "Automation exportable",
    "description": "",
    "apiVersion": "1",
    "deployedAt": "2026-02-02T11:01:09.045Z",
    "createdAt": "2026-02-02T10:58:09.998Z",
    "updatedAt": "2026-02-02T11:01:09.045Z",
    "disableMembershipNotifications": false,
    "groups": [],
    "state": "STOPPED",
    "visibility": "PRIVATE",
    "labels": [],
    "lifecycleState": "CREATED",
    "tags": [],
    "primaryOwner": {
      "id": "67f4af8a-e325-4289-b4af-8ae32532897d",
      "email": "benoit.bordigoni@graviteesource.com",
      "displayName": "Benoît Bordigoni",
      "type": "USER"
    },
    "categories": [],
    "originContext": {
      "origin": "MANAGEMENT"
    },
    "responseTemplates": {},
    "resources": [],
    "properties": []
  },
  "pages": [
    {
      "id": "d73e0602-460b-451d-be06-02460b251d00",
      "name": "Legal",
      "type": "FOLDER",
      "order": 0,
      "published": true,
      "visibility": "PUBLIC",
      "updatedAt": "2026-02-02T11:00:08.352Z",
      "configuration": {},
      "homepage": false,
      "excludedAccessControls": false,
      "accessControls": [],
      "metadata": {}
    },
    {
      "id": "dc0d8bda-6c16-4793-8d8b-da6c16279348",
      "name": "General Conditions",
      "type": "MARKDOWN",
      "content": "Anyone is free to use this example API",
      "order": 0,
      "published": true,
      "visibility": "PUBLIC",
      "updatedAt": "2026-02-02T11:00:42.905Z",
      "configuration": {},
      "homepage": false,
      "parentId": "d73e0602-460b-451d-be06-02460b251d00",
      "excludedAccessControls": false,
      "accessControls": [],
      "metadata": {}
    }
  ],
  "plans": [
    {
      "definitionVersion": "V4",
      "flows": [],
      "id": "fa417460-3df3-44b1-8174-603df3d4b178",
      "name": "Default Keyless (UNSECURED)",
      "description": "Default unsecured plan",
      "apiId": "1d5a2d23-df93-450c-9a2d-23df93950c33",
      "security": {
        "type": "KEY_LESS",
        "configuration": {}
      },
      "mode": "STANDARD",
      "characteristics": [],
      "commentRequired": false,
      "createdAt": "2026-02-02T10:58:12.414Z",
      "excludedGroups": [],
      "generalConditions": "dc0d8bda-6c16-4793-8d8b-da6c16279348",
      "order": 1,
      "publishedAt": "2026-02-02T10:58:12.557Z",
      "status": "PUBLISHED",
      "tags": [],
      "type": "API",
      "updatedAt": "2026-02-02T11:01:05.373Z",
      "validation": "MANUAL"
    }
  ],
  "apiMedia": []
}
```

### API in HCL format (before cleanup)

```terraform
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

```

### API in HCL format (after cleanup)

```terraform
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
```
