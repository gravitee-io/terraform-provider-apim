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

## Issues

We use GitHub issues to track bugs and enhancements. Found a bug in the source code? You want to propose new features or enhancements? 
Before you submit your issue, search the [issues repository](https://github.com/gravitee-io/issues/issues); maybe your question was already answered.

[Report a new issue](https://github.com/gravitee-io/issues/issues/new/choose)

> Issues are only to report bugs, request enhancements, or request new features. For general questions and discussions, use the [Community forum](https://community.gravitee.io?utm_source=contributing)

Providing the following information will help us to deal quickly with your issue :
* **Tag the issue with `project: GKO`**
* **Overview of the issue** : describe the issue and why this is a bug for you.
* **Gravitee.io version(s)** : possible regression?
* **Terraform CLI version** : 1.9.* ...
* **Terraform Provider version** : 0.3.1 ...
* **You have logs, traces, server error stack traces?** add these to the issue's description.

<!-- Start Table of Contents [toc] -->
## Table of Contents
<!-- $toc-max-depth=2 -->
* [Gravitee - APIM Terraform Provider](#gravitee-apim-terraform-provider)
  * [Issues](#issues)
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
      version = "0.4.0"
    }
  }
}

provider "apim" {
  server_url = "..." # Optional - can use APIM_SERVER_URL environment variable
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
