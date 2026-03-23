resource "apim_application" "backend-to-backend-multi-certs" {
  # should match the resource name
  hrid        = "backend-to-backend-multi-certs"
  name        = "[Terraform] Application for Backend to Backend OAuth"
  description = "Demonstrates applications for OAuth with certificate can be created with Terraform"
  domain      = "example.com"
  settings = {
    oauth = {
      application_type = "backend_to_backend"
      redirect_uris    = []
      grant_types = [
        "client_credentials"
      ]
    }
    tls = {
      client_certificates = [{
        name : "cert1"
        content : data.local_file.cert1.content
        endsAt : "2026-08-01T00:00:00Z"
        }, {
        name : "cert2"
        content : data.local_file.cert2.content
        startsAt : "2026-07-01T00:00:00Z"
      }]
    }
  }
}