// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type TCPListener struct {
	// Listener type.
	Type        ListenerType `json:"type"`
	Entrypoints []Entrypoint `json:"entrypoints,omitempty"`
	Servers     []string     `json:"servers,omitempty"`
	// A list of hostnames for which the API will match against SNI.  This must be unique for all TCP listener for a given server id. See 'servers' attribute
	Hosts []string `json:"hosts"`
}

func (o *TCPListener) GetType() ListenerType {
	if o == nil {
		return ListenerType("")
	}
	return o.Type
}

func (o *TCPListener) GetEntrypoints() []Entrypoint {
	if o == nil {
		return nil
	}
	return o.Entrypoints
}

func (o *TCPListener) GetServers() []string {
	if o == nil {
		return nil
	}
	return o.Servers
}

func (o *TCPListener) GetHosts() []string {
	if o == nil {
		return []string{}
	}
	return o.Hosts
}
