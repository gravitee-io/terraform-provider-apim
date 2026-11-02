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
		endpoint := strings.TrimRight(base, "/") + "/organizations/" + url.PathEscape(attributes["organization_id"]) + "/environments/" + url.PathEscape(attributes["environment_id"])
		if apiScoped {
			endpoint += "/apis/" + url.PathEscape(attributes["id"])
		}
		endpoint += "/portal-navigation-items?loadChildren=true"
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return err
		}
		if token := os.Getenv("APIM_SA_TOKEN"); token != "" {
			request.Header.Set("Authorization", "Bearer "+token)
		} else {
			request.SetBasicAuth(os.Getenv("APIM_USERNAME"), os.Getenv("APIM_PASSWORD"))
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return fmt.Errorf("read effective navigation: %w", err)
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("read effective navigation: HTTP %d", response.StatusCode)
		}
		var body struct {
			Items []struct {
				Title      string `json:"title"`
				Visibility string `json:"visibility"`
			} `json:"items"`
		}
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			return err
		}
		matches := 0
		for _, item := range body.Items {
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
