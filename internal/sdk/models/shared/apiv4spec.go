// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/internal/utils"
)

// APIV4Spec - ApiV4DefinitionSpec defines the desired state of ApiDefinition.
type APIV4Spec struct {
	// A unique human readable id identifying this object
	Hrid string `json:"hrid"`
	// the context where the api definition was created.
	//
	// Deprecated: This will be removed in a future release, please migrate away from it as soon as possible.
	DefinitionContext DefinitionContext `json:"definitionContext"`
	// API's name. Duplicate names can exists.
	Name string `json:"name"`
	// API's version. It's a simple string only used in the portal.
	Version string `json:"version"`
	// API's type.
	Type APIType `json:"type"`
	// API's description. A short description of your API.
	Description *string `json:"description,omitempty"`
	// The list of sharding tags associated with this API.
	Tags []string `json:"tags,omitempty"`
	// The list of listeners associated with this API.
	Listeners      []Listener        `json:"listeners"`
	EndpointGroups []EndpointGroupV4 `json:"endpointGroups"`
	Analytics      *Analytics        `json:"analytics,omitempty"`
	Failover       *FailoverV4       `json:"failover,omitempty"`
	Properties     []Property        `json:"properties,omitempty"`
	Resources      []Resource        `json:"resources,omitempty"`
	// Map of plan IDs to Plan objects
	Plans         map[string]PlanInput `json:"plans,omitempty"`
	FlowExecution *FlowExecution       `json:"flowExecution,omitempty"`
	// List of flows for the API
	Flows []FlowV4Input `json:"flows,omitempty"`
	// A list of Response Templates for the API (Not applicable for Native API)
	ResponseTemplates map[string]map[string]ResponseTemplate `json:"responseTemplates,omitempty"`
	Services          *APIServices                           `json:"services,omitempty"`
	// List of groups associated with the API.
	// This groups are id or name references to existing groups in APIM.
	Groups []string `json:"groups,omitempty"`
	// The visibility of the resource regarding the portal.
	Visibility *Visibility `json:"visibility,omitempty"`
	// The state of the API regarding the gateway(s).
	State        *LifecycleState `json:"state,omitempty"`
	PrimaryOwner *PrimaryOwner   `json:"primaryOwner,omitempty"`
	// List of labels of the API
	Labels []string `json:"labels,omitempty"`
	// The list of API's metadata.
	Metadata []Metadata `json:"metadata,omitempty"`
	// The status of the API regarding the console.
	LifecycleState APILifecycleState `json:"lifecycleState"`
	// The list of category keys associated with this API.
	Categories []string `json:"categories,omitempty"`
	// Set of members associated with the plan
	Members []Member `json:"members,omitempty"`
	// If true, new members added to the API spec will
	// be notified when the API is synced with APIM.
	NotifyMembers *bool `default:"true" json:"notifyMembers"`
}

func (a APIV4Spec) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(a, "", false)
}

func (a *APIV4Spec) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &a, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *APIV4Spec) GetHrid() string {
	if o == nil {
		return ""
	}
	return o.Hrid
}

func (o *APIV4Spec) GetDefinitionContext() DefinitionContext {
	if o == nil {
		return DefinitionContext{}
	}
	return o.DefinitionContext
}

func (o *APIV4Spec) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *APIV4Spec) GetVersion() string {
	if o == nil {
		return ""
	}
	return o.Version
}

func (o *APIV4Spec) GetType() APIType {
	if o == nil {
		return APIType("")
	}
	return o.Type
}

func (o *APIV4Spec) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *APIV4Spec) GetTags() []string {
	if o == nil {
		return nil
	}
	return o.Tags
}

func (o *APIV4Spec) GetListeners() []Listener {
	if o == nil {
		return []Listener{}
	}
	return o.Listeners
}

func (o *APIV4Spec) GetEndpointGroups() []EndpointGroupV4 {
	if o == nil {
		return []EndpointGroupV4{}
	}
	return o.EndpointGroups
}

func (o *APIV4Spec) GetAnalytics() *Analytics {
	if o == nil {
		return nil
	}
	return o.Analytics
}

func (o *APIV4Spec) GetFailover() *FailoverV4 {
	if o == nil {
		return nil
	}
	return o.Failover
}

func (o *APIV4Spec) GetProperties() []Property {
	if o == nil {
		return nil
	}
	return o.Properties
}

func (o *APIV4Spec) GetResources() []Resource {
	if o == nil {
		return nil
	}
	return o.Resources
}

func (o *APIV4Spec) GetPlans() map[string]PlanInput {
	if o == nil {
		return nil
	}
	return o.Plans
}

func (o *APIV4Spec) GetFlowExecution() *FlowExecution {
	if o == nil {
		return nil
	}
	return o.FlowExecution
}

func (o *APIV4Spec) GetFlows() []FlowV4Input {
	if o == nil {
		return nil
	}
	return o.Flows
}

func (o *APIV4Spec) GetResponseTemplates() map[string]map[string]ResponseTemplate {
	if o == nil {
		return nil
	}
	return o.ResponseTemplates
}

func (o *APIV4Spec) GetServices() *APIServices {
	if o == nil {
		return nil
	}
	return o.Services
}

func (o *APIV4Spec) GetGroups() []string {
	if o == nil {
		return nil
	}
	return o.Groups
}

func (o *APIV4Spec) GetVisibility() *Visibility {
	if o == nil {
		return nil
	}
	return o.Visibility
}

func (o *APIV4Spec) GetState() *LifecycleState {
	if o == nil {
		return nil
	}
	return o.State
}

func (o *APIV4Spec) GetPrimaryOwner() *PrimaryOwner {
	if o == nil {
		return nil
	}
	return o.PrimaryOwner
}

func (o *APIV4Spec) GetLabels() []string {
	if o == nil {
		return nil
	}
	return o.Labels
}

func (o *APIV4Spec) GetMetadata() []Metadata {
	if o == nil {
		return nil
	}
	return o.Metadata
}

func (o *APIV4Spec) GetLifecycleState() APILifecycleState {
	if o == nil {
		return APILifecycleState("")
	}
	return o.LifecycleState
}

func (o *APIV4Spec) GetCategories() []string {
	if o == nil {
		return nil
	}
	return o.Categories
}

func (o *APIV4Spec) GetMembers() []Member {
	if o == nil {
		return nil
	}
	return o.Members
}

func (o *APIV4Spec) GetNotifyMembers() *bool {
	if o == nil {
		return nil
	}
	return o.NotifyMembers
}
