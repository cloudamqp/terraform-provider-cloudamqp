package cloudamqp

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Extra environment variables needed to record this test. Used self signed certififactes and loaded
// them into environment variables for recording.
// export TEST_CERTIFICATE_CA=$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' certs/ca.pem)
// export TEST_CERTIFICATE_CERT=$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' certs/server.crt)
// export TEST_CERTIFICATE_PRIVATE_KEY=$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' certs/server.key)

// TestAccCustomCertificate_Basic: Creating dedicated AWS instance and upload custom certificate.
func TestAccCustomCertificate_Basic(t *testing.T) {
	t.Parallel()

	instanceResourceName := "cloudamqp_instance.instance"
	customCertificateResourceName := "cloudamqp_custom_certificate.custom_cert"

	// Set sanitized value for playback and use real value for recording
	testCertificateCA := "TEST_CERTIFICATE_CA"
	testCertificateCert := "TEST_CERTIFICATE_CERT"
	testCertificatePrivateKey := "TEST_CERTIFICATE_PRIVATE_KEY"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testCertificateCA = os.Getenv("TEST_CERTIFICATE_CA")
		testCertificateCert = os.Getenv("TEST_CERTIFICATE_CERT")
		testCertificatePrivateKey = os.Getenv("TEST_CERTIFICATE_PRIVATE_KEY")
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccCustomCertificate_Basic"
						plan   = "bunny-1"
						region = "amazon-web-services::us-east-1"
						tags   = []
					}
					
					resource "cloudamqp_custom_certificate" "custom_cert" {
						instance_id = cloudamqp_instance.instance.id
						ca          = "%s"
						cert        = "%s"
						private_key = "%s"
						sni_hosts   = "my.custom.domain"
					}
				`, testCertificateCA, testCertificateCert, testCertificatePrivateKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccCustomCertificate_Basic"),
					resource.TestCheckResourceAttr(customCertificateResourceName, "sni_hosts", "my.custom.domain"),
					resource.TestCheckResourceAttr(customCertificateResourceName, "version", "1"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccCustomCertificate_Basic"
						plan   = "bunny-1"
						region = "amazon-web-services::us-east-1"
						tags   = []
					}
					
					# Increment version to trigger force new, restore default and re-upload custom certificate.
					resource "cloudamqp_custom_certificate" "custom_cert" {
						instance_id = cloudamqp_instance.instance.id
						ca          = "%s"
						cert        = "%s"
						private_key = "%s"
						sni_hosts   = "my.custom.domain"
						version		  = 2
					}
				`, testCertificateCA, testCertificateCert, testCertificatePrivateKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccCustomCertificate_Basic"),
					resource.TestCheckResourceAttr(customCertificateResourceName, "sni_hosts", "my.custom.domain"),
					resource.TestCheckResourceAttr(customCertificateResourceName, "version", "2"),
				),
			},
		},
	})
}

// TestAccCustomCertificate_KeyID: Creating dedicated AWS instance, upload custom certificate and
// update with key identifier.
func TestAccCustomCertificate_KeyID(t *testing.T) {
	t.Parallel()

	instanceResourceName := "cloudamqp_instance.instance"
	customCertificateResourceName := "cloudamqp_custom_certificate.custom_cert"

	// Set sanitized value for playback and use real value for recording
	testCertificateCA := "TEST_CERTIFICATE_CA"
	testCertificateCert := "TEST_CERTIFICATE_CERT"
	testCertificatePrivateKey := "TEST_CERTIFICATE_PRIVATE_KEY"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testCertificateCA = os.Getenv("TEST_CERTIFICATE_CA")
		testCertificateCert = os.Getenv("TEST_CERTIFICATE_CERT")
		testCertificatePrivateKey = os.Getenv("TEST_CERTIFICATE_PRIVATE_KEY")
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccCustomCertificate_KeyID"
						plan   = "bunny-1"
						region = "amazon-web-services::us-east-1"
						tags   = []
					}
					
					resource "cloudamqp_custom_certificate" "custom_cert" {
						instance_id = cloudamqp_instance.instance.id
						ca          = "%s"
						cert        = "%s"
						private_key = "%s"
						sni_hosts   = "my.custom.domain"
						key_id		  = "a918beb8-fee4-4de1-b0d5-873e2cb0eba2"
					}
				`, testCertificateCA, testCertificateCert, testCertificatePrivateKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccCustomCertificate_KeyID"),
					resource.TestCheckResourceAttr(customCertificateResourceName, "sni_hosts", "my.custom.domain"),
					resource.TestCheckResourceAttr(customCertificateResourceName, "key_id", "a918beb8-fee4-4de1-b0d5-873e2cb0eba2"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccCustomCertificate_KeyID"
						plan   = "bunny-1"
						region = "amazon-web-services::us-east-1"
						tags   = []
					}
					
					# Increment version to trigger force new, restore default and re-upload custom certificate.
					resource "cloudamqp_custom_certificate" "custom_cert" {
						instance_id = cloudamqp_instance.instance.id
						ca          = "%s"
						cert        = "%s"
						private_key = "%s"
						sni_hosts   = "my.custom.domain"
						key_id		  = "53f188e8-a81d-4232-b5f1-7379b0223bb1"
					}
				`, testCertificateCA, testCertificateCert, testCertificatePrivateKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccCustomCertificate_KeyID"),
					resource.TestCheckResourceAttr(customCertificateResourceName, "sni_hosts", "my.custom.domain"),
					resource.TestCheckResourceAttr(customCertificateResourceName, "key_id", "53f188e8-a81d-4232-b5f1-7379b0223bb1"),
				),
			},
		},
	})
}
