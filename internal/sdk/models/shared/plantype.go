// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// PlanType - Plan type.
type PlanType string

const (
	PlanTypeAPI     PlanType = "API"
	PlanTypeCatalog PlanType = "CATALOG"
)

func (e PlanType) ToPointer() *PlanType {
	return &e
}
func (e *PlanType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "API":
		fallthrough
	case "CATALOG":
		*e = PlanType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for PlanType: %v", v)
	}
}
