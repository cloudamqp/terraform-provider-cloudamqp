package cloudamqp

import (
	"fmt"
	"os"
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

// TestAccTrustStore_HttpWithCA: Creating dedicated AWS instance and configure trust store with http
// provider and CA certificate.
func TestAccTrustStore_HttpWithCA(t *testing.T) {
	t.Parallel()

	instanceResourceName := "cloudamqp_instance.instance"
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
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccTrustStore_HttpWithCA"
						plan   = "bunny-1"
						region = "amazon-web-services::us-east-1"
						tags   = []
					}
					
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = cloudamqp_instance.instance.id
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
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccTrustStore_HttpWithCA"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "http.url", "https://trust-store.example.com/certs"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "60"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "version", "1"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccTrustStore_HttpWithCA"
						plan   = "bunny-1"
						region = "amazon-web-services::us-east-1"
						tags   = []
					}
					
					# Increment version to trigger update with new CA certificate.
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = cloudamqp_instance.instance.id
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
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccTrustStore_HttpWithCA"),
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

// TestAccTrustStore_File: Creating dedicated AWS instance and configure trust store with file certificates.
func TestAccTrustStore_File(t *testing.T) {
	t.Parallel()

	instanceResourceName := "cloudamqp_instance.instance"
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
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccTrustStore_File"
						plan   = "bunny-1"
						region = "amazon-web-services::us-east-1"
						tags   = []
					}
					
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = cloudamqp_instance.instance.id
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
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccTrustStore_File"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "60"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "version", "1"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccTrustStore_File"
						plan   = "bunny-1"
						region = "amazon-web-services::us-east-1"
						tags   = []
					}
					
					# Increment version to trigger update with new CA certificate.
					resource "cloudamqp_trust_store" "trust_store" {
						instance_id      = cloudamqp_instance.instance.id
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
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccTrustStore_File"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "refresh_interval", "60"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "version", "2"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(trustStoreResourceName, "timeout", "1800"),
				),
			},
		},
	})
}
