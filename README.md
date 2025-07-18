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
      version = "4.8.0-alpha.5"
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
* [apim_shared_policy_group](docs/resources/shared_policy_group.md)
### Data Sources

* [apim_apiv4](docs/data-sources/apiv4.md)
* [apim_shared_policy_group](docs/data-sources/shared_policy_group.md)
<!-- End Available Resources and Data Sources [operations] -->

<!-- No End Testing the provider locally [usage] -->

<!-- Testing Instruction -->

### Automated Testing

Automated tests take Terraform configuration(s) and then perform create, read, import, update, and delete against those using real Terraform commands. Automated testing also supports writing assertions against the Terraform plan or state per test step.

The testing code is conventionally written as a Go test file (`internal/provider/xxx_resource_test.go`) while configurations are in `internal/provider/testdata` directories named after the test. For example:

* [Example provider testing code](internal/provider/apiv4_resource_test.go)
* [Example provider testing configuration](internal/provider/testdata/TestAPIV4Resource_lifecycle/main.tf).

To test the APIM Terraform Provider you need to set environment variables:

```shell
export APIM_TOKEN=xxxx 
export APIM_SERVER_URL=http://localhost:30083/automation
export APIM_ORG_ID=DEFAULT
export APIM_ENV_ID=DEFAULT
export TF_ACC=1
```

Content of `.terraformrc` must be 

```hcl
provider_installation {

  dev_overrides {
      "registry.terraform.io/hashicorp/gravitee" = "/path/to/terraform-provider-apim"
  }
  
  direct {}
}
```

You can use GKO make targets to start an APIM cluster

```shell
cd ../gravitee-kubernetes-operator
make start-cluster
```

Run:

```shell
cd ../terraform-provider-apim
go test -count=1 -timeout=10m -v ./internal/provider
```


<!-- Placeholder for Future Speakeasy SDK Sections -->
