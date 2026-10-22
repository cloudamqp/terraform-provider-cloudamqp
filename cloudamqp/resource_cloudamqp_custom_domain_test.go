package cloudamqp

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccCustomDomain_Update: Create custom domain, then update hostname.
func TestAccCustomDomain_Update(t *testing.T) {
	t.Parallel()

	var (
		customDomainResourceName = "cloudamqp_custom_domain.custom_domain"
		testInstanceID           = "407805" // pre-existing instance

		// Recording benefits from the default sleep (10 s) between retries, matching
		// production behavior. Replaying the cassette doesn't need the real wait, so a
		// lower sleep speeds up the test.
		sleep = "1"
	)
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		sleep = "10"
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_custom_domain" "custom_domain" {
						instance_id = %s
						hostname    = "vcr-test.myddns.me"
						sleep       = %s
					}
				`, testInstanceID, sleep),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(customDomainResourceName, "hostname", "vcr-test.myddns.me"),
					resource.TestCheckResourceAttr(customDomainResourceName, "sleep", sleep),
					resource.TestCheckResourceAttr(customDomainResourceName, "timeout", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_custom_domain" "custom_domain" {
						instance_id = %s
						hostname    = "vcr-update.myddns.me"
						sleep       = %s
					}
				`, testInstanceID, sleep),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(customDomainResourceName, "hostname", "vcr-update.myddns.me"),
					resource.TestCheckResourceAttr(customDomainResourceName, "sleep", sleep),
					resource.TestCheckResourceAttr(customDomainResourceName, "timeout", "1800"),
				),
			},
		},
	})
}

// TestAccCustomDomain_Import: Import existing custom domain.
func TestAccCustomDomain_Import(t *testing.T) {
	t.Parallel()

	var (
		testInstanceID = "407805" // pre-existing instance
		sleep          = "1"
	)
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		sleep = "10"
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_custom_domain" "settings" {
						instance_id = %s
						hostname    = "vcr-test.myddns.me"
						sleep       = %s
					}
				`, testInstanceID, sleep),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_custom_domain.settings", "hostname", "vcr-test.myddns.me"),
					resource.TestCheckResourceAttr("cloudamqp_custom_domain.settings", "sleep", sleep),
					resource.TestCheckResourceAttr("cloudamqp_custom_domain.settings", "timeout", "1800"),
				),
			},
			{
				ResourceName:            "cloudamqp_custom_domain.settings",
				ImportStateId:           testInstanceID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sleep"},
			},
		},
	})
}
