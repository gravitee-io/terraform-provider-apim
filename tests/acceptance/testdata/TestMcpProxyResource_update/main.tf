provider "apim" {
  organization_id = "DEFAULT"
  environment_id  = "DEFAULT"
}

variable "hrid" {
  type = string
}

variable "name" {
  type = string
}

variable "description" {
  type    = string
  default = null
}

variable "context_path" {
  type = string
}

variable "server_url" {
  type = string
}

variable "state" {
  type = string
}

variable "plans" {
  type        = list(string)
  description = "Names of the plans, sorted"
}

resource "apim_mcp_proxy" "test" {
  hrid         = var.hrid
  entity_id    = "mcp-proxy.${var.hrid}"
  name         = var.name
  description  = var.description
  context_path = var.context_path
  mode         = "PROXY"
  state        = var.state
  proxy = {
    server_url = var.server_url
  }
  plans = [
    for plan in var.plans : {
      name = plan
      # An API takes one keyless plan at most.
      security = {
        api_key = {
          source = "HEADER"
        }
      }
    }
  ]
}
