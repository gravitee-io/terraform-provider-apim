---
page_title: "Simple Application"
subcategory: "Application"
---

# Create an application

The following example creates a simple application.

```HCL
resource "apim_application" "simple" {
  # should match the resource name
  hrid        = "simple"
  name        = "[Terraform] Simple Application"
  description = "Demonstrate applications can be created with Terraform"
  settings = {
    app = {
      type = "test"
    }
  }
  metadata = [{
    name   = "hello"
    format = "STRING"
  }]
}
```