package api

// VPC peering for AWS, using vpcID as identifier.

import (
	"fmt"
	"log"
	"time"
)

func (api *API) AcceptVpcPeeringWithVpcId(vpcID, peeringID string, sleep, timeout int) (
	map[string]any, error) {

	attempt, err := api.waitForPeeringStatusWithVpcID(vpcID, peeringID, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/request/%s", vpcID, peeringID)
	return api.retryAcceptVpcPeering(path, attempt, sleep, timeout)
}

func (api *API) ReadVpcInfoWithVpcId(vpcID string) (map[string]any, error) {
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/info", vpcID)
	// Initiale values, 5 attempts and 20 second sleep
	return api.readVpcInfoWithRetry(path, 5, 20)
}

func (api *API) ReadVpcPeeringRequestWithVpcId(vpcID, peeringID string) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%s/vpc-peering/request/%s", vpcID, peeringID)
	)

	log.Printf("[DEBUG] api::vpc_peering_withvpcid#request path: %s", path)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("[DEBUG] api::vpc_peering_withvpcid#request data: %v", data)
		return data, nil
	default:
		return nil, fmt.Errorf("read request failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) RemoveVpcPeeringWithVpcId(vpcID, peeringID string, sleep, timeout int) error {
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/%s", vpcID, peeringID)
	return api.retryRemoveVpcPeering(path, 1, sleep, timeout)
}

func (api *API) waitForPeeringStatusWithVpcID(vpcID, peeringID string, attempt, sleep,
	timeout int) (int, error) {

	time.Sleep(10 * time.Second)
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/status/%s", vpcID, peeringID)
	return api.waitForPeeringStatusWithRetry(path, peeringID, attempt, sleep, timeout)
}
