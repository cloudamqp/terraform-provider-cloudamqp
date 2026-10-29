package cloudamqp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccUpgradableVersionsDataSource_Basic: Read upgradable versions of an instance.
func TestAccUpgradableVersionsDataSource_Basic(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_instance" "instance" {
            name        = "TestAccUpgradableVersionsDataSource_Basic"
            region      = "amazon-web-services::eu-central-1"
            plan        = "squirrel-1"
            tags        = ["vcr-test"]
						rmq_version = "4.0.0"
				  }

          data "cloudamqp_upgradable_versions" "versions" {
            instance_id = cloudamqp_instance.instance.id
          }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudamqp_upgradable_versions.versions", "new_rabbitmq_version", "4.2.8"),
					resource.TestCheckResourceAttr("data.cloudamqp_upgradable_versions.versions", "new_erlang_version", "28.5.0.7"),
				),
			},
		},
	})
}
