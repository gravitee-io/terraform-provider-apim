package acceptance_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Automation declarations can omit visibility; the navigation tree exposes its effective value.
func checkNavigationVisibility(scopeAddress string, apiScoped bool, title, expected string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		scope, ok := state.RootModule().Resources[scopeAddress]
		if !ok || scope.Primary == nil {
			return fmt.Errorf("missing navigation scope %s", scopeAddress)
		}
		attributes := scope.Primary.Attributes
		base := os.Getenv("APIM_MANAGEMENT_SERVER_URL")
		if base == "" {
			base = strings.TrimSuffix(strings.TrimRight(os.Getenv("APIM_SERVER_URL"), "/"), "/automation") + "/management/v2"
		}
		environment := strings.TrimRight(base, "/") + "/organizations/" + url.PathEscape(attributes["organization_id"]) + "/environments/" + url.PathEscape(attributes["environment_id"])
		id := url.PathEscape(attributes["id"])
		endpoints := []string{environment + "/apis/" + id + "/portal-navigation-items?loadChildren=true"}
		if !apiScoped {
			// Recent APIM reads a portal's navigation under that portal: its environment-level
			// navigation only finds a default portal, and a portal created by automation is not one.
			endpoints = []string{
				environment + "/portals/" + id + "/navigation-items?loadChildren=true",
				environment + "/portal-navigation-items?loadChildren=true",
			}
		}
		items, err := readNavigation(endpoints)
		if err != nil {
			return err
		}
		matches := 0
		for _, item := range items {
			if item.Title != title {
				continue
			}
			matches++
			if item.Visibility != expected {
				return fmt.Errorf("navigation %q visibility %q, want %q", title, item.Visibility, expected)
			}
		}
		if matches != 1 {
			return fmt.Errorf("expected one navigation item %q, found %d", title, matches)
		}
		return nil
	}
}

type navigationItem struct {
	Title      string `json:"title"`
	Visibility string `json:"visibility"`
}

// readNavigation reads the first of the endpoints this APIM version serves.
func readNavigation(endpoints []string) ([]navigationItem, error) {
	for i, endpoint := range endpoints {
		items, status, err := getNavigation(endpoint)
		if status == http.StatusNotFound && i < len(endpoints)-1 {
			continue
		}
		return items, err
	}
	return nil, nil
}

func getNavigation(endpoint string) ([]navigationItem, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, err
	}
	if token := os.Getenv("APIM_SA_TOKEN"); token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	} else {
		request.SetBasicAuth(os.Getenv("APIM_USERNAME"), os.Getenv("APIM_PASSWORD"))
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, 0, fmt.Errorf("read effective navigation: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, fmt.Errorf("read effective navigation: HTTP %d from %s", response.StatusCode, endpoint)
	}
	var body struct {
		Items []navigationItem `json:"items"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, response.StatusCode, err
	}
	return body.Items, response.StatusCode, nil
}
