package api

// VPC peering for AWS, using vpcID as identifier.

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (api *API) waitForPeeringStatusWithVpcId(vpcID, peeringID string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::waitForPeeringStatus instance id: %s, peering id: %s", vpcID, peeringID)
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		time.Sleep(10 * time.Second)
		path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/status/%s", vpcID, peeringID)
		response, err := api.sling.New().Path(path).Receive(&data, &failed)

		if err != nil {
			return nil, err
		}
		if response.StatusCode != 200 {
			return nil, fmt.Errorf("Wait for peering status failed, status: %v, message: %s", response.StatusCode, failed)
		}
		switch data["status"] {
		case "active", "pending-acceptance":
			return data, nil
		}
	}
}

func (api *API) ReadVpcInfoWithVpcId(vpcID string) (map[string]interface{}, error) {
	// Initiale values, 5 attempts and 20 second sleep
	return api.readVpcInfoWithRetryWithVpcId(vpcID, 5, 20)
}

func (api *API) readVpcInfoWithRetryWithVpcId(vpcID string, attempts, sleep int) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::info vpc id: %s", vpcID)
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/info", vpcID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::info data: %v", data)

	if err != nil {
		return nil, err
	}

	statusCode := response.StatusCode
	log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::info statusCode: %d", statusCode)
	switch {
	case statusCode == 400:
		// Todo: Implement error code to be checked instead. To avoid using string comparison.
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			if attempts--; attempts > 0 {
				log.Printf("[INFO] go-api::vpc_peering_withvpcid::info Timeout talking to backend "+
					"attempts left %d and retry in %d seconds", attempts, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				return api.readVpcInfoWithRetryWithVpcId(vpcID, attempts, 2*sleep)
			} else {
				return nil, fmt.Errorf("ReadInfo failed, status: %v, message: %s", response.StatusCode, failed)
			}
		}
	}
	return data, nil
}

func (api *API) ReadVpcPeeringRequestWithVpcId(vpcID, peeringID string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::request vpc id: %v, peering id: %v", vpcID, peeringID)
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/request/%s", vpcID, peeringID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::request data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadRequest failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, nil
}

func (api *API) AcceptVpcPeeringWithVpcId(vpcID, peeringID string) (map[string]interface{}, error) {
	_, err := api.waitForPeeringStatusWithVpcId(vpcID, peeringID)

	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::accept vpc id: %s, peering id: %s", vpcID, peeringID)
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/request/%s", vpcID, peeringID)
	response, err := api.sling.New().Put(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::accept data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("AcceptVpcPeering failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, nil
}

func (api *API) RemoveVpcPeeringWithVpcId(vpcID, peeringID string) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::remove vpc id: %s, peering id: %s", vpcID, peeringID)
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/%s", vpcID, peeringID)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return fmt.Errorf("Remove VPC peering failed, status: %v, message: %s", response.StatusCode, failed)
	}
	return nil
}
