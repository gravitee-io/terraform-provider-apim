---
page_title: "Environment variables override"
---

# Environment variables override

Gravitee Terraform provider can be configured using the following environment variables

| Variable          | Corresponding provider field |
| ----------------- | ---------------------------- |
| APIM\_SERVER\_URL | server\_url                  |
| APIM\_SA\_TOKEN   | bearer\_auth                 |
| APIM\_ORG\_ID     | organization\_id             |
| APIM\_ENV\_ID     | environment\_id              |

If you use only environment variables to configure your provider, then your configuration block looks like this:

```hcl
provider "apim" {
  # Hooray, nothing is hardcoded!
}
```