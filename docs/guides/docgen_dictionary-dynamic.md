---
page_title: "Dynamic Dictionary"
subcategory: "Dictionary"
---

# Dynamic Dictionary

## How Dictionary and API work together

This example demonstrates how to use a dynamic dictionary calling an API. The API every 5 seconds.

Note that the API feeding the dictionary is the same endpoint used as a back-end.

There is what happens when this dictionary and API are applied:

1. API is called with an extra header
1. [JOLT Transformation](https://github.com/bazaarvoice/jolt) is run
1. `headers` JSON object returned is transformed into a Gravitee property (list of `key` and `value` object)
1. Properties are made available to the Gateway runtime
1. API can resolve it and add a header via the `transform-headers` policy

To understand how the transformation works, see below.

```terraform
resource "apim_dictionary" "dynamic" {
  hrid        = "dynamic"
  name        = "[Terraform] Dynamic dictionary"
  description = "Expose all headers of Gravitee echo API as properties"
  deployed    = true
  type        = "DYNAMIC"
  dynamic = {
    provider = {
      http = {
        type   = "HTTP"
        url    = "https://api.gravitee.io/echo"
        method = "GET"
        # This header will returned and then used
        # as a property in the API policy
        headers = [
          {
            name  = "X-Test-Specific"
            value = "ABCDEF"
          }
        ]
        specification = <<-EOT
        [
          {
            "operation": "shift",
            "spec": {
              "headers": {
                "*": {
                  "$": "[#2].key",
                  "@": "[#2].value"
                }
              }
            }
          }
        ]
        EOT
      }
    }
    trigger = {
      rate = 5
      unit = "SECONDS"
    }
  }
}



resource "apim_apiv4" "proxy-with-dyn-dictionary" {
  # should match the resource name
  hrid            = "proxy-with-dyn-dictionary"
  name            = "[Terraform] Simple PROXY API with dynamic dict ref"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PUBLIC"
  lifecycle_state = "PUBLISHED"
  listeners = [
    {
      http = {
        type = "HTTP"
        entrypoints = [
          {
            type = "http-proxy"
          }
        ]
        paths = [
          {
            path = "/proxy-with-dyn-dictionary/"
          }
        ]
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
          name                  = "Default HTTP proxy"
          type                  = "http-proxy"
          weight                = 1
          inherit_configuration = false
          # Configuration is JSON as is depends on the type schema
          configuration = jsonencode({
            target = "https://api.gravitee.io/echo"
          })
        }
      ]
    }
  ]
  flow_execution = {
    mode           = "DEFAULT"
    match_required = false
  }
  flows = [
    {
      enabled = true
      selectors = [
        {
          http = {
            type          = "HTTP"
            path          = "/"
            path_operator = "STARTS_WITH"
            methods       = []
          }
        }
      ]
      request = [
        {
          enabled = true
          name    = "Add 1 header"
          policy  = "transform-headers"
          # Configuration is JSON as the schema depends on the policy used
          configuration = jsonencode({
            scope = "REQUEST"
            addHeaders = [
              {
                name = "X-Env"
                # HRID of the dictionary is used as the key
                value = "{#dictionaries['dynamic']['X-Test-Specific']}"
              }
            ]
          })
        }
      ]
    }
  ]
  analytics = {
    enabled = false
  }
  plans = [
    {
      hrid        = "KeyLess"
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

## JOLT Transformation

The above transformation turns the API response below:

```JSON
{
    "headers": {
        "Host": "api.gravitee.io",
        "User-Agent": "Gravitee.io/4.12.0-SNAPSHOT (gio-apim-apis; Gravitee.io - Rest APIs; 386aedc2-8b91-4050-aaed-c28b91305016)",
        "X-Gravitee-Request-Id": "9b6810b4-f6f1-4ebc-a810-b4f6f1eebcf4",
        "X-Gravitee-Transaction-Id": "9b6810b4-f6f1-4ebc-a810-b4f6f1eebcf4",
        "X-Test-Specific": "ABCDEF",
        "accept-encoding": "deflate, gzip"
    },
    "query_params": {},
    "bodySize": 0
}
```

Into this structure that Gravitee can use to resolve properties in the `transform-headers` policy.

```JSON
[
    {
        "key":"Host",
        "value":"api.gravitee.io"
    },
    {
        "key":"User-Agent",
        "value":"Gravitee.io/4.12.0-SNAPSHOT (gio-apim-apis; Gravitee.io - Rest APIs; 386aedc2-8b91-4050-aaed-c28b91305016)"
    },
    {
        "key":"X-Gravitee-Request-Id",
        "value":"9b6810b4-f6f1-4ebc-a810-b4f6f1eebcf4"
    },
    {
        "key":"X-Gravitee-Transaction-Id",
        "value":"9b6810b4-f6f1-4ebc-a810-b4f6f1eebcf4"
    },
    {
        "key":"X-Test-Specific",
        "value":"ABCDEF"
    },
    {
        "key":"accept-encoding",
        "value":"deflate, gzip"
    }
]```

## Calling the API

Calling `GET /proxy-with-dyn-dictionary/` will return the following API response:

```JSON
:

```JSON
{
    "headers": {
        "Host": "api.gravitee.io",
        "User-Agent": "Gravitee.io/4.12.0-SNAPSHOT (gio-apim-apis; Gravitee.io - Rest APIs; 386aedc2-8b91-4050-aaed-c28b91305016)",
        "X-Gravitee-Request-Id": "9b6810b4-f6f1-4ebc-a810-b4f6f1eebcf4",
        "X-Gravitee-Transaction-Id": "9b6810b4-f6f1-4ebc-a810-b4f6f1eebcf4",
        # This is one of the value that was populated from the call
        # That has been resolved from the dictionary
        "X-Env": "ABCDEF",
        "accept-encoding": "deflate, gzip"
    },
    "query_params": {},
    "bodySize": 0
}
```