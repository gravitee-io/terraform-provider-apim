provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

resource "apim_group" "test" {
  hrid = var.hrid
  name = "Test"
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

