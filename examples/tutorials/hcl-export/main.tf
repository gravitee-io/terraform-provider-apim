terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

provider "apim" {
  http_headers = {
    "X-Gravitee-Set-Hrid" = "true"
  }
}

