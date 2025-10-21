---
page_title: "Setup"
---

# Setup

## Pre-requisites
Before configuring the provider, check the following information is well known:

* The host and port for your Management API 
  * It is automatically configured for [Gravitee Cloud](#gravitee-cloud)
* Automation API is enabled
  * It is enabled by default on Gravitee Cloud
  * Check with your DevOps team or with Gravitee if your installation is managed 
* Credentials
    * Service account [Define an APIM service account for Terraform](https://documentation.gravitee.io/apim/terraform/define-an-apim-service-account-for-terraform) (Gravitee documentation website)
    * or Gravitee Cloud token
* (Optional) For a multi-tenant installation:
    * Organization ID
    * Environment ID

## Environment variables

All the above can be configured using [environment variables overrides](../envvar).

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

~> `environment_id` becomes mandatory if the token was crafted to work on multiple environments.