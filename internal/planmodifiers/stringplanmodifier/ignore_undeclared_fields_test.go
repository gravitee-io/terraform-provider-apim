package stringplanmodifier_test

import (
	"context"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type testStep struct {
	Policy        types.String `tfsdk:"policy"`
	Configuration types.String `tfsdk:"configuration"`
}

type testFlow struct {
	Request []testStep `tfsdk:"request"`
}

type testEndpointGroup struct {
	Type                types.String `tfsdk:"type"`
	SharedConfiguration types.String `tfsdk:"shared_configuration"`
}

type testResource struct {
	Configuration  types.String        `tfsdk:"configuration"`
	Flows          []testFlow          `tfsdk:"flows"`
	EndpointGroups []testEndpointGroup `tfsdk:"endpoint_groups"`
}

func listOf(attributes map[string]schema.Attribute) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: attributes},
	}
}

var testSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"configuration": schema.StringAttribute{Optional: true},
		"flows": listOf(map[string]schema.Attribute{
			"request": listOf(map[string]schema.Attribute{
				"policy":        schema.StringAttribute{Optional: true},
				"configuration": schema.StringAttribute{Optional: true},
			}),
		}),
		"endpoint_groups": listOf(map[string]schema.Attribute{
			"type":                 schema.StringAttribute{Optional: true},
			"shared_configuration": schema.StringAttribute{Optional: true},
		}),
	},
}

var (
	rootConfiguration  = path.Root("configuration")
	firstStep          = path.Root("flows").AtListIndex(0).AtName("request").AtListIndex(0).AtName("configuration")
	firstEndpointGroup = path.Root("endpoint_groups").AtListIndex(0).AtName("shared_configuration")
)

// plan runs the modifier on the attribute at path, with config and state as the resource's configuration and state,
// and returns the planned value with the configured and the stored ones.
func plan(t *testing.T, at path.Path, config, state testResource) (planned, configured, stored types.String) {
	t.Helper()
	ctx := context.Background()

	configData := tfsdk.State{Schema: testSchema}
	if diags := configData.Set(ctx, config); diags.HasError() {
		t.Fatalf("building the configuration: %s", diags.Errors())
	}
	stateData := tfsdk.State{Schema: testSchema}
	if diags := stateData.Set(ctx, state); diags.HasError() {
		t.Fatalf("building the state: %s", diags.Errors())
	}
	if diags := configData.GetAttribute(ctx, at, &configured); diags.HasError() {
		t.Fatalf("reading the configuration at %s: %s", at, diags.Errors())
	}
	if diags := stateData.GetAttribute(ctx, at, &stored); diags.HasError() {
		t.Fatalf("reading the state at %s: %s", at, diags.Errors())
	}

	resp := &planmodifier.StringResponse{PlanValue: configured}
	stringplanmodifier.IgnoreUndeclaredFields().PlanModifyString(ctx, planmodifier.StringRequest{
		Path:        at,
		Config:      tfsdk.Config{Schema: testSchema, Raw: configData.Raw},
		State:       stateData,
		ConfigValue: configured,
		PlanValue:   configured,
		StateValue:  stored,
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	return resp.PlanValue, configured, stored
}

func assertPlanned(t *testing.T, planned, configured, stored types.String, keepsState bool) {
	t.Helper()
	expected := configured
	if keepsState {
		expected = stored
	}
	if !planned.Equal(expected) {
		t.Fatalf("expected plan value %s, got %s", expected, planned)
	}
}

func TestIgnoreUndeclaredFields_JSON(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config     types.String
		state      types.String
		keepsState bool
	}{
		"same fields in another order - keeps the state": {
			config:     types.StringValue(`{"a":1,"b":2}`),
			state:      types.StringValue(`{"b":2,"a":1}`),
			keepsState: true,
		},
		"APIM added a default - keeps the state": {
			config:     types.StringValue(`{"http":{"readTimeout":5000}}`),
			state:      types.StringValue(`{"http":{"readTimeout":5000,"maxWaitQueueSize":-1}}`),
			keepsState: true,
		},
		"declared field changed - plans the configuration": {
			config: types.StringValue(`{"http":{"readTimeout":5000}}`),
			state:  types.StringValue(`{"http":{"readTimeout":1000}}`),
		},
		"APIM added a default inside a list item - keeps the state": {
			config:     types.StringValue(`{"items":[{"name":"a"}]}`),
			state:      types.StringValue(`{"items":[{"name":"a","enabled":true}]}`),
			keepsState: true,
		},
		"list item removed - plans the configuration": {
			config: types.StringValue(`{"items":[{"name":"a"}]}`),
			state:  types.StringValue(`{"items":[{"name":"a"},{"name":"b"}]}`),
		},
		"same number written differently - keeps the state": {
			config:     types.StringValue(`{"limit":30000}`),
			state:      types.StringValue(`{"limit":30000.0}`),
			keepsState: true,
		},
		"declared null where APIM has a value - plans the configuration": {
			config: types.StringValue(`{"a":null}`),
			state:  types.StringValue(`{"a":1}`),
		},
		"declared field removed - keeps the state, APIM keeps the last value": {
			config:     types.StringValue(`{"http":{}}`),
			state:      types.StringValue(`{"http":{"readTimeout":5000}}`),
			keepsState: true,
		},
		"no state (creation) - plans the configuration": {
			config: types.StringValue(`{"a":1}`),
			state:  types.StringNull(),
		},
		"attribute left out - plans null": {
			config: types.StringNull(),
			state:  types.StringValue(`{"a":1}`),
		},
		"configuration is not JSON - plans the configuration": {
			config: types.StringValue(`not json`),
			state:  types.StringValue(`{"a":1}`),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			planned, configured, stored := plan(t, rootConfiguration,
				testResource{Configuration: tc.config},
				testResource{Configuration: tc.state},
			)
			assertPlanned(t, planned, configured, stored, tc.keepsState)
		})
	}
}

