package listvalidators

import "github.com/hashicorp/terraform-plugin-framework/schema/validator"

// SortedByServerAndTool returns a validator for the lists the platform reports
// sorted by catalog server, then by tool: the tools of an MCP studio.
func SortedByServerAndTool() validator.List {
	return sortedBy{attributes: []string{"server", "tool"}}
}
