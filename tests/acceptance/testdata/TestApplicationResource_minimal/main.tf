variable "hrid" {
  type = string
}

resource "apim_application" "test" {
  hrid        = var.hrid
  name        = "terraform example"
  description = "Conformance tests application minimal atttributes"
}
