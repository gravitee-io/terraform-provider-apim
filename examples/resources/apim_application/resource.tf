resource "apim_application" "my_application" {
  background     = "https://upload.wikimedia.org/wikipedia/commons/d/df/Green_Red_Gradient_Background.png"
  description    = "This is the documentation explaining purpose of this Application  ."
  domain         = "examples.com"
  environment_id = "...my_environment_id..."
  groups = [
    "..."
  ]
  hrid = "my_demo_api"
  members = [
    {
      role      = "REVIEWER"
      source    = "system"
      source_id = "john.doe@example.com"
    }
  ]
  metadata = [
    {
      default_value = "...my_default_value..."
      format        = "STRING"
      key           = "...my_key..."
      name          = "...my_name..."
      value         = "...my_value..."
    }
  ]
  name            = "Example Application"
  notify_members  = false
  organization_id = "...my_organization_id..."
  picture_url     = "https://upload.wikimedia.org/wikipedia/fr/0/09/Logo_App_Store_d%27Apple.png"
  primary_owner = {
    display_name = "John Doe"
    email        = "john.doe@example.com"
    id           = "00f8c9e7-78fc-4907-b8c9-e778fc790750"
    type         = "USER"
  }
  settings = {
    app = {
      client_id = "example-client-id"
      type      = "web"
    }
    oauth = {
      additional_client_metadata = {
        key = "value"
      }
      application_type = "browser"
      grant_types = [
        "implicit"
      ]
      redirect_uris = [
        "..."
      ]
    }
    tls = {
      client_certificate = "-----BEGIN CERTIFICATE-----\nMIIDfjCCAmagAwIBAgIUfHj3mygGaOfd1u1Uj09L6vY5stcwDQYJKoZIhvcNAQEL\nBQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM\nGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yNTA4MDUxNTUyMTBaFw0yNjA4\nMDUxNTUyMTBaMGkxCzAJBgNVBAYTAlVTMQ4wDAYDVQQIDAVTdGF0ZTENMAsGA1UE\nBwwEQ2l0eTEVMBMGA1UECgwMT3JnYW5pemF0aW9uMRMwEQYDVQQLDApEZXBhcnRt\nZW50MQ8wDQYDVQQDDAZjbGllbnQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDW862KHvjkq0EtwZJO/xw+QoTnRB0qm4E5+1wspC1er6tOm3hTJqCzfKwQ\ngZQKoP1Eq1PhM8GzceeqGjh8VZJaDmWwiJZdk5fprrZ1Lvwwl010lnh4MEhtN0Dw\nlwHSZCQ/vSvEDWJXugiE4F1OvAgi2+lIR5uYfyy2U6YbhlcVPdGAboBAFSQnxECF\n1gDpc3dFarPXfO/X3yf/BzAHys6IyMyqvBbur3K2UTO4gJL+59/DEyAwx7ofwukj\nTWpgGNDXlNFYwKk9qTSTbxdcofAVCjrBCEDTdoPkvrr5SxI7dV/ha5y33iOI4VPV\no6vN/58RJz+ZMI0mbOBeluqBW+xBAgMBAAGjQjBAMB0GA1UdDgQWBBTjpQ+KfcmK\nw4hCptY8iK/LX9BOhzAfBgNVHSMEGDAWgBQYdcUWurMS8FEEMzcJlFm2d4Dk3DAN\nBgkqhkiG9w0BAQsFAAOCAQEAoyv0RhgEbRNmyFF6WoTeH4durjmZRe3SCtum0Mnv\n4TOGT4sstPdz0l24psroL33z3jtsY8IrbqnSfTXWbziSCanDXnMHOewLykgN0ld0\nPHa2i5naU5tMeGdWeM80ZTXU7GMiiCkgrRai/V7GkXNKYTIdBontiLpbxUaGLpjY\naMYoCmHIEizazQP9xaAtm40CkYub1o40kgyQULyrwftqrlRtKshfYmHB6yxYVz60\npikgTVupVbhYcNMLOVXO7Q31UEYfC7fxMGqzybXg67EhvzoykXhhYo3YqAjho2yh\num2oEO8b5eQVAwRaooVLh0uqjZCpfN2ozscPpiTM9Pj3xQ==\n-----END CERTIFICATE-----\n"
    }
  }
  status = "ACTIVE"
}