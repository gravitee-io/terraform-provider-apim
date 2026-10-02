package utils

import (
	"os"
	"testing"
)

// mcpServerURLEnv names the MCP server the catalog tests register. The platform
// discovers it when a catalog server is applied, so the URL is the one APIM
// reaches it at: the in-cluster service on CI, a public server on a laptop.
const mcpServerURLEnv = "APIM_MCP_STUB_URL"

// McpServerURL returns the MCP server to register in the catalog, and skips the
// test when none is given.
func McpServerURL(t *testing.T) string {
	url := os.Getenv(mcpServerURLEnv)
	if url == "" {
		t.Skip("Skipping test: " + mcpServerURLEnv + " is not set, there is no MCP server to discover")
	}
	return url
}
