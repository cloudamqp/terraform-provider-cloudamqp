package cloudamqp

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Unit tests for cloudamqp_upgrade_lavinmq against the fake API in
// resource_cloudamqp_upgrade_rabbitmq_unit_test.go.

const upgradeLavinMQResourceName = "cloudamqp_upgrade_lavinmq.upgrade"

var upgradeLavinMQRequest = fmt.Sprintf("POST /api/instances/%d/actions/upgrade-lavinmq", upgradeTestInstanceID)

// TestUpgradeLavinMQ_Latest: without new_version no version is sent. The API always responds
// 200, which the provider reads as already at the latest version, so there is no wait for the
// nodes. Changing instance_id is an in-place update that sends nothing and keeps the id.
func TestUpgradeLavinMQ_Latest(t *testing.T) {
	t.Parallel()
	fake, server := newFakeUpgradeAPI(t)
	config := upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_lavinmq", "")

	upgradeUnitTest(t, fake, []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(upgradeLavinMQRequest),
				resource.TestCheckResourceAttr(upgradeLavinMQResourceName, "id", fmt.Sprint(upgradeTestInstanceID)),
				resource.TestCheckResourceAttr(upgradeLavinMQResourceName, "instance_id", fmt.Sprint(upgradeTestInstanceID)),
				resource.TestCheckNoResourceAttr(upgradeLavinMQResourceName, "new_version"),
			),
		},
		{
			Config:   config,
			PlanOnly: true,
		},
		{
			Config: strings.Replace(config, fmt.Sprintf("instance_id = %d", upgradeTestInstanceID), "instance_id = 9999", 1),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(upgradeLavinMQRequest),
				resource.TestCheckResourceAttr(upgradeLavinMQResourceName, "id", fmt.Sprint(upgradeTestInstanceID)),
				resource.TestCheckResourceAttr(upgradeLavinMQResourceName, "instance_id", "9999"),
			),
		},
	}, upgradeLavinMQRequest)
}

// TestUpgradeLavinMQ_Specific: new_version sends the version and waits for the nodes. Changing
// it, or removing it, replaces the resource, which upgrades again.
func TestUpgradeLavinMQ_Specific(t *testing.T) {
	t.Parallel()
	fake, server := newFakeUpgradeAPI(t)
	first := upgradeLavinMQRequest + ` {"version":"2.1.0"}`
	second := upgradeLavinMQRequest + ` {"version":"2.2.0"}`

	upgradeUnitTest(t, fake, []resource.TestStep{
		{
			Config: upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_lavinmq", `new_version = "2.1.0"`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(first, upgradeNodesRequest),
				resource.TestCheckResourceAttr(upgradeLavinMQResourceName, "id", fmt.Sprint(upgradeTestInstanceID)),
				resource.TestCheckResourceAttr(upgradeLavinMQResourceName, "new_version", "2.1.0"),
			),
		},
		{
			Config: upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_lavinmq", `new_version = "2.2.0"`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(first, upgradeNodesRequest, second, upgradeNodesRequest),
				resource.TestCheckResourceAttr(upgradeLavinMQResourceName, "new_version", "2.2.0"),
			),
		},
		{
			Config: upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_lavinmq", ""),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(first, upgradeNodesRequest, second, upgradeNodesRequest, upgradeLavinMQRequest),
				resource.TestCheckNoResourceAttr(upgradeLavinMQResourceName, "new_version"),
			),
		},
	}, first, upgradeNodesRequest, second, upgradeNodesRequest, upgradeLavinMQRequest)
}

// TestUpgradeLavinMQ_Error: a failed upgrade gives an error and nothing is stored in the state.
func TestUpgradeLavinMQ_Error(t *testing.T) {
	t.Parallel()
	fake, server := newFakeUpgradeAPI(t)
	fake.failStatus = http.StatusUnprocessableEntity
	fake.failBody = map[string]any{"message": "LavinMQ updates are temporarily disabled, please try again later"}
	config := upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_lavinmq", `new_version = "2.1.0"`)

	upgradeUnitTest(t, fake, []resource.TestStep{
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`unexpected status code: 422`),
		},
		{
			Config:             config,
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	}, upgradeLavinMQRequest+` {"version":"2.1.0"}`)
}
