terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

provider "apim" {
  # Can be set/overriden with APIM_SERVER_URL environment variable
  # Automatically set with when cloud_token is set
  # server_url = "http://localhost:8083/automation"

  # Can be set/overriden with APIM_SA_TOKEN environment variable
  # bearer_auth = "xxx"       # your account personal token

  # Can be set/overriden with APIM_CLOUD_TOKEN environment variable
  # cloud_token = "xxx"       # your Gravitee Cloud token

  # Can be set/overriden with APIM_ORGANIZATION_ID environment variable
  # Automatically set with when cloud_token is set
  # organization_id = "xxx"   # your account personal token

  # Can be set/overriden with APIM_ENVIRONMENT_ID environment variable
  # Automatically set with when cloud_token is set
  # environment_id = "xxx"    # your account personal token

  # Basic auth is also supported but discourage for security reason and should only be used for tests.
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"

}
