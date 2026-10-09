provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}
variable "portal_hrid" {
  type = string
}

resource "apim_portal" "portal" {
  hrid = var.portal_hrid
  name = "Test Portal"
  navigation = [
    {
      path = "/docs-${var.portal_hrid}"
    }
  ]
}

resource "apim_documentation_portal" "test" {
  portal_hrid = apim_portal.portal.hrid
  hrid        = var.hrid
  name        = "Getting Started"
  type        = "GRAVITEE_MARKDOWN"
  content     = "# Hello"
  location    = "/docs-${var.portal_hrid}"
}
