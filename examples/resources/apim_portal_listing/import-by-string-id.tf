import {
  to = apim_portal_listing.public-apis
  id = jsonencode({
    environment_id  = "DEFAULT"
    hrid            = "public-apis"
    organization_id = "DEFAULT"
    portal_hrid     = "developer-portal"
  })
}
