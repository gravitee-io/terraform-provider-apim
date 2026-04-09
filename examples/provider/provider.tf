terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

# Gravitee Terraform provider to manage Gravitee API Management assets
# server_url and bearer_auth is required only when cloud_auth is empty.
# Usage of APIM_* environment variables is encouraged.
provider "apim" {

  # Gravitee Cloud Token
  #
  # Can be set/overriden with APIM_CLOUD_TOKEN environment variable.
  # organization_id and environment_id server_url and automatically
  # set when cloud_auth contains a Gravitee cloud token.
  cloud_auth = "eyJhbGciOiJSUzI1NiIsICJ0eXAiOiAiSldUIn0...."

  # Gravitee Automation API URL
  #
  # Can be set/overriden with APIM_SERVER_URL environment variable
  # Automatically set with when cloud_auth is set
  server_url = "http://localhost:8083/automation"

  # API Management Service Account token
  #
  # Can be set/overriden with APIM_SA_TOKEN environment variable
  # Ignored when cloud_auth is set.
  bearer_auth = "2435eff45e45-..."

  # Gravitee Organization UUID (defaulted to "DEFAULT")
  # Can be set/overriden with APIM_ORGANIZATION_ID environment variable
  # Automatically set with when cloud_auth is set
  organization_id = "xxx"

  # Gravitee Environment UUID (defaulted to "DEFAULT")
  # Can be set/overriden with APIM_ENVIRONMENT_ID environment variable
  # Automatically set with when cloud_auth is set
  environment_id = "xxx"

  # Basic auth (username/password) is also supported but discourage
  # for security reason and should only be used for testing purposes.

}