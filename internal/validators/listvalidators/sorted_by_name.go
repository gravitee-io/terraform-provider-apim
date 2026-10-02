package listvalidators

import "github.com/hashicorp/terraform-plugin-framework/schema/validator"

// SortedByName returns a validator for the lists the platform reports sorted by
// name: the plans and the identity providers of an MCP proxy.
func SortedByName() validator.List {
	return sortedBy{attributes: []string{"name"}}
}
