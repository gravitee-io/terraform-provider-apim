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
      source  = "gravitee-io/apim"
      version = "0.3.0"
    }
  }
}

provider "apim" {
  # Configuration options
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

### Resources

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
