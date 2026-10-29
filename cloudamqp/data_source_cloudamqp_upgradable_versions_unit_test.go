package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Unit tests for the cloudamqp_upgradable_versions data source against the fake API in
// resource_cloudamqp_upgrade_rabbitmq_unit_test.go.

const upgradableVersionsDataSourceName = "data.cloudamqp_upgradable_versions.versions"

func upgradableVersionsUnitTestConfig(baseURL string) string {
	return fmt.Sprintf(`
provider "cloudamqp" {
  apikey  = "test"
  baseurl = %q
}

data "cloudamqp_upgradable_versions" "versions" {
  instance_id = %d
}
`, baseURL, upgradeTestInstanceID)
}

// TestUpgradableVersions_Read: the versions are read and the id is the instance id.
func TestUpgradableVersions_Read(t *testing.T) {
	t.Parallel()
	fake, server := newFakeUpgradeAPI(t)
	fake.versions = map[string]any{"new_rabbitmq_version": "4.2.7", "new_erlang_version": "28.4.2"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: upgradeUnitTestProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: upgradableVersionsUnitTestConfig(server.URL),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(upgradableVersionsDataSourceName, "id", fmt.Sprint(upgradeTestInstanceID)),
					resource.TestCheckResourceAttr(upgradableVersionsDataSourceName, "instance_id", fmt.Sprint(upgradeTestInstanceID)),
					resource.TestCheckResourceAttr(upgradableVersionsDataSourceName, "new_rabbitmq_version", "4.2.7"),
					resource.TestCheckResourceAttr(upgradableVersionsDataSourceName, "new_erlang_version", "28.4.2"),
				),
			},
		},
	})
}

// TestUpgradableVersions_NullVersions: the API can respond with null versions, e.g. when the
// instance has no RabbitMQ version yet. They are read as empty, not as a "<nil>" string.
func TestUpgradableVersions_NullVersions(t *testing.T) {
	t.Parallel()
	fake, server := newFakeUpgradeAPI(t)
	fake.versions = map[string]any{"new_rabbitmq_version": nil, "new_erlang_version": nil}

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: upgradeUnitTestProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: upgradableVersionsUnitTestConfig(server.URL),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(upgradableVersionsDataSourceName, "id", fmt.Sprint(upgradeTestInstanceID)),
					resource.TestCheckResourceAttr(upgradableVersionsDataSourceName, "new_rabbitmq_version", ""),
					resource.TestCheckResourceAttr(upgradableVersionsDataSourceName, "new_erlang_version", ""),
				),
			},
		},
	})
}
