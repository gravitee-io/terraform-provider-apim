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
        # This header will returned and then used
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
