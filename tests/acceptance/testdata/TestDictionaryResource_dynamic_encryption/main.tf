provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "encrypt_on_fetch" {
  type    = bool
  default = null
}

resource "apim_dictionary" "test" {
  hrid     = var.hrid
  name     = "[Terraform] Encrypted dynamic dictionary"
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
    encryption = var.encrypt_on_fetch == null ? null : {
      encrypt_on_fetch = var.encrypt_on_fetch
    }
  }
}
