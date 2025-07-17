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
  * [Testing the provider locally](#testing-the-provider-locally)

<!-- End Table of Contents [toc] -->

<!-- Start Installation [installation] -->
## Installation

To install this provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```hcl
terraform {
  required_providers {
    apim = {
      source  = "gravitee-io/apim"
      version = "0.1.0"
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

<!-- Start Testing the provider locally [usage] -->
## Testing the provider locally

#### Local Provider

Should you want to validate a change locally, the `--debug` flag allows you to execute the provider against a terraform instance locally.

This also allows for debuggers (e.g. delve) to be attached to the provider.

```sh
go run main.go --debug
# Copy the TF_REATTACH_PROVIDERS env var
# In a new terminal
cd examples/your-example
TF_REATTACH_PROVIDERS=... terraform init
TF_REATTACH_PROVIDERS=... terraform apply
```

#### Compiled Provider

Terraform allows you to use local provider builds by setting a `dev_overrides` block in a configuration file called `.terraformrc`. This block overrides all other configured installation methods.

1. Execute `go build` to construct a binary called `terraform-provider-apim`
2. Ensure that the `.terraformrc` file is configured with a `dev_overrides` section such that your local copy of terraform can see the provider binary

Terraform searches for the `.terraformrc` file in your home directory and applies any configuration settings you set.

```
provider_installation {

  dev_overrides {
      "registry.terraform.io/gravitee-io/apim" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```
<!-- End Testing the provider locally [usage] -->

<!-- Testing Instruction -->

### Automated Testing

Automated tests take Terraform configuration(s) and then perform create, read, import, update, and delete against those using real Terraform commands. Automated testing also supports writing assertions against the Terraform plan or state per test step.

The testing code is conventionally written as a Go test file (`internal/provider/xxx_resource_test.go`) while configurations are in `internal/provider/testdata` directories named after the test. For example:

* [Example provider testing code](https://github.com/speakeasy-sdks/terraform-provider-gravitee/blob/master/internal/provider/apiv4_resource_test.go)
* [Example provider testing configuration](https://github.com/speakeasy-sdks/terraform-provider-gravitee/blob/master/internal/provider/testdata/TestAPIV4Resource_lifecycle/main.tf).

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
cd gravitee-kubernetes-operator
make start-cluster
```

Run:

```shell
go test -count=1 -timeout=10m -v ./internal/provider
```


<!-- Placeholder for Future Speakeasy SDK Sections -->
