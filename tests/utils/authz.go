package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// AuthzKind names a collection of the AuthZ automation API.
type AuthzKind string

const (
	AuthzSchemas  AuthzKind = "schemas"
	AuthzPdps     AuthzKind = "pdps"
	AuthzEntities AuthzKind = "entities"
	AuthzPolicies AuthzKind = "policies"
)

// SkipUnlessAuthzServes skips the test unless the APIM under test serves each
// kind. It reads an HRID that exists nowhere: a 404 carrying the technical code
// hrid.notFound comes from the AuthZ module. Another 404 means the module, or
// that kind, is not there: APIM 4.13 and older, or a master that bundles an
// earlier module. Any other answer fails the test, a 403 included: it means the
// licence lacks the authz feature, which a skip would hide.
func SkipUnlessAuthzServes(t *testing.T, kinds ...AuthzKind) {
	t.Helper()
	for _, kind := range kinds {
		served, err := authzServes(kind)
		if err != nil {
			t.Fatalf("probing the authz %s: %v", kind, err)
		}
		if !served {
			t.Skipf("Skipping test: the APIM under test does not serve the authz %s", kind)
		}
	}
}

func authzServes(kind AuthzKind) (bool, error) {
	server := strings.TrimSuffix(os.Getenv("APIM_SERVER_URL"), "/")
	if server == "" {
		return false, fmt.Errorf("APIM_SERVER_URL is not set")
	}
	url := fmt.Sprintf("%s/organizations/DEFAULT/environments/DEFAULT/authz/%s/probe-%s", server, kind, acctest.RandString(10))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	if token := os.Getenv("APIM_SA_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		req.SetBasicAuth(os.Getenv("APIM_USERNAME"), os.Getenv("APIM_PASSWORD"))
	}

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusNotFound {
		return false, fmt.Errorf("GET %s answered %d: %s", url, resp.StatusCode, body)
	}
	var answer struct {
		TechnicalCode string `json:"technicalCode"`
	}
	_ = json.Unmarshal(body, &answer)
	return answer.TechnicalCode == "hrid.notFound", nil
}
