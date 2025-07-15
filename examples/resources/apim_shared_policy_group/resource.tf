resource "apim_shared_policy_group" "my_sharedpolicygroup" {
  api_type             = "MESSAGE"
  description          = "this is a shared policy group"
  environment_id       = "...my_environment_id..."
  hrid                 = "...my_hrid..."
  name                 = "My Shared Policy Group"
  organization_id      = "...my_organization_id..."
  phase                = "RESPONSE"
  prerequisite_message = "the resource cache \"my-cache\" is required"
  steps = [
    {
      condition         = "...my_condition..."
      configuration     = "...my_configuration..."
      description       = "...my_description..."
      enabled           = true
      message_condition = "...my_message_condition..."
      name              = "...my_name..."
      policy            = "...my_policy..."
    }
  ]
}