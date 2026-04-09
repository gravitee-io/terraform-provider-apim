---
page_title: "Certificate rotation"
subcategory: "Tutorials"
---

This tutorial demonstrates how to perform mTLS certificates rotation.

This is available starting with APIM 4.11

## Provider configuration

We use external local files, but content can be inlined.

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

provider "apim" {
  // all is configured via env vars
}

data "local_file" "cert1" {
  filename = "client1.crt"
}

data "local_file" "cert2" {
  filename = "client2.crt"
}
```

## Resource

In this state, `cert1` is active until the end of August and `cert2` will be active the 1 of July.

```terraform
resource "apim_application" "backend-to-backend-multi-certs" {
  # should match the resource name
  hrid        = "backend-to-backend-multi-certs"
  name        = "[Terraform] Application for Backend to Backend OAuth"
  description = "Demonstrates applications for OAuth with certificate can be created with Terraform"
  domain      = "example.com"
  settings = {
    oauth = {
      application_type = "backend_to_backend"
      redirect_uris    = []
      grant_types = [
        "client_credentials"
      ]
    }
    tls = {
      client_certificates = [{
        name    = "cert1"
        content = data.local_file.cert1.content
        ends_at = "2026-08-01T00:00:00Z"
        }, {
        name      = "cert2"
        content   = data.local_file.cert2.content
        starts_at = "2026-04-01T00:00:00Z"
      }]
    }
  }
}

```

The rotation can be done "manually" (i.e., without dates) by adding the new cert (`cert2`) when it needs to be added.
Applying the config will propagate this new cert and make it available to call APIM Gateway on top of `cert1`.

Later one can remove `cert1` from the config file when all clients are using `cert2`, apply to eventually remove `cert1` from the APIM Gateway.
