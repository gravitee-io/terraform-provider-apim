// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Resource struct {
	Configuration types.String `tfsdk:"configuration"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
}
