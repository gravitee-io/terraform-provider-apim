// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// PlanValidation - Plan validation type.
type PlanValidation string

const (
	PlanValidationAuto   PlanValidation = "AUTO"
	PlanValidationManual PlanValidation = "MANUAL"
)

func (e PlanValidation) ToPointer() *PlanValidation {
	return &e
}
func (e *PlanValidation) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "AUTO":
		fallthrough
	case "MANUAL":
		*e = PlanValidation(v)
		return nil
	default:
		return fmt.Errorf("invalid value for PlanValidation: %v", v)
	}
}
