---
page_title: "Shared API resource configuration"
subcategory: "Tutorials"
---


---
page_title: "Shared API resource configuration"
subcategory: "Tutorials"
---

# Share API resource configuration between APIs

~> API resource is different from Terraform resource. It is used by Gravitee policies to use access data.

In this example we show that APIs using an *"In memory users resource"* and share its configuration so all users are the same across APIs.

## Provider configuration and resource configuration file

```terraform
terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
    local = {
      source = "hashicorp/local"
    }
  }
}

# Using "local" datasource to read the file.
data "local_file" "api-resource-basic-auth" {
  filename = "basic-auth-config.json"
}
```

This is the content of the resource configuration file `basic-auth-config.json`, it creates two users that can be used by the basic authentication policy.

```JSON
{
  "users": [
    {
      "username": "admin",
      "password": "admin",
      "roles": []
    },
    {
      "username": "user",
      "password": "password",
      "roles": []
    }
  ]
}
```

## Usage in APIs

Below, two APIs using the same API resource configuration to perform basic authentication.

### API 1

```terraform
resource "apim_apiv4" "simple-api-shared-resource-1" {
  # should match the resource name
  hrid            = "simple-api-shared-resource-1"
  name            = "[Terraform] Simple API with shared resource [1/2]"
  description     = "A simple API that routes traffic to gravitee echo API. Using basic auth configured in a shared resource"
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
            path = "/simple-api-shared-resource-1/"
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
          name = "Default HTTP proxy"
          type = "http-proxy"
          # Configuration is JSON as it is owned by the "http-proxy" endpoint plugin
          configuration = jsonencode({
            target = "https://api.gravitee.io/echo"
          })
        }
      ]
    }
  ]
  flows = [
    {
      enabled = true
      selectors = [
        {
          http = {
            type         = "HTTP"
            path         = "/"
            pathOperator = "STARTS_WITH"
          }
        }
      ]
      request = [
        {
          # Authentication policy
          name    = "Basic Authentication",
          enabled = true,
          policy  = "policy-basic-authentication",
          # Configuration is JSON as is depends on the
          configuration = jsonencode({
            authenticationProviders = [
              "In memory users"
            ]
            realm = "gravitee.io"
          })
        }
      ]
    }
  ]
  resources = [
    {
      enabled = true
      name    = "In memory users"
      type    = "auth-provider-inline-resource"
      # Where configuraiton file is included in the API resource
      configuration = data.local_file.api-resource-basic-auth.content
    }
  ]
  analytics = {
    enabled = false
  }
  plans = [
    {
      hrid        = "keyless"
      name        = "Key Less"
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

### API 2

```terraform
resource "apim_apiv4" "simple-api-shared-resource-1" {
  # should match the resource name
  hrid            = "simple-api-shared-resource-1"
  name            = "[Terraform] Simple API with shared resource [1/2]"
  description     = "A simple API that routes traffic to gravitee echo API. Using basic auth configured in a shared resource"
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
            path = "/simple-api-shared-resource-1/"
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
          name = "Default HTTP proxy"
          type = "http-proxy"
          # Configuration is JSON as it is owned by the "http-proxy" endpoint plugin
          configuration = jsonencode({
            target = "https://api.gravitee.io/echo"
          })
        }
      ]
    }
  ]
  flows = [
    {
      enabled = true
      selectors = [
        {
          http = {
            type         = "HTTP"
            path         = "/"
            pathOperator = "STARTS_WITH"
          }
        }
      ]
      request = [
        {
          # Authentication policy
          name    = "Basic Authentication",
          enabled = true,
          policy  = "policy-basic-authentication",
          # Configuration is JSON as is depends on the
          configuration = jsonencode({
            authenticationProviders = [
              "In memory users"
            ]
            realm = "gravitee.io"
          })
        }
      ]
    }
  ]
  resources = [
    {
      enabled = true
      name    = "In memory users"
      type    = "auth-provider-inline-resource"
      # Where configuraiton file is included in the API resource
      configuration = data.local_file.api-resource-basic-auth.content
    }
  ]
  analytics = {
    enabled = false
  }
  plans = [
    {
      hrid        = "keyless"
      name        = "Key Less"
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



# Share API resource configuration between APIs

~> API resource is different from Terraform resource. It is used by a Gravitee policies to use access data.

In this example we show that APIs using an *"In memory users resource"* and share its configuration so all users are the same between APIs.

## Provider configuration and resource configuration file

```terraform
terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
    local = {
      source = "hashicorp/local"
    }
  }
}

# Using "local" datasource to read the file.
data "local_file" "api-resource-basic-auth" {
  filename = "basic-auth-config.json"
}
```

This is the content of the resource configuration file, it creates two users that can be used by the basic authentication policy.

```JSON

```

## Usage in APIs

Below you will two APIs that use the same API resource configuration used to perform basic authentication.

### API 1

```terraform

```

### API 2

```terraform

```

