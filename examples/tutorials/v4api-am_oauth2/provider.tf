terraform {
  # Write-only attributes, such as the value_wo of am_client_secret, need Terraform 1.11 or later.
  required_version = ">= 1.11"
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
    am = {
      source = "gravitee-io/am"
    }
  }
}

# Set with APIM_SERVER_URL and APIM_SA_TOKEN.
provider "apim" {}

# Set with AM_SERVER_URL and AM_SA_TOKEN.
provider "am" {}
