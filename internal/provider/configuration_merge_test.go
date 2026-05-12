package provider

import (
	"encoding/json"
	"testing"

	tfTypes "github.com/gravitee-io/terraform-provider-apim/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMergePreserveState_PreservesWriteOnlyField(t *testing.T) {
	apiJSON := []byte(`{"cacheKey":"zoho-access-token","defaultOperation":"SET","resource":"cache-global","value":"{#context.attributes['zoho_access_token']}"}`)
	stateJSON := []byte(`{"cacheKey":"zoho-access-token","defaultOperation":"SET","resource":"cache-global","timeToLiveSeconds":3500,"value":"{#context.attributes['zoho_access_token']}"}`)

	result, err := mergePreserveState(apiJSON, stateJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(result, &m); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	raw, ok := m["timeToLiveSeconds"]
	if !ok {
		t.Fatal("expected timeToLiveSeconds to be preserved from state but it is missing")
	}
	if string(raw) != "3500" {
		t.Fatalf("expected timeToLiveSeconds=3500, got %s", raw)
	}
}

func TestMergePreserveState_APIValueWins(t *testing.T) {
	apiJSON := []byte(`{"cacheKey":"new-key"}`)
	stateJSON := []byte(`{"cacheKey":"old-key","timeToLiveSeconds":3500}`)

	result, err := mergePreserveState(apiJSON, stateJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(result, &m); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	if string(m["cacheKey"]) != `"new-key"` {
		t.Fatalf("expected API value to win: cacheKey=%s", m["cacheKey"])
	}
	if string(m["timeToLiveSeconds"]) != "3500" {
		t.Fatalf("expected timeToLiveSeconds preserved from state, got %s", m["timeToLiveSeconds"])
	}
}

func TestMergePreserveState_NilStateReturnsAPI(t *testing.T) {
	apiJSON := []byte(`{"cacheKey":"my-key"}`)

	result, err := mergePreserveState(apiJSON, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != string(apiJSON) {
		t.Fatalf("expected api JSON unchanged, got %s", result)
	}
}

func TestMergePreserveState_NonObjectAPIJSON(t *testing.T) {
	apiJSON := []byte(`"just-a-string"`)
	stateJSON := []byte(`{"key":"value"}`)

	result, err := mergePreserveState(apiJSON, stateJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != string(apiJSON) {
		t.Fatalf("expected non-object API JSON returned unchanged, got %s", result)
	}
}

// TestPreserveWriteOnlyConfigs_PlanFlow verifies that after a Read in which the API
// omits a write-only field (e.g. timeToLiveSeconds), the prior state value for that
// field is preserved in r.Plans flows.
func TestPreserveWriteOnlyConfigs_PlanFlow(t *testing.T) {
	priorConfig := `{"cacheKey":"my-key","defaultOperation":"SET","timeToLiveSeconds":3500}`
	afterReadConfig := `{"cacheKey":"my-key","defaultOperation":"SET"}` // API dropped timeToLiveSeconds

	r := &Apiv4ResourceModel{
		Plans: []tfTypes.PlanV4{
			{
				Flows: []tfTypes.FlowV4{
					{
						Request: []tfTypes.StepV4{
							{
								Policy:        types.StringValue("data-cache"),
								Configuration: jsontypes.NewNormalizedValue(afterReadConfig),
							},
						},
					},
				},
			},
		},
	}

	priorPlans := []tfTypes.PlanV4{
		{
			Flows: []tfTypes.FlowV4{
				{
					Request: []tfTypes.StepV4{
						{
							Policy:        types.StringValue("data-cache"),
							Configuration: jsontypes.NewNormalizedValue(priorConfig),
						},
					},
				},
			},
		},
	}

	r.preserveWriteOnlyConfigs(priorPlans, nil)

	gotConfig := r.Plans[0].Flows[0].Request[0].Configuration.ValueString()
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(gotConfig), &m); err != nil {
		t.Fatalf("configuration is not valid JSON: %v", err)
	}
	if string(m["timeToLiveSeconds"]) != "3500" {
		t.Fatalf("expected timeToLiveSeconds=3500 to be preserved, got config: %s", gotConfig)
	}
}

// TestPreserveWriteOnlyConfigs_TopLevelFlow verifies the same for top-level flows.
func TestPreserveWriteOnlyConfigs_TopLevelFlow(t *testing.T) {
	priorConfig := `{"cacheKey":"my-key","timeToLiveSeconds":3600}`
	afterReadConfig := `{"cacheKey":"my-key"}` // API dropped timeToLiveSeconds

	r := &Apiv4ResourceModel{
		Flows: []tfTypes.FlowV4{
			{
				Request: []tfTypes.StepV4{
					{
						Policy:        types.StringValue("data-cache"),
						Configuration: jsontypes.NewNormalizedValue(afterReadConfig),
					},
				},
			},
		},
	}

	priorFlows := []tfTypes.FlowV4{
		{
			Request: []tfTypes.StepV4{
				{
					Policy:        types.StringValue("data-cache"),
					Configuration: jsontypes.NewNormalizedValue(priorConfig),
				},
			},
		},
	}

	r.preserveWriteOnlyConfigs(nil, priorFlows)

	gotConfig := r.Flows[0].Request[0].Configuration.ValueString()
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(gotConfig), &m); err != nil {
		t.Fatalf("configuration is not valid JSON: %v", err)
	}
	if string(m["timeToLiveSeconds"]) != "3600" {
		t.Fatalf("expected timeToLiveSeconds=3600 to be preserved, got config: %s", gotConfig)
	}
}

// TestPreserveWriteOnlyConfigs_NewStepNoStateToCopy verifies that new steps added
// after the last state (no prior state entry) are left unchanged.
func TestPreserveWriteOnlyConfigs_NewStepNoStateToCopy(t *testing.T) {
	afterReadConfig := `{"cacheKey":"brand-new"}`

	r := &Apiv4ResourceModel{
		Plans: []tfTypes.PlanV4{
			{
				Flows: []tfTypes.FlowV4{
					{
						Request: []tfTypes.StepV4{
							{Configuration: jsontypes.NewNormalizedValue(afterReadConfig)},
						},
					},
				},
			},
		},
	}

	// No prior plans at all — simulate first-time creation.
	r.preserveWriteOnlyConfigs(nil, nil)

	gotConfig := r.Plans[0].Flows[0].Request[0].Configuration.ValueString()
	if gotConfig != afterReadConfig {
		t.Fatalf("expected config unchanged, got %s", gotConfig)
	}
}
