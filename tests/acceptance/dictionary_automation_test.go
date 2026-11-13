package acceptance_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The value APIM returns in place of an encrypted dictionary value, from 4.13.3.
const dictionaryMask = "••••••••••••"

// The fields of a dictionary read that a write accepts back.
var dictionarySpecFields = []string{"hrid", "name", "description", "deployed", "type", "manual", "dynamic"}

// automationRequest calls the Automation API of the DEFAULT environment with the credentials the
// provider uses, to read or change a dictionary as someone outside Terraform would.
func automationRequest(method, path string, body any) (int, []byte, error) {
	endpoint := strings.TrimRight(os.Getenv("APIM_SERVER_URL"), "/") + "/organizations/DEFAULT/environments/DEFAULT" + path
	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		payload = bytes.NewReader(encoded)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, method, endpoint, payload)
	if err != nil {
		return 0, nil, err
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token := os.Getenv("APIM_SA_TOKEN"); token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	} else {
		request.SetBasicAuth(os.Getenv("APIM_USERNAME"), os.Getenv("APIM_PASSWORD"))
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	return response.StatusCode, content, err
}

func getDictionary(hrid string) (map[string]any, []byte, error) {
	status, content, err := automationRequest(http.MethodGet, "/dictionaries/"+url.PathEscape(hrid), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("read dictionary %s: %w", hrid, err)
	}
	if status != http.StatusOK {
		return nil, nil, fmt.Errorf("read dictionary %s: HTTP %d: %s", hrid, status, content)
	}
	var dictionary map[string]any
	if err := json.Unmarshal(content, &dictionary); err != nil {
		return nil, nil, fmt.Errorf("read dictionary %s: %w", hrid, err)
	}
	return dictionary, content, nil
}

