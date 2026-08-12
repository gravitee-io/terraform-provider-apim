variable "environment_id" {
  type = string
}

variable "hrid" {
  type = string
}

variable "organization_id" {
  type = string
}

resource "apim_shared_policy_group" "test" {
  api_type        = "PROXY"
  environment_id  = var.environment_id
  hrid            = var.hrid
  name            = "terraform_example"
  organization_id = var.organization_id
  phase           = "REQUEST"
  steps = [
    {
      enabled     = true
      description = "Simple Groovy script"
      name        = "Simple Groovy script"
      policy      = "groovy"
      configuration = jsonencode({
        script                 = <<-EOT
          return 'test'"
        EOT
        readContent            = true
        overrideContent        = false
        strictExecutionTimeout = false
      }),
    },
  ]
}
