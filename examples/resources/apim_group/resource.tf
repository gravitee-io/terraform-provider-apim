resource "apim_group" "example" {
  hrid = "example"
  name = "Example"
  default_member_roles = {
    api         = "USER"
    application = "USER"
    api_product = "USER"
  }
  members = [
    {
      roles = {
        API         = "OWNER"
        APPLICATION = "USER"
        INTEGRATION = "USER"
      }
      source    = "memory"
      source_id = "api1"
    }
  ]
}
