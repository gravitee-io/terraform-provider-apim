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

provider "apim" {
  // all is configured via env vars
}

data "local_file" "cert1" {
  filename = "client1.crt"
}

data "local_file" "cert2" {
  filename = "client2.crt"
}