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


// TODO test
// /management/v2/environments/DEFAULT/apis/_import/definition
// remove automation from APIM_SERVER_URL
// call import with payload
// get the ID
// run tf plan -generate-config-out=exported.tf
// check that API can be applied without plan changes
// destroy
// delete plan
// delete API