provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "method" {
  type = string
}

variable "deployed" {
  type = bool
}

variable "headers" {
  type = list(
    object({
      name  = string
      value = string
    })
  )
}

resource "apim_dictionary" "test" {
  hrid        = var.hrid
  name        = "[Terraform] Dynamic dictionary"
  description = "Expose all heders of Gravitee echo API"
  deployed    = var.deployed
  type        = "DYNAMIC"
  dynamic = {
    provider = {
      http = {
        type          = "HTTP"
        url           = "https://api.gravitee.io/echo"
        method        = var.method
        headers       = var.headers
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
