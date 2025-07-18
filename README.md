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
      version = "4.8.0-alpha"
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

### Run acceptance tests

Automated tests take Terraform configuration(s) and then perform create, read, import, update, and delete against those using real Terraform commands. Automated testing also supports writing assertions against the Terraform plan or state per test step.

The testing code is conventionally written as a Go test file (`internal/provider/xxx_resource_test.go`) while configurations are in `internal/provider/testdata` directories named after the test. For example:

* [Example provider testing code](v4/internal/provider/apiv4_resource_test.go)
* [Example provider testing configuration](v4/internal/provider/testdata/TestAPIV4Resource_lifecycle/main.tf).

To test the APIM Terraform Provider you need to set environment variables (examples):

```shell
export APIM_SA_TOKEN=xxxx 
# or
export APIM_USERNAME=admin
export APIM_PASSWORD=admin
# to test elsewhere than master env
export APIM_SERVER_URL=http://localhost:30083/automation
# optionally
export APIM_ORG_ID=DEFAULT
export APIM_ENV_ID=DEFAULT
```

You can use GKO make targets to start an APIM local cluster

```shell
cd ../gravitee-kubernetes-operator
make start-cluster
```

Run:

```shell
cd ../terraform-provider-apim
make acceptance-tests
```

### Debug 
```sh
go run main.go --debug
# Copy the TF_REATTACH_PROVIDERS env var
# In a new terminal
cd examples/your-example
TF_REATTACH_PROVIDERS=... terraform init
TF_REATTACH_PROVIDERS=... terraform apply
```

<!-- Placeholder for Future Speakeasy SDK Sections -->
