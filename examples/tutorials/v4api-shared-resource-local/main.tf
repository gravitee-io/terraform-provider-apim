terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
    local = {
      source = "hashicorp/local"
    }
  }
}

data "local_file" "api-resource-basic-auth" {
  filename = "basic-auth-config.json"
}