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
			}
			return nil, fmt.Errorf("ReadInfo failed, status: %v, message: %s", response.StatusCode, failed)
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

func (api *API) retryAcceptVpcPeeringWithVpcId(vpcID, peeringID string, attempt, sleep, timeout int) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/request/%s", vpcID, peeringID)
	response, err := api.sling.New().Put(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if attempt*sleep >= timeout {
		return nil, fmt.Errorf("AcceptVpcPeeringWithVpcId failed, reached timeout of %d seconds", timeout)
	} else if response.StatusCode == 400 {
		errorCode := failed["error_code"].(float64)
		if errorCode == 40001 {
			log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::accept firewall not finished configuring will retry "+
				"accept VPC peering, attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.retryAcceptVpcPeeringWithVpcId(vpcID, peeringID, attempt, sleep, timeout)
		}
		return nil, fmt.Errorf("AcceptVpcPeeringWithVpcId failed, status: %v, message: %s", response.StatusCode, failed)
	} else if response.StatusCode != 200 {
		return nil, fmt.Errorf("AcceptVpcPeeringWithVpcId failed, status: %v, message: %s", response.StatusCode, failed)
	}
	return data, nil
}

func (api *API) AcceptVpcPeeringWithVpcId(vpcID, peeringID string, sleep, timeout int) (map[string]interface{}, error) {
	_, err := api.waitForPeeringStatusWithVpcId(vpcID, peeringID)
	if err != nil {
		return nil, err
	}
	return api.retryAcceptVpcPeeringWithVpcId(vpcID, peeringID, 1, sleep, timeout)
}

func (api *API) retryRemoveVpcPeeringWithVpcId(vpcID, peeringID string, attempt, sleep, timeout int) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::remove vpc id: %s, peering id: %s", vpcID, peeringID)
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/%s", vpcID, peeringID)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if attempt*sleep >= timeout {
		return fmt.Errorf("RemoveVpcPeeringWithVpcId failed, reached timeout of %d seconds", timeout)
	} else if response.StatusCode == 400 {
		errorCode := failed["error_code"].(float64)
		if errorCode == 40001 {
			log.Printf("[DEBUG] go-api::vpc_peering_withvpcid::remove firewall not finished configuring will retry "+
				"accept VPC peering, attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.retryRemoveVpcPeeringWithVpcId(vpcID, peeringID, attempt, sleep, timeout)
		}
		return fmt.Errorf("RemoveVpcPeeringWithVpcId failed, status: %v, message: %s", response.StatusCode, failed)
	} else if response.StatusCode != 204 {
		return fmt.Errorf("RemoveVpcPeeringWithVpcId failed, status: %v, message: %s", response.StatusCode, failed)
	}
	return nil
}

func (api *API) RemoveVpcPeeringWithVpcId(vpcID, peeringID string, sleep, timeout int) error {
	return api.retryRemoveVpcPeeringWithVpcId(vpcID, peeringID, 1, sleep, timeout)
}
