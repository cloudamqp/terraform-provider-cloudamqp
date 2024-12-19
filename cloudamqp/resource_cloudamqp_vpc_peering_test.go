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
		ProviderFactories:         testAccProviderFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			// {
			// 	ResourceName:            vpcResourceName,
			// 	ImportStateId:           vpcID,
			// 	ImportState:             true,
			// 	ImportStateVerify:       true,
			// 	ImportStateVerifyIgnore: []string{},
			// },
			// {
			// 	ResourceName:            instanceResourceName,
			// 	ImportStateId:           instanceID,
			// 	ImportState:             true,
			// 	ImportStateVerify:       true,
			// 	ImportStateVerifyIgnore: []string{},
			// },
			// {
			// 	Config: configuration.GetTemplatedConfig(t, fileNames, params),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr(vpcResourceName, "name", params["VpcName"]),
			// 		resource.TestCheckResourceAttr(vpcResourceName, "id", vpcID),
			// 		resource.TestCheckResourceAttr(vpcResourceName, "subnet", params["VpcSubnet"]),
			// 		resource.TestCheckResourceAttr(vpcResourceName, "region", params["VpcRegion"]),
			// 		resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
			// 		resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
			// 		resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
			// 	),
			// },
			// {
			// 	ResourceName:            vpcPeeringResourceName,
			// 	ImportStateIdFunc:       testAccImportVpcPeeringStateIdFunc("vpc", vpcResourceName, peeringID),
			// 	ImportState:             true,
			// 	ImportStateVerify:       true,
			// 	ImportStateVerifyIgnore: []string{},
			// },
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
