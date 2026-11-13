---
page_title: "Dictionary Encryption"
subcategory: "Dictionary"
---

# Dictionary Encryption

-> **NOTE:** Dictionary encryption requires APIM 4.13.2 or later. Earlier versions ignore
`encrypted_properties` and `encryption` and store nothing. APIM 4.13.3 or later is recommended: it
reads secrets back masked and rejects a secret declared as plain, as described below.

This example stores a backend API key as a secret in a manual dictionary, encrypts every value a
dynamic dictionary fetches, and has an API send the secret to its backend.

## Secret properties

APIM stores the values of `encrypted_properties` encrypted. From APIM 4.13.3, a read returns each
key with its value masked. Terraform keeps the values it sent in its state, in plain text: protect
the state as you would the secret itself. From APIM 4.13.3, a value of `••••••••••••` keeps the
secret APIM stores for that key, for instance after an import.

A secret added or deleted in the Console plans a change: the next apply removes it or restores it.
A secret whose value changes in the Console does not: the mask only proves that a value is stored.
If a secret is made plain in the Console, the plan shows the secret changing without printing the
plain value, and the apply encrypts it again.

Removing a key from `encrypted_properties`, or removing the whole attribute, deletes the secret.
A key cannot be in both `properties` and `encrypted_properties`: APIM rejects it.

A manual dictionary needs at least one property, plain or secret.

## A key encrypted outside Terraform

If a key declared in `properties` is encrypted in the Console, APIM 4.13.3 and later reject every
apply of the dictionary with `dictionary.property.encryptedToPlain`. Move the key to
`encrypted_properties` to apply again. APIM 4.13.2 keeps the key encrypted instead, and Terraform
plans to add it to `properties` on every run until the key moves to `encrypted_properties`.

## Making a secret plain again

It takes two applies: remove the key from `encrypted_properties` and apply, then add it to
`properties` and apply. From APIM 4.13.3, a single apply moving the key is rejected with
`dictionary.property.encryptedToPlain`; APIM 4.13.2 keeps it encrypted.

## Fetched values

`encrypt_on_fetch = true` encrypts every value the dynamic dictionary's provider fetches. Omitting
`encryption` turns encryption off, so removing the block takes effect. Values already encrypted
stay encrypted: APIM keeps each key's classification.

## Using a secret

An API reads a secret property like a plain one, with the dictionary HRID and the key:
`{#dictionaries['encrypted']['backend-api-key']}`. The gateway decrypts it.

```terraform
variable "backend_api_key" {
  type      = string
  sensitive = true
  default   = "change-me"
}

resource "apim_dictionary" "encrypted" {
  # should match the resource name
  hrid     = "encrypted"
  name     = "[Terraform] Dictionary with a secret"
  deployed = true
  type     = "MANUAL"
  manual = {
    properties = {
      env = "test"
    }
    # Stored encrypted by APIM; from 4.13.3, a read returns it masked
    encrypted_properties = {
      "backend-api-key" = {
        value = var.backend_api_key
      }
    }
  }
}

resource "apim_dictionary" "dynamic-encrypted" {
  # should match the resource name
  hrid     = "dynamic-encrypted"
  name     = "[Terraform] Dynamic dictionary with encrypted values"
  deployed = true
  type     = "DYNAMIC"
  dynamic = {
    provider = {
      http = {
        type   = "HTTP"
        url    = "https://api.gravitee.io/echo"
        method = "GET"
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
    # Every fetched value is encrypted
    encryption = {
      encrypt_on_fetch = true
    }
  }
}

resource "apim_apiv4" "proxy-with-secret" {
  # should match the resource name
  hrid            = "proxy-with-secret"
  name            = "[Terraform] Simple PROXY API with a dictionary secret"
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
            path = "/proxy-with-secret/"
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
          name    = "Add the API key header"
          policy  = "transform-headers"
          configuration = jsonencode({
            scope = "REQUEST"
            addHeaders = [
              {
                name = "X-Api-Key"
                # Secret properties resolve like plain ones: HRID of the dictionary, then the key
                value = "{#dictionaries['encrypted']['backend-api-key']}"
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
