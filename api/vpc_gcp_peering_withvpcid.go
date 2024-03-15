package api

// VPC peering for GCP, using vpcID as identifier.

import (
	"fmt"
	"log"
)

// RequestVpcGcpPeeringWithVpcId: requests a VPC peering from an instance.
func (api *API) RequestVpcGcpPeeringWithVpcId(vpcID string, params map[string]interface{},
	waitOnStatus bool, sleep, timeout int) (map[string]interface{}, error) {

	path := fmt.Sprintf("api/vpcs/%s/vpc-peering", vpcID)
	attempt, data, err := api.requestVpcGcpPeeringWithRetry(path, params, waitOnStatus, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}

	if waitOnStatus {
		log.Printf("[DEBUG] go-api::vpc_gcp_peering_withvpcid::request waiting for active state")
		err = api.waitForGcpPeeringStatus(path, data["peering"].(string), attempt, sleep, timeout)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (api *API) ReadVpcGcpPeeringWithVpcId(vpcID string, sleep, timeout int) (
	map[string]interface{}, error) {

	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering", vpcID)
	_, data, err := api.readVpcGcpPeeringWithRetry(path, 1, sleep, timeout)
	return data, err
}

// UpdateVpcGcpPeeringWithVpcId: updates the VPC peering from the API
func (api *API) UpdateVpcGcpPeeringWithVpcId(vpcID string, sleep, timeout int) (
	map[string]interface{}, error) {

	// NOP just read out the VPC peering
	return api.ReadVpcGcpPeeringWithVpcId(vpcID, sleep, timeout)
}

// RemoveVpcGcpPeeringWithVpcId: removes the VPC peering from the API
func (api *API) RemoveVpcGcpPeeringWithVpcId(vpcID, peerID string) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/vpcs/%s/vpc-peering/%s", vpcID, peerID)
	)

	log.Printf("[DEBUG] go-api::vpc_gcp_peering_withvpcid::remove vpc id: %s, peering id: %s",
		vpcID, peerID)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("remove VPC peering failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

// ReadVpcGcpInfoWithVpcId: reads the VPC info from the API
func (api *API) ReadVpcGcpInfoWithVpcId(vpcID string, sleep, timeout int) (
	map[string]interface{}, error) {

	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/info", vpcID)
	_, data, err := api.readVpcGcpPeeringWithRetry(path, 1, sleep, timeout)
	return data, err
}
