// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

// HTTPMethod - The method of the selector
type HTTPMethod string

const (
	HTTPMethodConnect HTTPMethod = "CONNECT"
	HTTPMethodDelete  HTTPMethod = "DELETE"
	HTTPMethodGet     HTTPMethod = "GET"
	HTTPMethodHead    HTTPMethod = "HEAD"
	HTTPMethodOptions HTTPMethod = "OPTIONS"
	HTTPMethodPatch   HTTPMethod = "PATCH"
	HTTPMethodPost    HTTPMethod = "POST"
	HTTPMethodPut     HTTPMethod = "PUT"
	HTTPMethodTrace   HTTPMethod = "TRACE"
	HTTPMethodOther   HTTPMethod = "OTHER"
)

func (e HTTPMethod) ToPointer() *HTTPMethod {
	return &e
}
func (e *HTTPMethod) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "CONNECT":
		fallthrough
	case "DELETE":
		fallthrough
	case "GET":
		fallthrough
	case "HEAD":
		fallthrough
	case "OPTIONS":
		fallthrough
	case "PATCH":
		fallthrough
	case "POST":
		fallthrough
	case "PUT":
		fallthrough
	case "TRACE":
		fallthrough
	case "OTHER":
		*e = HTTPMethod(v)
		return nil
	default:
		return fmt.Errorf("invalid value for HTTPMethod: %v", v)
	}
}
