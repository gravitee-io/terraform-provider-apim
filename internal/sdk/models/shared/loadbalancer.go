// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/internal/utils"
)

// LoadBalancerType - Load balancer type.
type LoadBalancerType string

const (
	LoadBalancerTypeRandom             LoadBalancerType = "RANDOM"
	LoadBalancerTypeRoundRobin         LoadBalancerType = "ROUND_ROBIN"
	LoadBalancerTypeWeightedRandom     LoadBalancerType = "WEIGHTED_RANDOM"
	LoadBalancerTypeWeightedRoundRobin LoadBalancerType = "WEIGHTED_ROUND_ROBIN"
)

func (e LoadBalancerType) ToPointer() *LoadBalancerType {
	return &e
}
func (e *LoadBalancerType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "RANDOM":
		fallthrough
	case "ROUND_ROBIN":
		fallthrough
	case "WEIGHTED_RANDOM":
		fallthrough
	case "WEIGHTED_ROUND_ROBIN":
		*e = LoadBalancerType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for LoadBalancerType: %v", v)
	}
}

type LoadBalancer struct {
	// Load balancer type.
	Type *LoadBalancerType `default:"ROUND_ROBIN" json:"type"`
}

func (l LoadBalancer) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(l, "", false)
}

func (l *LoadBalancer) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &l, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *LoadBalancer) GetType() *LoadBalancerType {
	if o == nil {
		return nil
	}
	return o.Type
}
