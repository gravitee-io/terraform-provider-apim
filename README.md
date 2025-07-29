# Gravitee - APIM Terraform Provider

<picture>
  <source media="(prefers-color-scheme: dark)" srcset=".assets/gravitee-logo-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset=".assets/gravitee-logo-light.svg">
  <img alt="Gravitee.io" width="400">
</picture>

<!-- Start Summary [summary] -->
## Summary

Gravitee: APIM Terraform Provider (alpha)

Manage APIs and Shared Policy Groups with Terraform

Compatible with APIM 4.8 and above
<!-- End Summary [summary] -->

<!-- Start Table of Contents [toc] -->
## Table of Contents
<!-- $toc-max-depth=2 -->
* [Gravitee - APIM Terraform Provider](#gravitee-apim-terraform-provider)
  * [Installation](#installation)
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
      version = "0.3.2"
    }
  }
}

provider "apim" {
  # Configuration options
}
```
<!-- End Installation [installation] -->

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
