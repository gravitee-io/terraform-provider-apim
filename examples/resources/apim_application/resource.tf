resource "apim_application" "my_application" {
  background     = "...my_background..."
  description    = "My Demo Application"
  domain         = "...my_domain..."
  environment_id = "...my_environment_id..."
  groups = [
    "..."
  ]
  hrid = "my_demo_api"
  members = [
    {
      id        = "...my_id..."
      role      = "...my_role..."
      source    = "...my_source..."
      source_id = "...my_source_id..."
    }
  ]
  metadata = [
    {
      default_value = "...my_default_value..."
      format        = "DATE"
      hidden        = false
      key           = "...my_key..."
      name          = "...my_name..."
      value         = "...my_value..."
    }
  ]
  name            = "My Application"
  notify_members  = false
  organization_id = "...my_organization_id..."
  picture_url     = "...my_picture_url..."
  primary_owner = {
    display_name = "John Doe"
    email        = "...my_email..."
    id           = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
    type         = "USER"
  }
  settings = {
    app = {
      client_id = "...my_client_id..."
      type      = "...my_type..."
    }
    oauth = {
      application_type = "...my_application_type..."
      grant_types = [
        "..."
      ]
      redirect_uris = [
        "..."
      ]
    }
    tls = {
      client_certificate = "...my_client_certificate..."
    }
  }
  status = "ARCHIVED"
}