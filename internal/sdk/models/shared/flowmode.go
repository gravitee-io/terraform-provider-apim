// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// FlowMode - API's flow mode.
type FlowMode string

const (
	FlowModeBestMatch FlowMode = "BEST_MATCH"
	FlowModeDefault   FlowMode = "DEFAULT"
)

func (e FlowMode) ToPointer() *FlowMode {
	return &e
}
func (e *FlowMode) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "BEST_MATCH":
		fallthrough
	case "DEFAULT":
		*e = FlowMode(v)
		return nil
	default:
		return fmt.Errorf("invalid value for FlowMode: %v", v)
	}
}
