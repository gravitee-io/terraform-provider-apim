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

# Using "local" datasource to read the file.
data "local_file" "api-resource-basic-auth" {
  filename = "basic-auth-config.json"
}