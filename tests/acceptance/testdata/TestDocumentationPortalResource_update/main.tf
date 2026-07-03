provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "content" {
  type = string
}

variable "hrid" {
  type = string
}

variable "name" {
  type = string
}

variable "order" {
  type = number
}

variable "portal_hrid" {
  type = string
}

resource "apim_portal" "portal" {
  hrid = var.portal_hrid
  name = "Test Portal"
  navigation = [
    {
      path = "/docs-${var.hrid}"
    },
    {
      path = "/reference-${var.hrid}"
    }
  ]
}

resource "apim_documentation_portal" "test" {
  portal_hrid = var.portal_hrid
  hrid        = var.hrid
  name        = var.name
  type        = "GRAVITEE_MARKDOWN"
  content     = var.content
  location    = var.name == "Getting Started" ? "/docs-${var.hrid}" : "/reference-${var.hrid}"
  order       = var.order
}
