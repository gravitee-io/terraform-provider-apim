resource "apim_application" "with_group" {
  # should match the resource name
  hrid        = "with_group"
  name        = "[Terraform] Application bound to a group"
  description = "Demonstrate an application with a group"
  settings = {
    app = {
      type = "test"
    }
  }
  metadata = [{
    name   = "hello"
    format = "STRING"
  }]
  groups = [
    apim_group.developers.hrid
  ]
}

resource "apim_group" "developers" {
  hrid = "app-developers"
  name = "App Developers"
  members = [
    {
      source    = "memory"
      source_id = "application1"
      roles = {
        API         = "USER"
        APPLICATION = "OWNER"
        INTEGRATION = "USER"
      }
    }
  ]
}
