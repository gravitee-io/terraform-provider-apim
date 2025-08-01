// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LoggingV4 struct {
	Condition        types.String      `tfsdk:"condition"`
	Content          *LoggingContentV4 `tfsdk:"content"`
	MessageCondition types.String      `tfsdk:"message_condition"`
	Mode             *LoggingModeV4    `tfsdk:"mode"`
	Phase            *LoggingPhase     `tfsdk:"phase"`
}
