// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LoggingPhase struct {
	Request  types.Bool `tfsdk:"request"`
	Response types.Bool `tfsdk:"response"`
}
