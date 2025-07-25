// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/internal/utils"
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/shared"
	"net/http"
)

type GetAPIGlobals struct {
	// Id of an organization.
	OrganizationID *string `default:"DEFAULT" pathParam:"style=simple,explode=false,name=orgId"`
	// Id of an environment.
	EnvironmentID *string `default:"DEFAULT" pathParam:"style=simple,explode=false,name=envId"`
}

func (g GetAPIGlobals) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GetAPIGlobals) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *GetAPIGlobals) GetOrganizationID() *string {
	if o == nil {
		return nil
	}
	return o.OrganizationID
}

func (o *GetAPIGlobals) GetEnvironmentID() *string {
	if o == nil {
		return nil
	}
	return o.EnvironmentID
}

type GetAPIRequest struct {
	// Id of an organization.
	OrganizationID *string `default:"DEFAULT" pathParam:"style=simple,explode=false,name=orgId"`
	// Id of an environment.
	EnvironmentID *string `default:"DEFAULT" pathParam:"style=simple,explode=false,name=envId"`
	// Human-readable ID of a spec
	Hrid string `pathParam:"style=simple,explode=false,name=hrid"`
}

func (g GetAPIRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GetAPIRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *GetAPIRequest) GetOrganizationID() *string {
	if o == nil {
		return nil
	}
	return o.OrganizationID
}

func (o *GetAPIRequest) GetEnvironmentID() *string {
	if o == nil {
		return nil
	}
	return o.EnvironmentID
}

func (o *GetAPIRequest) GetHrid() string {
	if o == nil {
		return ""
	}
	return o.Hrid
}

type GetAPIResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful retrieval
	APIV4State *shared.APIV4State
}

func (o *GetAPIResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetAPIResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetAPIResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetAPIResponse) GetAPIV4State() *shared.APIV4State {
	if o == nil {
		return nil
	}
	return o.APIV4State
}
