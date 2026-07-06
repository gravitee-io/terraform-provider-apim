---
page_title: "Environment variables override"
---

# Environment variables override

Gravitee Terraform provider can be configured using the following environment variables

| Variable         | Corresponding provider field |
|------------------|------------------------------|
| APIM_SERVER_URL  | server_url                   |
| APIM_SA_TOKEN    | bearer_auth                  |
| APIM_ORG_ID      | organization_id              |
| APIM_ENV_ID      | environment_id               |
| APIM_CLOUD_TOKEN | cloud_auth                   |

If you use only environment variables to configure your provider, then your configuration block looks like this:

```hcl
provider "apim" {
  # Hooray, nothing is hardcoded!
}
```