# Gravitee - APIM Terraform Provider

<picture>
  <source media="(prefers-color-scheme: dark)" srcset=".assets/gravitee-logo-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset=".assets/gravitee-logo-light.svg">
  <img alt="Gravitee.io" width="400">
</picture>

<!-- Start Summary [summary] -->
## Summary

Gravitee: Gravitee API Management Terraform Provider (beta)

You can manage with Terraform the following:
* APIs
* Shared Policy Groups
* Applications
* Subscriptions

[Go to our documentation web site for more about configuration, capabilities and examples](https://documentation.gravitee.io/apim/terraform) 

Compatible with APIM 4.9 and above

Checkout other sections to configure, authenticate and start working with Gravitee resources
<!-- End Summary [summary] -->

<!-- Start Table of Contents [toc] -->
## Table of Contents
<!-- $toc-max-depth=2 -->
* [Gravitee - APIM Terraform Provider](#gravitee-apim-terraform-provider)
  * [Installation](#installation)
  * [Authentication](#authentication)
  * [Available Resources and Data Sources](#available-resources-and-data-sources)

<!-- End Table of Contents [toc] -->

<!-- Start Installation [installation] -->
## Installation

To install this provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```hcl
terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

# Gravitee Terraform provider to manage Gravitee API Management assets
# server_url and bearer_auth is required only when cloud_auth is empty.
# Usage of APIM_* environment variables is encouraged.
provider "apim" {

  # Gravitee Cloud Token
  #
  # Can be set/overriden with APIM_CLOUD_TOKEN environment variable.
  # organization_id and environment_id server_url and automatically
  # set when cloud_auth contains a Gravitee cloud token.
  cloud_auth = "eyJhbGciOiJSUzI1NiIsICJ0eXAiOiAiSldUIn0...."

  # Gravitee Automation API URL
  #
  # Can be set/overriden with APIM_SERVER_URL environment variable
  # Automatically set with when cloud_auth is set
  server_url = "http://localhost:8083/automation"

  # API Management Service Account token
  #
  # Can be set/overriden with APIM_SA_TOKEN environment variable
  # Ignored when cloud_auth is set.
  bearer_auth = "2435eff45e45-..."

  # Gravitee Organization UUID (defaulted to "DEFAULT")
  # Can be set/overriden with APIM_ORGANIZATION_ID environment variable
  # Automatically set with when cloud_auth is set
  organization_id = "xxx"

  # Gravitee Environment UUID (defaulted to "DEFAULT")
  # Can be set/overriden with APIM_ENVIRONMENT_ID environment variable
  # Automatically set with when cloud_auth is set
  environment_id = "xxx"

  # Basic auth (username/password) is also supported but discourage
  # for security reason and should only be used for testing purposes.

}
```
<!-- End Installation [installation] -->

<!-- Start Authentication [security] -->
## Authentication

This provider supports authentication configuration via environment variables and provider configuration.

The configuration precedence is:

- Provider configuration
- Environment variables

Available configuration:

| Provider Attribute | Description |
|---|---|
| `bearer_auth` | Service account authentication. Configurable via environment variable `APIM_SA_TOKEN`. |
| `cloud_auth` | Gravitee Cloud Token authentication. Configurable via environment variable `APIM_CLOUD_TOKEN`. |
| `password` | Basic authentication password. Configurable via environment variable `APIM_PASSWORD`. |
| `username` | Basic authentication username. Configurable via environment variable `APIM_USERNAME`. |
<!-- End Authentication [security] -->

<!-- Start Available Resources and Data Sources [operations] -->
## Available Resources and Data Sources

### Managed Resources

* [apim_apiv4](docs/resources/apiv4.md)
* [apim_application](docs/resources/application.md)
* [apim_shared_policy_group](docs/resources/shared_policy_group.md)
* [apim_subscription](docs/resources/subscription.md)

### Data Sources

* [apim_apiv4](docs/data-sources/apiv4.md)
* [apim_application](docs/data-sources/application.md)
* [apim_shared_policy_group](docs/data-sources/shared_policy_group.md)
* [apim_subscription](docs/data-sources/subscription.md)
<!-- End Available Resources and Data Sources [operations] -->

<!-- No End Testing the provider locally [usage] -->

<!-- Placeholder for Future Speakeasy SDK Sections -->
