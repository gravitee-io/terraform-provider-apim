import {
  to = apim_authz_policy.github_create_issue
  id = jsonencode({
    environment_id  = "DEFAULT"
    hrid            = "engineering-create-issue"
    organization_id = "DEFAULT"
  })
}
