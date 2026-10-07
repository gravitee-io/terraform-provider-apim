resource "apim_group" "example" {
  hrid = "example"
  name = "Example"
  default_member_roles = {
    API         = "USER"
    APPLICATION = "USER"
    API_PRODUCT = "USER"
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
