terraform {
  required_providers {
    apim = {
      source  = "gravitee-io/apim"
      version = "0.5.2"
    }
  }
}

provider "apim" {
  server_url = "..." # Optional - can use APIM_SERVER_URL environment variable
}