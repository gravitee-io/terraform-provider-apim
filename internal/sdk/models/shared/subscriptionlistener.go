// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type SubscriptionListener struct {
	// Listener type.
	Type        ListenerType `json:"type"`
	Entrypoints []Entrypoint `json:"entrypoints,omitempty"`
	Servers     []string     `json:"servers,omitempty"`
}

func (o *SubscriptionListener) GetType() ListenerType {
	if o == nil {
		return ListenerType("")
	}
	return o.Type
}

func (o *SubscriptionListener) GetEntrypoints() []Entrypoint {
	if o == nil {
		return nil
	}
	return o.Entrypoints
}

func (o *SubscriptionListener) GetServers() []string {
	if o == nil {
		return nil
	}
	return o.Servers
}
