package acceptance_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gravitee-io/terraform-provider-apim/tests/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Verifies the create, read, import, and delete lifecycle of the
// `apim_application` resource.
func TestApplicationResource_minimal(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_application.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
			},
			// Verifies resource import.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"organization_id": config.StringVariable(organizationId),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}

func getPrivateKey(t *testing.T) *rsa.PrivateKey {
	const pemPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCxoeCUW5KJxNPxMp+KmCxKLc1Zv9Ny+4CFqcUXVUYH69L3mQ7v
IWrJ9GBfcaA7BPQqUlWxWM+OCEQZH1EZNIuqRMNQVuIGCbz5UQ8w6tS0gcgdeGX7
J7jgCQ4RK3F/PuCM38QBLaHx988qG8NMc6VKErBjctCXFHQt14lerd5KpQIDAQAB
AoGAYrf6Hbk+mT5AI33k2Jt1kcweodBP7UkExkPxeuQzRVe0KVJw0EkcFhywKpr1
V5eLMrILWcJnpyHE5slWwtFHBG6a5fLaNtsBBtcAIfqTQ0Vfj5c6SzVaJv0Z5rOd
7gQF6isy3t3w9IF3We9wXQKzT6q5ypPGdm6fciKQ8RnzREkCQQDZwppKATqQ41/R
vhSj90fFifrGE6aVKC1hgSpxGQa4oIdsYYHwMzyhBmWW9Xv/R+fPyr8ZwPxp2c12
33QwOLPLAkEA0NNUb+z4ebVVHyvSwF5jhfJxigim+s49KuzJ1+A2RaSApGyBZiwS
rWvWkB471POAKUYt5ykIWVZ83zcceQiNTwJBAMJUFQZX5GDqWFc/zwGoKkeR49Yi
MTXIvf7Wmv6E++eFcnT461FlGAUHRV+bQQXGsItR/opIG7mGogIkVXa3E1MCQARX
AAA7eoZ9AEHflUeuLn9QJI/r0hyQQLEtrpwv6rDT1GCWaLII5HJ6NUFVf4TTcqxo
6vdM4QGKTJoO+SaCyP0CQFdpcxSAuzpFcKv0IlJ8XzS/cy+mweCMwyJ1PFEc4FX6
wg/HcAJWY60xZTJDFN+Qfx8ZQvBEin6c2/h+zZi5IVY=
-----END RSA PRIVATE KEY-----
`
	block, _ := pem.Decode([]byte(pemPrivateKey))

	testPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse private key: %s", err)
	}

	return testPrivateKey
}

func getClientTLSCert(t *testing.T) string {
	random := rand.Reader

	ecdsaPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key: %s", err)
	}

	pub := &ecdsaPriv.PublicKey
	priv := getPrivateKey(t)
	sigAlgo := x509.SHA256WithRSA

	commonName := "test.example.com"
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Gravitee"},
			Country:      []string{"FR"},
		},
		NotBefore: time.Now().Add(-24 * time.Hour),
		NotAfter:  time.Now().Add(24 * time.Hour),

		SignatureAlgorithm: sigAlgo,

		SubjectKeyId: []byte{1, 2, 3, 4},
		KeyUsage:     x509.KeyUsageDigitalSignature,

		BasicConstraintsValid: true,
		IsCA:                  false,

		OCSPServer:            []string{"http://ocsp.example.com"},
		IssuingCertificateURL: []string{"http://crt.example.com/ca1.crt"},
	}

	derBytes, err := x509.CreateCertificate(random, &template, &template, pub, priv)
	if err != nil {
		t.Errorf("failed to create certificate: %s", err)
	}

	certData := strings.Builder{}
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}

	if err := pem.Encode(&certData, block); err != nil {
		t.Fatalf("failed to encode certificate: %s", err)
	}

	return certData.String()
}

func TestApplicationResource_withOrgIdFromProvider(t *testing.T) {
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
				ProtoV6ProviderFactories: testProviders(),
				// TODO: We only have one provider in the test environment. This test should be updated once we create more.
				ExpectError: regexp.MustCompile(`Invalid organization or environment`),
			},
		},
	})
}

func TestApplicationResource_withEnvIdFromProvider(t *testing.T) {
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
				ProtoV6ProviderFactories: testProviders(),
				// TODO: We only have one provider in the test environment. This test should be updated once we create more.
				ExpectError: regexp.MustCompile(`Invalid organization or environment`),
			},
		},
	})
}

func TestApplicationResource_overrideOrgIdAndEnvIdFromProvider(t *testing.T) {
	t.Parallel()

	randomId := "test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
				},
				ProtoV6ProviderFactories: testProviders(),
			},
		},
	})
}

// Verifies the create, read, import, and delete lifecycle of the
// `apim_application` resource with as many fields as possible
func TestApplicationResource_all(t *testing.T) {

	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10)

	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId, certObj1, certObj2 := prepareMultiCertsVars(t)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":      config.StringVariable(environmentId),
					"hrid":                config.StringVariable(randomId),
					"organization_id":     config.StringVariable(organizationId),
					"client_certificates": config.ListVariable(certObj1, certObj2),
				},
			},
		},
	})
}

// Verifies the create, read, and delete lifecycle of the
// `apim_application` resource with as many fields as possible
// this is for all versions of apim
func TestApplicationResource_all_versions(t *testing.T) {

	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	cert := getClientTLSCert(t)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":     config.StringVariable(environmentId),
					"hrid":               config.StringVariable(randomId),
					"organization_id":    config.StringVariable(organizationId),
					"client_certificate": config.StringVariable(cert),
				},
			},
		},
	})
}

// Verifies the create, read, import, and delete lifecycle of the
// `apim_application` resource with as many fields as possible
func TestApplicationResource_tlsLegacy(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	cert := getClientTLSCert(t)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":     config.StringVariable(environmentId),
					"hrid":               config.StringVariable(randomId),
					"organization_id":    config.StringVariable(organizationId),
					"client_certificate": config.StringVariable(cert),
				},
			},
		},
	})
}

// Verifies certificate rotation by creating an application with one certificate,
// adding a second, then removing the first.
func TestApplicationResource_certificateRotation(t *testing.T) {
	utils.SkipFor(t, utils.ApimV4_9, utils.ApimV4_10)

	t.Parallel()

	randomId, certObj1, certObj2 := prepareMultiCertsVars(t)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Step 1: Create with one certificate.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":                config.StringVariable(randomId),
					"client_certificates": config.ListVariable(certObj1),
				},
			},
			// Step 2: Add a second certificate (both present).
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":                config.StringVariable(randomId),
					"client_certificates": config.ListVariable(certObj1, certObj2),
				},
			},
			// Step 3: Remove the first certificate (rotation complete).
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid":                config.StringVariable(randomId),
					"client_certificates": config.ListVariable(certObj2),
				},
			},
		},
	})
}

func prepareMultiCertsVars(t *testing.T) (string, config.Variable, config.Variable) {
	randomId := "test-" + acctest.RandString(10)
	cert1 := getClientTLSCert(t)
	cert2 := getClientTLSCert(t)
	yesterday := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	tomorrow := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	certObj1 := config.ObjectVariable(config.Variables{
		"name":      config.StringVariable("cert1"),
		"content":   config.StringVariable(cert1),
		"starts_at": config.StringVariable(yesterday),
		"ends_at":   config.StringVariable(tomorrow),
	})
	certObj2 := config.ObjectVariable(config.Variables{
		"name":      config.StringVariable("cert2"),
		"content":   config.StringVariable(cert2),
		"starts_at": config.StringVariable(yesterday),
		"ends_at":   config.StringVariable(tomorrow),
	})
	return randomId, certObj1, certObj2
}

// Verifies the update the name of `apim_application` resource.
func TestApplicationResource_update(t *testing.T) {
	t.Parallel()

	environmentId := "DEFAULT"
	organizationId := "DEFAULT"
	randomId := "test-" + acctest.RandString(10)
	resourceAddress := "apim_application.test"

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Verifies resource create and read.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"hrid": config.StringVariable(randomId),
					"name": config.StringVariable(randomId + "-original"),
				},
			},
			// Verifies resource import.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"name":            config.StringVariable(randomId + "-original"),
					"organization_id": config.StringVariable(organizationId),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: importStateIDFunc(resourceAddress, []string{"environment_id", "hrid", "organization_id"}, nil),
				ImportStateVerify: true,
			},
			// Verifies resource update.
			{
				ProtoV6ProviderFactories: testProviders(),
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"environment_id":  config.StringVariable(environmentId),
					"hrid":            config.StringVariable(randomId),
					"name":            config.StringVariable(randomId + "-updated"),
					"organization_id": config.StringVariable(organizationId),
				},
			},
			// Testing framework implicitly verifies resource delete.
		},
	})
}
