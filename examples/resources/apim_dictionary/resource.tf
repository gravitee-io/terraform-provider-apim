resource "apim_dictionary" "my_dictionary" {
  deployed    = false
  description = "...my_description..."
  dynamic = {
    provider = {
      http = {
        body = "...my_body..."
        headers = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        method           = "CONNECT"
        specification    = "...my_specification..."
        type             = "HTTP"
        url              = "https://cautious-futon.com/"
        use_system_proxy = false
      }
    }
    trigger = {
      rate = 4
      unit = "MILLISECONDS"
    }
  }
  environment_id = "a44e0d1b-9fa9-4d64-8b76-3634623a2e27"
  hrid           = "demo_api"
  manual = {
    properties = {
      key = "value"
    }
  }
  name            = "...my_name..."
  organization_id = "dedd0e0f-b3e9-4d2f-89cd-b2a9de7cb145"
  type            = "DYNAMIC"
}