package listvalidators

import (
	"context"
	"fmt"
	"strings"

	"github.com/gravitee-io/terraform-provider-apim/internal/listkeys"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.List = sortedBy{}

// sortedBy refuses a list of objects declared in another order than the one the
// platform reports it in. The provider keeps the declared order in the state
// after an apply and takes the platform's on the next read: left unchecked, the
// difference shows as a change on every plan.
type sortedBy struct {
	attributes []string
}

func (v sortedBy) Description(_ context.Context) string {
	return "elements must be sorted by " + strings.Join(v.attributes, " then ")
}

func (v sortedBy) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v sortedBy) ValidateList(_ context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	keys, known := listkeys.Of(req.ConfigValue, v.attributes...)
	if !known || listkeys.IsSorted(keys) {
		return
	}

	expected := make([]string, 0, len(keys))
	for _, key := range listkeys.Sorted(keys) {
		expected = append(expected, strings.Join(key, "/"))
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"List declared in another order than the platform reports",
		fmt.Sprintf(
			"%s must be sorted by %s. Declare it in this order: %s.",
			req.Path, strings.Join(v.attributes, " then "), strings.Join(expected, ", "),
		),
	)
}
