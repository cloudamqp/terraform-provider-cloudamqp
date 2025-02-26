package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccVpcPeering_Basic: Accept VPC peering and import.
func TestAccVpcPeering_Basic(t *testing.T) {
	var (
		fileNames              = []string{"vpc_peering"}
		vpcPeeringResourceName = "cloudamqp_vpc_peering.accepter"

		params = map[string]string{
			"VpcID":     "186",
			"PeeringID": "pcx-03ebf8d9ceac304d7",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vpcPeeringResourceName, "vpc_id", params["VpcID"]),
					resource.TestCheckResourceAttr(vpcPeeringResourceName, "peering_id", params["PeeringID"]),
					resource.TestCheckResourceAttr(vpcPeeringResourceName, "status", "active"),
				),
			},
			{
				ResourceName:            vpcPeeringResourceName,
				ImportStateId:           fmt.Sprintf("%s,%s,%s", "vpc", params["VpcID"], params["PeeringID"]),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}
