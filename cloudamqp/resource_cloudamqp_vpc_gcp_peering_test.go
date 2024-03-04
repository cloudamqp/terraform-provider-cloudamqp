package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccVpcGcpPeering_Basic: Peer two VPCs hosted in GCP
func TestAccVpcGcpPeering_Basic(t *testing.T) {
	var (
		fileNames = []string{"vpcs_and_instances", "vpc_gcp_peerings"}

		vpcNameFirst           = "cloudamqp_vpc.vpc_first"
		instanceNameFirst      = "cloudamqp_instance.instance_first"
		dataVpcInfoFirst       = "data.cloudamqp_vpc_gcp_info.vpc_info_first"
		vpcPeeringRequestFirst = "cloudamqp_vpc_gcp_peering.vpc_peering_request_first"

		vpcNameSecond           = "cloudamqp_vpc.vpc_second"
		instanceNameSecond      = "cloudamqp_instance.instance_second"
		dataVpcInfoSecond       = "data.cloudamqp_vpc_gcp_info.vpc_info_second"
		vpcPeeringRequestSecond = "cloudamqp_vpc_gcp_peering.vpc_peering_request_second"

		params = map[string]string{
			"VpcNameFirst":        "TestAccVpcGcpPeering_Basic_First",
			"VpcIDFirst":          fmt.Sprintf("%s.id", vpcNameFirst),
			"VpcRegionFirst":      "google-compute-engine::us-east1",
			"InstanceNameFirst":   "TestAccVpcGcpPeering_Basic_First",
			"InstanceIDFirst":     fmt.Sprintf("%s.id", instanceNameFirst),
			"InstanceRegionFirst": "google-compute-engine::us-east1",
			"PeerNetworkUriFirst": fmt.Sprintf("%s.network", dataVpcInfoFirst),

			"VpcNameSecond":        "TestAccVpcGcpPeering_Basic_Second",
			"VpcIDSecond":          fmt.Sprintf("%s.id", vpcNameSecond),
			"VpcRegionSecond":      "google-compute-engine::us-east1",
			"InstanceNameSecond":   "TestAccVpcGcpPeering_Basic_Second",
			"InstanceIDSecond":     fmt.Sprintf("%s.id", instanceNameSecond),
			"InstanceRegionSecond": "google-compute-engine::us-east1",
			"PeerNetworkUriSecond": fmt.Sprintf("%s.network", dataVpcInfoSecond),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vpcPeeringRequestFirst, "state", "ACTIVE"),
					resource.TestCheckResourceAttr(vpcPeeringRequestSecond, "state", "ACTIVE"),
				),
			},
		},
	})
}
