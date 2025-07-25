---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "apim Provider"
subcategory: ""
description: |-
  Gravitee: APIM Terraform Provider (alpha)
  Manage APIs and Shared Policy Groups with Terraform
---

# apim Provider

Gravitee: APIM Terraform Provider (alpha)

Manage APIs and Shared Policy Groups with Terraform

## Example Usage

```terraform
terraform {
  required_providers {
    apim = {
      source  = "gravitee-io/apim"
      version = "0.2.11"
    }
  }
}

provider "apim" {
  # Configuration options
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `bearer_auth` (String, Sensitive)
- `environment_id` (String) Id of an environment.
- `organization_id` (String) Id of an organization.
- `password` (String, Sensitive)
- `server_url` (String) Server URL (defaults to https://apim-master-api.team-apim.gravitee.dev/automation)
- `username` (String, Sensitive)
