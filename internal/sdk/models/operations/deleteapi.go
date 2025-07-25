// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/internal/utils"
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/shared"
	"net/http"
)

type DeleteAPIGlobals struct {
	// Id of an organization.
	OrganizationID *string `default:"DEFAULT" pathParam:"style=simple,explode=false,name=orgId"`
	// Id of an environment.
	EnvironmentID *string `default:"DEFAULT" pathParam:"style=simple,explode=false,name=envId"`
}

func (d DeleteAPIGlobals) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(d, "", false)
}

func (d *DeleteAPIGlobals) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &d, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *DeleteAPIGlobals) GetOrganizationID() *string {
	if o == nil {
		return nil
	}
	return o.OrganizationID
}

func (o *DeleteAPIGlobals) GetEnvironmentID() *string {
	if o == nil {
		return nil
	}
	return o.EnvironmentID
}

type DeleteAPIRequest struct {
	// Id of an organization.
	OrganizationID *string `default:"DEFAULT" pathParam:"style=simple,explode=false,name=orgId"`
	// Id of an environment.
	EnvironmentID *string `default:"DEFAULT" pathParam:"style=simple,explode=false,name=envId"`
	// Human-readable ID of a spec
	Hrid string `pathParam:"style=simple,explode=false,name=hrid"`
}

func (d DeleteAPIRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(d, "", false)
}

func (d *DeleteAPIRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &d, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *DeleteAPIRequest) GetOrganizationID() *string {
	if o == nil {
		return nil
	}
	return o.OrganizationID
}

func (o *DeleteAPIRequest) GetEnvironmentID() *string {
	if o == nil {
		return nil
	}
	return o.EnvironmentID
}

func (o *DeleteAPIRequest) GetHrid() string {
	if o == nil {
		return ""
	}
	return o.Hrid
}

type DeleteAPIResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Validation error
	Error *shared.Error
}

func (o *DeleteAPIResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *DeleteAPIResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *DeleteAPIResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *DeleteAPIResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
