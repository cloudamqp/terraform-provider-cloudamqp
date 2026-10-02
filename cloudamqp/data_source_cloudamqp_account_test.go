package cloudamqp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccAccount_Basic: Read account information.
func TestAccAccount_Basic(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
				  resource "cloudamqp_instance" "instance-01" {
            name   = "TestAccAccount_Basic-01"
            plan   = "penguin-1"
            region = "amazon-web-services::us-east-1"
            tags   = ["vcr-test"]
          }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_instance.instance-01", "name", "TestAccAccount_Basic-01"),
					resource.TestCheckResourceAttr("cloudamqp_instance.instance-01", "plan", "penguin-1"),
					resource.TestCheckResourceAttr("cloudamqp_instance.instance-01", "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr("cloudamqp_instance.instance-01", "tags.#", "1"),
					resource.TestCheckResourceAttr("cloudamqp_instance.instance-01", "tags.0", "vcr-test"),
				),
			},
			{
				Config: `
				  resource "cloudamqp_instance" "instance-01" {
            name   = "TestAccAccount_Basic-01"
            plan   = "penguin-1"
            region = "amazon-web-services::us-east-1"
            tags   = ["vcr-test"]
          }

					resource "cloudamqp_instance" "instance-02" {
            name   = "TestAccAccount_Basic-02"
            plan   = "bunny-1"
            region = "amazon-web-services::us-east-1"
            tags   = ["vcr-test"]
          }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_instance.instance-02", "name", "TestAccAccount_Basic-02"),
					resource.TestCheckResourceAttr("cloudamqp_instance.instance-02", "plan", "bunny-1"),
					resource.TestCheckResourceAttr("cloudamqp_instance.instance-02", "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr("cloudamqp_instance.instance-02", "tags.#", "1"),
					resource.TestCheckResourceAttr("cloudamqp_instance.instance-02", "tags.0", "vcr-test"),
				),
			},
			{
				Config: `
				  resource "cloudamqp_instance" "instance-01" {
            name   = "TestAccAccount_Basic-01"
            plan   = "penguin-1"
            region = "amazon-web-services::us-east-1"
            tags   = ["vcr-test"]
          }

          resource "cloudamqp_instance" "instance-02" {
            name   = "TestAccAccount_Basic-02"
            plan   = "bunny-1"
            region = "amazon-web-services::us-east-1"
            tags   = ["vcr-test"]
          }

          data "cloudamqp_account" "this" {
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudamqp_account.this", "instances.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("data.cloudamqp_account.this", "instances.*", map[string]string{
						"name":   "TestAccAccount_Basic-01",
						"plan":   "penguin-1",
						"region": "amazon-web-services::us-east-1",
						"tags.#": "1",
						"tags.0": "vcr-test",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.cloudamqp_account.this", "instances.*", map[string]string{
						"name":   "TestAccAccount_Basic-02",
						"plan":   "bunny-1",
						"region": "amazon-web-services::us-east-1",
						"tags.#": "1",
						"tags.0": "vcr-test",
					}),
				),
			},
		},
	})
}
