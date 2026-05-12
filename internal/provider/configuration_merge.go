package provider

import (
	"encoding/json"

	tfTypes "github.com/gravitee-io/terraform-provider-apim/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

// mergePreserveState returns apiJSON with any top-level keys from stateJSON added
// back if they are absent in apiJSON. Keys present in apiJSON are never overwritten —
// this only fills gaps caused by the Gravitee API not round-tripping write-only fields
// such as data-cache's timeToLiveSeconds.
func mergePreserveState(apiJSON, stateJSON []byte) ([]byte, error) {
	var apiMap map[string]json.RawMessage
	if err := json.Unmarshal(apiJSON, &apiMap); err != nil {
		return apiJSON, nil
	}
	var stateMap map[string]json.RawMessage
	if err := json.Unmarshal(stateJSON, &stateMap); err != nil {
		return apiJSON, nil
	}
	for k, v := range stateMap {
		if _, exists := apiMap[k]; !exists {
			apiMap[k] = v
		}
	}
	return json.Marshal(apiMap)
}

// mergeStepConfig merges write-only fields from priorConfig into current and returns the result.
func mergeStepConfig(current jsontypes.Normalized, prior jsontypes.Normalized) jsontypes.Normalized {
	if current.IsNull() || current.IsUnknown() {
		return current
	}
	if prior.IsNull() || prior.IsUnknown() {
		return current
	}
	merged, err := mergePreserveState([]byte(current.ValueString()), []byte(prior.ValueString()))
	if err != nil {
		return current
	}
	return jsontypes.NewNormalizedValue(string(merged))
}

// mergeFlowSteps merges write-only configuration fields from priorSteps into steps by index.
func mergeFlowSteps(steps []tfTypes.StepV4, priorSteps []tfTypes.StepV4) {
	for i := range steps {
		if i >= len(priorSteps) {
			break
		}
		steps[i].Configuration = mergeStepConfig(steps[i].Configuration, priorSteps[i].Configuration)
	}
}

// mergeFlowStepConfigs merges write-only configuration fields from priorFlow into flow for all step types.
func mergeFlowStepConfigs(flow *tfTypes.FlowV4, priorFlow *tfTypes.FlowV4) {
	mergeFlowSteps(flow.EntrypointConnect, priorFlow.EntrypointConnect)
	mergeFlowSteps(flow.Interact, priorFlow.Interact)
	mergeFlowSteps(flow.Publish, priorFlow.Publish)
	mergeFlowSteps(flow.Request, priorFlow.Request)
	mergeFlowSteps(flow.Response, priorFlow.Response)
	mergeFlowSteps(flow.Subscribe, priorFlow.Subscribe)
}

// preserveWriteOnlyConfigs re-adds write-only policy configuration fields that the
// Gravitee API accepts on write but omits from GET responses (e.g. data-cache's
// timeToLiveSeconds). It must be called after RefreshFromSharedApiv4State with the
// Plans and Flows values captured before that call.
func (r *Apiv4ResourceModel) preserveWriteOnlyConfigs(priorPlans []tfTypes.PlanV4, priorFlows []tfTypes.FlowV4) {
	for planIdx := range r.Plans {
		if planIdx >= len(priorPlans) {
			break
		}
		for flowIdx := range r.Plans[planIdx].Flows {
			if flowIdx >= len(priorPlans[planIdx].Flows) {
				break
			}
			mergeFlowStepConfigs(&r.Plans[planIdx].Flows[flowIdx], &priorPlans[planIdx].Flows[flowIdx])
		}
	}
	for flowIdx := range r.Flows {
		if flowIdx >= len(priorFlows) {
			break
		}
		mergeFlowStepConfigs(&r.Flows[flowIdx], &priorFlows[flowIdx])
	}
}
