resource "apim_group" "example" {
  hrid = "example"
  name = "Example"
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
