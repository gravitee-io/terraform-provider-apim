terraform {
  required_providers {
    apim = {
      source  = "gravitee-io/apim"
      version = "0.3.1"
    }
  }
}

provider "apim" {
  bearer_auth = "access-token"
}