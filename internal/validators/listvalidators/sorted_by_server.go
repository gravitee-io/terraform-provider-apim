package listvalidators

import "github.com/hashicorp/terraform-plugin-framework/schema/validator"

// SortedByServer returns a validator for the lists the platform reports sorted
// by catalog server: the upstream authentications of an MCP studio.
func SortedByServer() validator.List {
	return sortedBy{attributes: []string{"server"}}
}
