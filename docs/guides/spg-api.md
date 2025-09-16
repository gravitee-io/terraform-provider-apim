---
page_title: "API & Shared Policy Group"
subcategory: "Shared Policy Group"
---

# Use a Shared Policy Group to curate the headers of a v4 HTTP proxy API

The following example configures the Shared Policy Group resource to use the Gravitee
[Transform Headers policy](https://documentation.gravitee.io/apim/create-and-configure-apis/apply-policies/policy-reference/transform-headers)
on the Request phase of a v4 HTTP proxy API.
The resource removes the header "User-Agent" and adds a header named "X-Content-Path" that contains the API's context path.

```HCL
resource "apim_shared_policy_group" "curate_headers" {
  # should match the resource name
  hrid        = "curate_headers"
  name        = "[Terraform] Curated headers"
  description = "Simple Shared Policy Group that contains one step to remove User-Agent header and add X-Content-Path that contains this API context path"
  api_type    = "PROXY"
  phase       = "REQUEST"
  steps = [
    {
      enabled = true
      name    = "Curate headers"
      policy = "transform-headers"
      # Configuration is JSON as the schema depends on the policy used
      configuration = jsonencode({
        scope = "REQUEST"
        addHeaders = [
          {
            name  = "X-Context-Path"
            value = "{#request.contextPath}"
          }
        ],
        removeHeaders = ["User-Agent"]
      })
    }
  ]
}

resource "apim_apiv4" "simple-api-with-spg" {
  # should match the resource name
  hrid            = "simple-api-with-spg"
  name            = "[Terraform] Simple PROXY API With Shared Policy Group"
  description     = "A simple API that routes traffic to gravitee echo API. It curates headers using curate_headers Shared Policy Group."
  version         = "1"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PRIVATE"
  lifecycle_state = "UNPUBLISHED"
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
            path = "/simple-api-with-spg/"
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
          name   = "Default HTTP proxy"
          type   = "http-proxy"
          weight = 1
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
            type         = "HTTP"
            path         = "/"
            pathOperator = "STARTS_WITH"
            methods = []
          }
        }
      ]
      request = [
        {
          enabled = true
          # Known limitation, this name should the same as SPG to avoid terraform plan to remain inconsistent
          name   = "[Terraform] Curated headers"
          policy = "shared-policy-group-policy",
          # Configuration is JSON as the schema depends on the policy used
          configuration = jsonencode({
            hrid = apim_shared_policy_group.curate_headers.hrid
          })
        }
      ]
    }
  ]
  analytics = {
    enabled = false
  }
  # known limitation, some default value is returned by default which appears to be remove during plan
  metadata=[{}]
  # known limitation: will dispear in next version
  definition_context = {}
  plans = {
    # known limitation, key have to match name to avoid terraform plan to remain inconsistent
    KeyLess = {
      name        = "KeyLess"
      type        = "API"
      mode        = "STANDARD"
      validation  = "AUTO"
      status      = "PUBLISHED"
      description = "This plan does not require any authentication"
      security = {
        type = "KEY_LESS"
      }
    }
  }
}

```
