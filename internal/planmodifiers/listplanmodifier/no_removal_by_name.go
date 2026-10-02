package listplanmodifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/gravitee-io/terraform-provider-apim/internal/listkeys"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// NoRemovalByName returns a plan modifier that refuses a configuration dropping
// an element the state holds, elements being identified by their name. It is for
// the lists the platform only adds to, such as the identity providers of an MCP
// proxy: the platform keeps a provider the declaration no longer names and
// reports it on every read, so the removal would be planned again and again.
func NoRemovalByName() planmodifier.List {
	return noRemovalByName{}
}

type noRemovalByName struct{}

func (m noRemovalByName) Description(_ context.Context) string {
	return "An element cannot be removed on an existing resource."
}

func (m noRemovalByName) MarkdownDescription(_ context.Context) string {
	return "An element cannot be removed on an existing resource."
}

func (m noRemovalByName) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.ConfigValue.IsUnknown() {
		return
	}
	held, known := listkeys.Of(req.StateValue, "name")
	if !known {
		return
	}

	declared := map[string]bool{}
	if !req.ConfigValue.IsNull() {
		names, known := listkeys.Of(req.ConfigValue, "name")
		if !known {
			return
		}
		for _, name := range names {
			declared[name[0]] = true
		}
	}

	var removed []string
	for _, name := range held {
		if !declared[name[0]] {
			removed = append(removed, name[0])
		}
	}
	if len(removed) == 0 {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Element removed from a list the platform only adds to",
		fmt.Sprintf(
			"%s no longer declares %s. The platform keeps it and reports it on every read. "+
				"Declare it again, or replace the resource to remove it.",
			req.Path, strings.Join(removed, ", "),
		),
	)
}
