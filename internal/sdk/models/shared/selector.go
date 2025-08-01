// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/internal/utils"
)

type SelectorType string

const (
	SelectorTypeHTTP      SelectorType = "HTTP"
	SelectorTypeChannel   SelectorType = "CHANNEL"
	SelectorTypeCondition SelectorType = "CONDITION"
)

type Selector struct {
	HTTPSelector      *HTTPSelector      `queryParam:"inline"`
	ChannelSelector   *ChannelSelector   `queryParam:"inline"`
	ConditionSelector *ConditionSelector `queryParam:"inline"`

	Type SelectorType
}

func CreateSelectorHTTP(http HTTPSelector) Selector {
	typ := SelectorTypeHTTP

	typStr := HTTPSelectorType(typ)
	http.Type = typStr

	return Selector{
		HTTPSelector: &http,
		Type:         typ,
	}
}

func CreateSelectorChannel(channel ChannelSelector) Selector {
	typ := SelectorTypeChannel

	typStr := ChannelSelectorType(typ)
	channel.Type = typStr

	return Selector{
		ChannelSelector: &channel,
		Type:            typ,
	}
}

func CreateSelectorCondition(condition ConditionSelector) Selector {
	typ := SelectorTypeCondition

	typStr := ConditionSelectorType(typ)
	condition.Type = typStr

	return Selector{
		ConditionSelector: &condition,
		Type:              typ,
	}
}

func (u *Selector) UnmarshalJSON(data []byte) error {

	type discriminator struct {
		Type string `json:"type"`
	}

	dis := new(discriminator)
	if err := json.Unmarshal(data, &dis); err != nil {
		return fmt.Errorf("could not unmarshal discriminator: %w", err)
	}

	switch dis.Type {
	case "HTTP":
		httpSelector := new(HTTPSelector)
		if err := utils.UnmarshalJSON(data, &httpSelector, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Type == HTTP) type HTTPSelector within Selector: %w", string(data), err)
		}

		u.HTTPSelector = httpSelector
		u.Type = SelectorTypeHTTP
		return nil
	case "CHANNEL":
		channelSelector := new(ChannelSelector)
		if err := utils.UnmarshalJSON(data, &channelSelector, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Type == CHANNEL) type ChannelSelector within Selector: %w", string(data), err)
		}

		u.ChannelSelector = channelSelector
		u.Type = SelectorTypeChannel
		return nil
	case "CONDITION":
		conditionSelector := new(ConditionSelector)
		if err := utils.UnmarshalJSON(data, &conditionSelector, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Type == CONDITION) type ConditionSelector within Selector: %w", string(data), err)
		}

		u.ConditionSelector = conditionSelector
		u.Type = SelectorTypeCondition
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for Selector", string(data))
}

func (u Selector) MarshalJSON() ([]byte, error) {
	if u.HTTPSelector != nil {
		return utils.MarshalJSON(u.HTTPSelector, "", true)
	}

	if u.ChannelSelector != nil {
		return utils.MarshalJSON(u.ChannelSelector, "", true)
	}

	if u.ConditionSelector != nil {
		return utils.MarshalJSON(u.ConditionSelector, "", true)
	}

	return nil, errors.New("could not marshal union type Selector: all fields are null")
}
