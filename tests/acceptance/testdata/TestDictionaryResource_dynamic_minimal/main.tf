provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

resource "apim_dictionary" "test" {
  hrid     = var.hrid
  name     = "[Terraform] Dynamic dictionary"
  deployed = false
  type     = "DYNAMIC"
  dynamic = {
    provider = {
      http = {
        type          = "HTTP"
        url           = "https://api.gravitee.io/echo"
        method        = "GET"
        specification = <<-EOT
        [
          {
            "operation": "shift",
            "spec": {
              "headers": {
                "*": {
                  "$": "[#2].key",
                  "@": "[#2].value"
                }
              }
            }
          }
        ]
        EOT
      }
    }
    trigger = {
      rate = 5
      unit = "SECONDS"
    }
  }
}