func putDictionary(spec map[string]any) error {
	status, content, err := automationRequest(http.MethodPut, "/dictionaries", spec)
	if err != nil {
		return fmt.Errorf("write dictionary %v: %w", spec["hrid"], err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("write dictionary %v: HTTP %d: %s", spec["hrid"], status, content)
	}
	return nil
}

func deleteDictionary(hrid string) {
	_, _, _ = automationRequest(http.MethodDelete, "/dictionaries/"+url.PathEscape(hrid), nil)
}

// manualMaps returns the plain and the encrypted properties of a dictionary read, creating them
// when the read omits them.
func manualMaps(dictionary map[string]any) (map[string]any, map[string]any) {
	manual, _ := dictionary["manual"].(map[string]any)
	if manual == nil {
		manual = map[string]any{}
		dictionary["manual"] = manual
	}
	plain, _ := manual["properties"].(map[string]any)
	if plain == nil {
		plain = map[string]any{}
		manual["properties"] = plain
	}
	secret, _ := manual["encryptedProperties"].(map[string]any)
	if secret == nil {
		secret = map[string]any{}
		manual["encryptedProperties"] = secret
	}
	return plain, secret
}

// rewriteDictionary reads a dictionary, lets change edit it, and writes it back. Encrypted values
// go back masked, which keeps them.
func rewriteDictionary(hrid string, change func(plain, secret map[string]any)) error {
	dictionary, _, err := getDictionary(hrid)
	if err != nil {
		return err
	}
	spec := make(map[string]any, len(dictionarySpecFields))
	for _, field := range dictionarySpecFields {
		if value, ok := dictionary[field]; ok {
			spec[field] = value
		}
	}
	plain, secret := manualMaps(spec)
	change(plain, secret)
	return putDictionary(spec)
}

// encryptOutsideTerraform moves a plain key to the encrypted properties, as a Console user
// encrypting it would.
func encryptOutsideTerraform(t *testing.T, hrid, key, plaintext string) func() {
	return func() {
		err := rewriteDictionary(hrid, func(plain, secret map[string]any) {
			delete(plain, key)
			secret[key] = map[string]any{"value": plaintext}
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

// addSecretOutsideTerraform adds a secret, as a Console user would.
func addSecretOutsideTerraform(t *testing.T, hrid, key, plaintext string) func() {
	return func() {
		err := rewriteDictionary(hrid, func(_, secret map[string]any) {
			secret[key] = map[string]any{"value": plaintext}
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

// deleteSecretOutsideTerraform deletes a secret, as a Console user would.
func deleteSecretOutsideTerraform(t *testing.T, hrid, key string) func() {
	return func() {
		err := rewriteDictionary(hrid, func(_, secret map[string]any) {
			delete(secret, key)
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

// makePlainOutsideTerraform turns an encrypted key plain. APIM refuses to move it in one write,
// so it takes two, as it does in Terraform.
func makePlainOutsideTerraform(t *testing.T, hrid, key, plaintext string) func() {
	return func() {
		err := rewriteDictionary(hrid, func(_, secret map[string]any) {
			delete(secret, key)
		})
		if err == nil {
			err = rewriteDictionary(hrid, func(plain, _ map[string]any) {
				plain[key] = plaintext
			})
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

func checkEncryptedKey(hrid, key, plaintext string) error {
	dictionary, content, err := getDictionary(hrid)
	if err != nil {
		return err
	}
	plain, secret := manualMaps(dictionary)
	if _, ok := plain[key]; ok {
		return fmt.Errorf("dictionary %s: %s is plain, expected it encrypted", hrid, key)
	}
	entry, _ := secret[key].(map[string]any)
	if entry == nil {
		return fmt.Errorf("dictionary %s: %s is not encrypted", hrid, key)
	}
	if entry["value"] != dictionaryMask {
		return fmt.Errorf("dictionary %s: %s is not masked", hrid, key)
	}
	if bytes.Contains(content, []byte(plaintext)) {
		return fmt.Errorf("dictionary %s: the read returns the value of %s", hrid, key)
	}
	return nil
}

// checkDictionaryEncryptsKey verifies that APIM stores the key encrypted and reads it back masked,
// without ever returning its value.
func checkDictionaryEncryptsKey(hrid, key, plaintext string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		return checkEncryptedKey(hrid, key, plaintext)
	}
}

// checkDictionaryEncryptsOnFetch verifies the encryption policy APIM stores for a dynamic dictionary.
func checkDictionaryEncryptsOnFetch(hrid string, expected bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		dictionary, _, err := getDictionary(hrid)
		if err != nil {
			return err
		}
		dynamic, _ := dictionary["dynamic"].(map[string]any)
		encryption, _ := dynamic["encryption"].(map[string]any)
		actual, ok := encryption["encryptOnFetch"].(bool)
		if !ok {
			return fmt.Errorf("dictionary %s: APIM returns no encryption policy", hrid)
		}
		if actual != expected {
			return fmt.Errorf("dictionary %s: encryptOnFetch is %v in APIM, expected %v", hrid, actual, expected)
		}
		return nil
	}
}

// requireMaskedRead fails the test unless APIM reads the key back masked, so a step proving that
// the mask causes no change does not pass against an APIM returning no value.
func requireMaskedRead(t *testing.T, hrid, key, plaintext string) func() {
	return func() {
		if err := checkEncryptedKey(hrid, key, plaintext); err != nil {
			t.Fatal(err)
		}
	}
}

// checkDictionaryLacksKey verifies that APIM holds the key neither plain nor encrypted.
func checkDictionaryLacksKey(hrid, key string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		dictionary, _, err := getDictionary(hrid)
		if err != nil {
			return err
		}
		plain, secret := manualMaps(dictionary)
		_, isPlain := plain[key]
		_, isSecret := secret[key]
		if isPlain || isSecret {
			return fmt.Errorf("dictionary %s: APIM still holds %s", hrid, key)
		}
		return nil
	}
}

// planWithoutValue fails when a value appears anywhere in the plan, prior state included.
type planWithoutValue struct {
	value string
}

func (c planWithoutValue) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	content, err := json.Marshal(req.Plan)
	if err != nil {
		resp.Error = err
		return
	}
	if bytes.Contains(content, []byte(c.value)) {
		resp.Error = fmt.Errorf("the plan shows a value it must hide")
	}
}

func expectPlanWithoutValue(value string) plancheck.PlanCheck {
	return planWithoutValue{value: value}
}
