---
page_title: "Setup"
---

# Setup

## Pre-requisites
Before you configure your provider, ensure you have the following:

* The host and port for your Management API
* Automation API must be enabled
  * It is enabled on by default on Gravitee Cloud
  * Check with your DevOps team for on-prem or Gravitee for hosted environments
* Credentials:
    * Service account [Define an APIM service account for Terraform](https://documentation.gravitee.io/apim/terraform/define-an-apim-service-account-for-terraform) (Gravitee documentation website)
    * or Gravitee Cloud token
* (Optional) For a multi-tenant setup:
    * Organization ID
    * Environment ID

## Simple instance with a service account

```hcl
provider "apim" {
  server_url      = "https://<mAPI host and port>/automation"
  bearer_auth     = "service_account_token_goes_here"
}
```

In that case organization and environment IDs will be set to "`DEFAULT`".

## Multi-tenant instance with a service account

```hcl
provider "apim" {
  server_url      = "https://<mAPI host and port>/automation"
  bearer_auth     = "service_account_token_goes_here"
  organization_id = "organization_uuid_goes_here"
  environment_id = "environment_uuid_goes_here"
}
```

## Gravitee Cloud 

```hcl
provider "apim" {
  cloud_auth     = "your_gravitee_cloud_jwt_goes_here"
}
```

This will set `server_url`, `bearer_auth`, `organization_id`, `environment_id` for you.

~> `environment_id` becomes mandatory if the token is valid for multiple environment 