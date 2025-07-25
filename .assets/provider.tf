terraform {
  required_providers {
    apim = {
      source  = "gravitee-io/apim"
      version = "4.8.0-alpha"
    }
  }
}

provider "apim" {
  # Can be overriden with APIM_SERVER_URL environment variable
  # server_url = "http://localhost:8083/automation"
  # Can be overriden with APIM_SA_TOKEN environment variable
  # beaer_auth = "xxx" # your account personal token
  # Can be overriden with APIM_CLOUD_TOKEN environment variable
  # cloud_auth = "xxx" # Gravitee JWT Cloud token. You must not use bearer_auth in that case and server_url should be https://eu.cloudgate.gravitee.io
  # Basic auth is also supported but discourage for security reason and should only be used for tests.
}
