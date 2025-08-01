terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

provider "apim" {
  # Can be overriden with APIM_SERVER_URL environment variable
  # server_url = "http://localhost:8083/automation"
  # Can be overriden with APIM_SA_TOKEN environment variable
  # bearer_auth = "xxx" # your account personal token
  # Basic auth is also supported but discourage for security reason and should only be used for tests.
}
