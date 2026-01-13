package cloudamqp

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Extra environment variables needed to record this test. Used self signed certificates and loaded
// them into environment variables for recording.
// export TEST_TRUST_STORE_CA=$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' certs/ca.pem)

// TestAccTrustStore_Http: Creating dedicated AWS instance and configure trust store with http
// provider, minimal required values and import.
func TestAccTrustStore_Http(t *testing.T) {
	t.Parallel()

	instanceResourceName := "cloudamqp_instance.instance"
	trustStoreResourceName := "cloudamqp_trust_store.trust_store"

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccTrustStore_Http"
						plan   = "bunny-1"
						region = "amazon-web-services::us-east-1"
						tags   = []
					}
					
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id  = cloudamqp_instance.instance.id
						http {
							url = "https://valid.example.com/trust-store"
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccTrustStore_Http"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "http.url", "https://valid.example.com/trust-store"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "30"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
			{
				ResourceName:      trustStoreResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"sleep",
					"timeout",
				},
			},
		},
	})
}

// TestAccTrustStore_HttpWithCaVersion: Pre-created dedicated instance, configure trust store with http
// provider and CA certificate with version trigger.
func TestAccTrustStore_HttpWithCaVersion(t *testing.T) {
	t.Parallel()

	trustStoreResourceName := "cloudamqp_trust_store.trust_store"

	// Set sanitized value for playback and use real value for recording
	testTrustStoreCA := "TEST_TRUST_STORE_CA"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testTrustStoreCA = os.Getenv("TEST_TRUST_STORE_CA")
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = 922
						http {
							url    = "https://trust-store.example.com/certs"
							cacert = "%s"
						}
						refresh_interval = 60
						version          = 1
						sleep            = 10
						timeout          = 1800
					}
				`, testTrustStoreCA),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(trustStoreResourceName, "http.url", "https://trust-store.example.com/certs"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "60"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "version", "1"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(`					
					# Increment version to trigger update CA certificate.
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = 922
						http {
							url    = "https://trust-store.example.com/certs"
							cacert = "%s"
						}
						refresh_interval = 60
						version          = 2
						sleep            = 10
						timeout          = 1800
					}
				`, testTrustStoreCA),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(trustStoreResourceName, "http.url", "https://trust-store.example.com/certs"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "60"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "version", "2"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
		},
	})
}

// TestAccTrustStore_HttpWithCA: Pre-created dedicated instance, configure trust store with http
// provider and CA certificate with key identifier trigger.
func TestAccTrustStore_HttpWithCaKeyID(t *testing.T) {
	t.Parallel()

	trustStoreResourceName := "cloudamqp_trust_store.trust_store"

	// Set sanitized value for playback and use real value for recording
	testTrustStoreCA := "TEST_TRUST_STORE_CA"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testTrustStoreCA = os.Getenv("TEST_TRUST_STORE_CA")
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = 922
						http {
							url    = "https://trust-store.example.com/certs"
							cacert = "%s"
						}
						refresh_interval = 60
						key_id           = "a918beb8-fee4-4de1-b0d5-873e2cb0eba2"
						sleep            = 10
						timeout          = 1800
					}
				`, testTrustStoreCA),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(trustStoreResourceName, "http.url", "https://trust-store.example.com/certs"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "60"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "key_id", "a918beb8-fee4-4de1-b0d5-873e2cb0eba2"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(`
					# Update key identifier to trigger update of CA certificate.
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = 922
						http {
							url    = "https://trust-store.example.com/certs"
							cacert = "%s"
						}
						refresh_interval = 60
						key_id           = "53f188e8-a81d-4232-b5f1-7379b0223bb1"
						sleep            = 10
						timeout          = 1800
					}
				`, testTrustStoreCA),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(trustStoreResourceName, "http.url", "https://trust-store.example.com/certs"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "60"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "key_id", "53f188e8-a81d-4232-b5f1-7379b0223bb1"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
		},
	})
}

// TestAccTrustStore_FileVersion: Pre-created dedicated instance, configure trust store with file
// certificates with version trigger.
func TestAccTrustStore_FileVersion(t *testing.T) {
	t.Parallel()

	trustStoreResourceName := "cloudamqp_trust_store.trust_store"

	// Set sanitized value for playback and use real value for recording
	testTrustStoreCert := "TEST_TRUST_STORE_CERT"
	testTrustStoreCert_02 := "TEST_TRUST_STORE_CERT_2"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testTrustStoreCert = os.Getenv("TEST_TRUST_STORE_CERT")
		testTrustStoreCert_02 = os.Getenv("TEST_TRUST_STORE_CERT_2")
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = 922
						file {
							certificates = ["%s", "%s"]
						}
						refresh_interval = 60
						version          = 1
						sleep            = 10
						timeout          = 1800
					}
				`, testTrustStoreCert, testTrustStoreCert_02),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "60"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "version", "1"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(`
					# Increment version to trigger update certificates.
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = 922
						file {
							certificates = ["%s", "%s"]
						}
						refresh_interval = 60
						version          = 2
						sleep            = 10
						timeout          = 1800
					}
				`, testTrustStoreCert, testTrustStoreCert_02),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "60"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "version", "2"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
		},
	})
}

// TestAccTrustStore_FileKeyID: Pre-created dedicated instance, configure trust store with file
// certificates with key identifier trigger.
func TestAccTrustStore_FileKeyID(t *testing.T) {
	t.Parallel()

	trustStoreResourceName := "cloudamqp_trust_store.trust_store"

	// Set sanitized value for playback and use real value for recording
	testTrustStoreCert := "TEST_TRUST_STORE_CERT"
	testTrustStoreCert_02 := "TEST_TRUST_STORE_CERT_2"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testTrustStoreCert = os.Getenv("TEST_TRUST_STORE_CERT")
		testTrustStoreCert_02 = os.Getenv("TEST_TRUST_STORE_CERT_2")
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = 922
						file {
							certificates = ["%s", "%s"]
						}
						refresh_interval = 60
						key_id           = "d31cb790-a57b-400b-98cc-c13ad5a9e860"
						sleep            = 10
						timeout          = 1800
					}
				`, testTrustStoreCert, testTrustStoreCert_02),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "60"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "key_id", "d31cb790-a57b-400b-98cc-c13ad5a9e860"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(`
					# Update key identifier to trigger update of certificates.
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = 922
						file {
							certificates = ["%s", "%s"]
						}
						refresh_interval = 60
						key_id           = "6e0e9b7a-268a-4621-9cd8-1d66af48d045"
						sleep            = 10
						timeout          = 1800
					}
				`, testTrustStoreCert, testTrustStoreCert_02),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "60"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "key_id", "6e0e9b7a-268a-4621-9cd8-1d66af48d045"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
		},
	})
}

// TestAccTrustStore_MissingProvider: Test missing provider block error handling.
func TestAccTrustStore_MissingProvider(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccTrustStore_MissingProvider"
						plan   = "bunny-1"
						region = "amazon-web-services::us-east-1"
						tags   = []
					}

					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = cloudamqp_instance.instance.id
						refresh_interval = 60
						version          = 1
						sleep            = 10
						timeout          = 1800
					}
				`,
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}