func TestIgnoreUndeclaredFields_ListItemGuards(t *testing.T) {
	t.Parallel()

	const (
		quotaConfigured = `{"limit":10}`
		quotaStored     = `{"limit":10,"periodTime":1,"periodTimeUnit":"MINUTES","addHeaders":true,"async":false}`
		addHeaders      = `{"addHeaders":[{"name":"x-team","value":"devex"}],"scope":"REQUEST"}`
		httpConfigured  = `{"http":{"readTimeout":5000}}`
		httpStored      = `{"http":{"readTimeout":5000,"connectTimeout":3000,"keepAlive":true}}`
	)
	step := func(policy, configuration string) testStep {
		return testStep{Policy: types.StringValue(policy), Configuration: types.StringValue(configuration)}
	}
	flow := func(steps ...testStep) testFlow {
		return testFlow{Request: steps}
	}
	endpointGroup := func(kind, configuration string) testEndpointGroup {
		return testEndpointGroup{Type: types.StringValue(kind), SharedConfiguration: types.StringValue(configuration)}
	}

	testCases := map[string]struct {
		at         path.Path
		config     testResource
		state      testResource
		keepsState bool
	}{
		"same step, configuration contained in the stored one - keeps the state": {
			at:         firstStep,
			config:     testResource{Flows: []testFlow{flow(step("quota", quotaConfigured))}},
			state:      testResource{Flows: []testFlow{flow(step("quota", quotaStored))}},
			keepsState: true,
		},
		"step inserted before the quota step - plans the configuration": {
			at: firstStep,
			config: testResource{Flows: []testFlow{flow(
				step("transform-headers", `{}`),
				step("quota", quotaConfigured),
			)}},
			state: testResource{Flows: []testFlow{flow(step("quota", quotaStored))}},
		},
		"step inserted before a step of the same policy - plans the configuration": {
			at: firstStep,
			config: testResource{Flows: []testFlow{flow(
				step("transform-headers", `{}`),
				step("transform-headers", addHeaders),
			)}},
			state: testResource{Flows: []testFlow{flow(step("transform-headers", addHeaders))}},
		},
		"flow inserted before the flow holding the step - plans the configuration": {
			at: firstStep,
			config: testResource{Flows: []testFlow{
				flow(step("quota", quotaConfigured)),
				flow(step("quota", quotaConfigured)),
			}},
			state: testResource{Flows: []testFlow{flow(step("quota", quotaStored))}},
		},
		"step replaced by another policy at the same position - plans the configuration": {
			at:     firstStep,
			config: testResource{Flows: []testFlow{flow(step("rate-limit", `{}`))}},
			state:  testResource{Flows: []testFlow{flow(step("quota", quotaStored))}},
		},
		"endpoint group type changed at the same position - plans the configuration": {
			at:     firstEndpointGroup,
			config: testResource{EndpointGroups: []testEndpointGroup{endpointGroup("mock", httpConfigured)}},
			state:  testResource{EndpointGroups: []testEndpointGroup{endpointGroup("http-proxy", httpStored)}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			planned, configured, stored := plan(t, tc.at, tc.config, tc.state)
			assertPlanned(t, planned, configured, stored, tc.keepsState)
		})
	}
}
