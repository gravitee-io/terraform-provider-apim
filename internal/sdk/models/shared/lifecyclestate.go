// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// LifecycleState - The state of the API regarding the gateway(s).
type LifecycleState string

const (
	LifecycleStateClosed      LifecycleState = "CLOSED"
	LifecycleStateInitialized LifecycleState = "INITIALIZED"
	LifecycleStateStarted     LifecycleState = "STARTED"
	LifecycleStateStopped     LifecycleState = "STOPPED"
	LifecycleStateStopping    LifecycleState = "STOPPING"
)

func (e LifecycleState) ToPointer() *LifecycleState {
	return &e
}
func (e *LifecycleState) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "CLOSED":
		fallthrough
	case "INITIALIZED":
		fallthrough
	case "STARTED":
		fallthrough
	case "STOPPED":
		fallthrough
	case "STOPPING":
		*e = LifecycleState(v)
		return nil
	default:
		return fmt.Errorf("invalid value for LifecycleState: %v", v)
	}
}
