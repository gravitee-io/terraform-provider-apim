provider "apim" {
  organization_id = "e60e2468-c42d-45f5-aab7-0052a3d0251b"
  environment_id  = "0f3791aa-d8f8-4288-8a74-8f4ee6824b47"
}

variable "hrid" {
  type = string
}

resource "apim_application" "test" {
  hrid            = var.hrid
  name            = "terraform example"
  description     = "Conformance tests application minimal atttributes"
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}