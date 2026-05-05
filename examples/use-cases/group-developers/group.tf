resource "apim_group" "developers" {
  hrid = "developers"
  name = "Developers"
  members = [
    {
      source    = "memory"
      source_id = "api1"
      roles = {
        API         = "OWNER"
        APPLICATION = "USER"
        INTEGRATION = "USER"
      }
    },
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
