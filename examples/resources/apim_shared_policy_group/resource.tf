resource "apim_shared_policy_group" "example" {
  # should match the resource name
  hrid        = "example"
  name        = "[Terraform] Example headers"
  description = "Simple Shared Policy Group that contains one step to remove User-Agent header and add X-Content-Path that contains this API context path"
  api_type    = "PROXY"
  phase       = "REQUEST"
  steps = [
    {
      enabled = true
      name    = "Curate headers"
      policy  = "transform-headers"
      # Configuration is JSON as the schema depends on the policy used
      configuration = jsonencode({
        scope = "REQUEST"
        addHeaders = [
          {
            name  = "X-Context-Path"
            value = "{#request.contextPath}"
          }
        ],
        removeHeaders = ["User-Agent"]
      })
    }
  ]
}

