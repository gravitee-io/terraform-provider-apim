package stringplanmodifier_test

import (
	"context"
	"testing"

	"github.com/gravitee-io/terraform-provider-apim/internal/planmodifiers/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNormalizePath_SuppressTrailingSlashDiff(t *testing.T) {
	t.Parallel()

	plan := types.StringValue("/my-api")
	resp := &planmodifier.StringResponse{PlanValue: plan}
	stringplanmodifier.NormalizePath().PlanModifyString(context.Background(), planmodifier.StringRequest{
		PlanValue:  plan,
		StateValue: types.StringValue("/my-api/"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.PlanValue.ValueString() != "/my-api/" {
		t.Fatalf("expected plan value to be /my-api/, got %s", resp.PlanValue.ValueString())
	}
}

func TestNormalizePath_NoChangeWhenAlreadyHasSlash(t *testing.T) {
	t.Parallel()

	plan := types.StringValue("/my-api/")
	resp := &planmodifier.StringResponse{PlanValue: plan}
	stringplanmodifier.NormalizePath().PlanModifyString(context.Background(), planmodifier.StringRequest{
		PlanValue:  plan,
		StateValue: types.StringValue("/my-api/"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.PlanValue.ValueString() != "/my-api/" {
		t.Fatalf("expected plan value to be /my-api/, got %s", resp.PlanValue.ValueString())
	}
}

func TestNormalizePath_NoChangeOnRealDiff(t *testing.T) {
	t.Parallel()

	plan := types.StringValue("/new-api")
	resp := &planmodifier.StringResponse{PlanValue: plan}
	stringplanmodifier.NormalizePath().PlanModifyString(context.Background(), planmodifier.StringRequest{
		PlanValue:  plan,
		StateValue: types.StringValue("/old-api/"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.PlanValue.ValueString() != "/new-api" {
		t.Fatalf("expected plan value to remain /new-api, got %s", resp.PlanValue.ValueString())
	}
}

func TestNormalizePath_NullPlanValue(t *testing.T) {
	t.Parallel()

	plan := types.StringNull()
	resp := &planmodifier.StringResponse{PlanValue: plan}
	stringplanmodifier.NormalizePath().PlanModifyString(context.Background(), planmodifier.StringRequest{
		PlanValue:  plan,
		StateValue: types.StringValue("/my-api/"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.IsNull() {
		t.Fatal("expected plan value to remain null")
	}
}

func TestNormalizePath_NullStateValue(t *testing.T) {
	t.Parallel()

	plan := types.StringValue("/my-api")
	resp := &planmodifier.StringResponse{PlanValue: plan}
	stringplanmodifier.NormalizePath().PlanModifyString(context.Background(), planmodifier.StringRequest{
		PlanValue:  plan,
		StateValue: types.StringNull(),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.PlanValue.ValueString() != "/my-api" {
		t.Fatalf("expected plan value to remain /my-api, got %s", resp.PlanValue.ValueString())
	}
}

func TestNormalizePath_RootPath(t *testing.T) {
	t.Parallel()

	plan := types.StringValue("/")
	resp := &planmodifier.StringResponse{PlanValue: plan}
	stringplanmodifier.NormalizePath().PlanModifyString(context.Background(), planmodifier.StringRequest{
		PlanValue:  plan,
		StateValue: types.StringValue("/"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.PlanValue.ValueString() != "/" {
		t.Fatalf("expected plan value to be /, got %s", resp.PlanValue.ValueString())
	}
}
