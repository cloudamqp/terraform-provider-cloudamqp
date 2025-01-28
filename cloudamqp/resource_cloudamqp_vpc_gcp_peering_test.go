package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccVpcGcpPeering_Basic: Create standalone VPC and instance, enable VPC Connect and import.
func TestAccVpcGcpPeering_Basic(t *testing.T) {
	var (
		fileNames                 = []string{"vpc_and_instance", "vpc_gcp_peering"}
		vpcResourceName           = "cloudamqp_vpc.vpc"
		instanceResourceName      = "cloudamqp_instance.instance"
		vpcGcpPeeringResourceName = "cloudamqp_vpc_gcp_peering.vpc_peering"
		region                    = "google-compute-engine::europe-west1"
		peerNetworkUri            = "https://www.googleapis.com/compute/v1/projects/playground-84codes/global/networks/vpc-ljopathx"

		params = map[string]string{
			"VpcName":             "TestAccVpcGcpPeering_Basic",
			"VpcRegion":           region,
			"VpcSubnet":           "10.56.72.0/24",
			"InstanceName":        "TestAccVpcGcpPeering_Basic",
			"InstanceRegion":      region,
			"InstanceID":          fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":        "bunny-1",
			"PeerNetworkUri":      peerNetworkUri,
			"WaitOnPeeringStatus": "true",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vpcResourceName, "name", params["VpcName"]),
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(vpcGcpPeeringResourceName, "auto_create_routes", "true"),
					resource.TestCheckResourceAttr(vpcGcpPeeringResourceName, "wait_on_peering_status", "true"),
					resource.TestCheckResourceAttr(vpcGcpPeeringResourceName, "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            vpcGcpPeeringResourceName,
				ImportStateIdFunc:       testAccImportVpcPeeringStateIdFunc("vpc", vpcResourceName, peerNetworkUri),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_on_peering_status"},
			},
		},
	})
}

func testAccImportVpcPeeringStateIdFunc(importType, resourceName, peerNetworkUri string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("No resource id set")
		}
		resourceID := rs.Primary.ID
		return fmt.Sprintf("%s,%s,%s", importType, resourceID, peerNetworkUri), nil
	}
}
