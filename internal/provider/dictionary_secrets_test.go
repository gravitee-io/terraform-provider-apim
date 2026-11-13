package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tfTypes "github.com/gravitee-io/terraform-provider-apim/internal/provider/types"
	"github.com/gravitee-io/terraform-provider-apim/internal/sdk/models/shared"
)

const dictionaryMask = "••••••••••••"

func secretValues(values map[string]string) map[string]tfTypes.EncryptableValue {
	if values == nil {
		return nil
	}
	secrets := make(map[string]tfTypes.EncryptableValue, len(values))
	for key, value := range values {
		secrets[key] = tfTypes.EncryptableValue{Value: types.StringValue(value)}
	}
	return secrets
}

func plainValues(values map[string]string) map[string]types.String {
	if values == nil {
		return nil
	}
	plain := make(map[string]types.String, len(values))
	for key, value := range values {
		plain[key] = types.StringValue(value)
	}
	return plain
}

func maskedAnswer(plain map[string]string, secretKeys ...string) *shared.DictionaryState {
	manual := &shared.ManualDictionarySpec{Properties: plain}
	if len(secretKeys) > 0 {
		manual.EncryptedProperties = make(map[string]shared.EncryptableValue, len(secretKeys))
		for _, key := range secretKeys {
			manual.EncryptedProperties[key] = shared.EncryptableValue{Value: dictionaryMask}
		}
	}
	return &shared.DictionaryState{Hrid: "secrets", Name: "secrets", Type: shared.DictionaryTypeManual, Manual: manual}
}

func TestDictionaryRefreshKeepingSecrets(t *testing.T) {
	testCases := map[string]struct {
		prior          *tfTypes.ManualDictionarySpec
		answer         *shared.DictionaryState
		expectedPlain  map[string]types.String
		expectedSecret map[string]tfTypes.EncryptableValue
	}{
		"a secret APIM reads back masked keeps its value": {
			prior:          &tfTypes.ManualDictionarySpec{Properties: plainValues(map[string]string{"env": "test"}), EncryptedProperties: secretValues(map[string]string{"api-key": "s3cr3t"})},
			answer:         maskedAnswer(map[string]string{"env": "test"}, "api-key"),
			expectedPlain:  plainValues(map[string]string{"env": "test"}),
			expectedSecret: secretValues(map[string]string{"api-key": "s3cr3t"}),
		},
		"a secret added outside Terraform keeps the mask": {
			prior:          &tfTypes.ManualDictionarySpec{EncryptedProperties: secretValues(map[string]string{"api-key": "s3cr3t"})},
			answer:         maskedAnswer(nil, "api-key", "token"),
			expectedSecret: secretValues(map[string]string{"api-key": "s3cr3t", "token": dictionaryMask}),
		},
		"a secret added outside Terraform to a dictionary without secrets keeps the mask": {
			prior:          &tfTypes.ManualDictionarySpec{Properties: plainValues(map[string]string{"env": "test"})},
			answer:         maskedAnswer(map[string]string{"env": "test"}, "token"),
			expectedPlain:  plainValues(map[string]string{"env": "test"}),
			expectedSecret: secretValues(map[string]string{"token": dictionaryMask}),
		},
		"a secret deleted outside Terraform is dropped": {
			prior:          &tfTypes.ManualDictionarySpec{Properties: plainValues(map[string]string{"env": "test"}), EncryptedProperties: secretValues(map[string]string{"api-key": "s3cr3t"})},
			answer:         maskedAnswer(map[string]string{"env": "test"}),
			expectedPlain:  plainValues(map[string]string{"env": "test"}),
			expectedSecret: secretValues(map[string]string{}),
		},
		"a secret made plain outside Terraform stays out of properties": {
			prior:          &tfTypes.ManualDictionarySpec{Properties: plainValues(map[string]string{"env": "test"}), EncryptedProperties: secretValues(map[string]string{"api-key": "s3cr3t"})},
			answer:         maskedAnswer(map[string]string{"env": "test", "api-key": "now-plain"}),
			expectedPlain:  plainValues(map[string]string{"env": "test"}),
			expectedSecret: secretValues(map[string]string{}),
		},
		"a dictionary holding only a secret made plain outside Terraform has no plain properties": {
			prior:          &tfTypes.ManualDictionarySpec{EncryptedProperties: secretValues(map[string]string{"api-key": "s3cr3t"})},
			answer:         maskedAnswer(map[string]string{"api-key": "now-plain"}),
			expectedSecret: secretValues(map[string]string{}),
		},
		"an empty secrets map stays empty": {
			prior:          &tfTypes.ManualDictionarySpec{Properties: plainValues(map[string]string{"env": "test"}), EncryptedProperties: secretValues(map[string]string{})},
			answer:         maskedAnswer(map[string]string{"env": "test"}),
			expectedPlain:  plainValues(map[string]string{"env": "test"}),
			expectedSecret: secretValues(map[string]string{}),
		},
		"an import keeps the mask": {
			answer:         maskedAnswer(map[string]string{"env": "test"}, "api-key"),
			expectedPlain:  plainValues(map[string]string{"env": "test"}),
			expectedSecret: secretValues(map[string]string{"api-key": dictionaryMask}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			data := DictionaryResourceModel{Manual: tc.prior}

			diags := data.refreshKeepingSecrets(context.Background(), tc.answer)

			require.False(t, diags.HasError(), diags)
			require.NotNil(t, data.Manual)
			assert.Equal(t, tc.expectedPlain, data.Manual.Properties)
			assert.Equal(t, tc.expectedSecret, data.Manual.EncryptedProperties)
		})
	}
}
