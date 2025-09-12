---
page_title: "Gravitee APIM Provider Configuration Reference"
subcategory: ""
description: |-
  Documentation for the APIM Terraform provider.
---

# Gravitee APIM Provider Configuration Reference

The `apim` code block defines the organization settings, credentials, and endpoints that are required by the Gravitee provider. These parameters identify and provide access to a specific Gravitee instance.

This example shows how to define a self-hosted Gravitee deployment:

```terraform
provider "apim" {
  server_url      = "https://<mAPI host and port>/automation"
  bearer_auth     = "c7783347-f1bc-45fd-8199-d2ef18d24717" 
}
```

For more information about the provider, see the [APIM provider documentation](https://documentation.gravitee.io/apim/terraform/configure-the-gravitee-provider#configuration-block).