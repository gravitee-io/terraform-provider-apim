resource "apim_group" "developers" {
  hrid           = "developers"
  name           = "Developers"
  notify_members = false
  # Members come from the identity provider's group mapping, not from Terraform.
  ignore_members = true
  default_member_roles = {
    api         = "USER"
    application = "USER"
    api_product = "USER"
  }
}
