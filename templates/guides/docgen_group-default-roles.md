---
page_title: "Group with default member roles"
subcategory: "Group"
---

# Group with default member roles

This example declares a group whose members come from an identity provider.
When an OIDC or SAML login maps a user into the group, the user receives the roles in `default_member_roles`:
`USER` on APIs, applications and API products here.

`ignore_members = true` stops Terraform from removing those mapped members at each apply.
Without it, an apply that declares no `members` removes every member of the group.

`default_member_roles` sets `api`, `application` and `api_product`. Declared, it is the whole set: a scope left out loses its default role.
When the attribute is not set, the group keeps the default roles set in the Console.

```terraform
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

```
