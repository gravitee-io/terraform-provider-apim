resource "apim_dictionary" "dynamic" {
  hrid        = "dynamic"
  name        = "[Terraform] Dynamic dictionary"
  description = "Expose all headers of Gravitee echo API as properties"
  deployed    = true
  type        = "DYNAMIC"
  dynamic = {
    provider = {
      http = {
        type   = "HTTP"
        url    = "https://api.gravitee.io/echo"
        method = "GET"
        # This header will be returned and then used
        # as a property in the API policy
        headers = [
          {
            name  = "X-Test-Specific"
            value = "ABCDEF"
          }
        ]
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

variable "api_key" {
  type      = string
  sensitive = true
}

# Requires APIM 4.13 or later
resource "apim_dictionary" "manual" {
  hrid     = "manual"
  name     = "[Terraform] Manual dictionary"
  deployed = true
  type     = "MANUAL"
  manual = {
    properties = {
      env = "production"
    }
    # Encrypted by APIM and never returned
    encrypted_properties = {
      "api-key" = {
        value = var.api_key
      }
    }
  }
}
