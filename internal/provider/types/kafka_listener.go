// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KafkaListener struct {
	Entrypoints []Entrypoint   `tfsdk:"entrypoints"`
	Host        types.String   `tfsdk:"host"`
	Port        types.Int64    `tfsdk:"port"`
	Servers     []types.String `tfsdk:"servers"`
	Type        types.String   `tfsdk:"type"`
}
