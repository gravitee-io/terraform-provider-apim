terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

provider "apim" {
  # Can be set/overriden with APIM_SERVER_URL environment variable
  # Automatically set with when cloud_auth is set
  #server_url = "http://localhost:8083/automation"

  # Can be set/overriden with APIM_SA_TOKEN environment variable
  #bearer_auth = "xxx"       # your account personal token

  # Can be set/overriden with APIM_CLOUD_TOKEN environment variable
  #cloud_auth = "..." # your Gravitee Cloud token

  # organization_id and environment_id are set to "DEFAULT" by default.
  # To use a different value, override them.

  # Can be set/overriden with APIM_ORGANIZATION_ID environment variable
  # Automatically set with when cloud_auth is set
  #organization_id = "xxx"

  # Can be set/overriden with APIM_ENVIRONMENT_ID environment variable
  # Automatically set with when cloud_auth is set
  #environment_id = "xxx"

  # Basic auth is also supported but discourage for security reason and should only be used for testing purposes.

}
