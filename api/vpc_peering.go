package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (api *API) waitForPeeringStatus(instanceID int, peeringID string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::vpc_peering::waitForPeeringStatus instance id: %v, peering id: %v", instanceID, peeringID)
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		path := fmt.Sprintf("/api/instances/%v/vpc-peering/status/%v", instanceID, peeringID)
		response, err := api.sling.New().Path(path).Receive(&data, &failed)

		if err != nil {
			return nil, err
		}
		if response.StatusCode != 200 {
			return nil, fmt.Errorf("waitForPeeringStatus failed, status: %v, message: %s", response.StatusCode, failed)
		}
		switch data["status"] {
		case "active", "pending-acceptance":
			return data, nil
		}
		time.Sleep(10 * time.Second)
	}
}

func (api *API) ReadVpcInfo(instanceID int) (map[string]interface{}, error) {
	// Initiale values, 5 attempts and 20 second sleep
	return api.readVpcInfo(instanceID, 5, 20)
}

func (api *API) readVpcInfo(instanceID, attempts, sleep int) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering::info instance id: %v", instanceID)
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/info", instanceID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_peering::info data: %v", data)

	if err != nil {
		return nil, err
	}

	statusCode := response.StatusCode
	log.Printf("[DEBUG] go-api::vpc_peering::info statusCode: %d", statusCode)
	switch {
	case statusCode == 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			if attempts--; attempts > 0 {
				log.Printf("[INFO] go-api::vpc_peering::info Timeout talking to backend "+
					"attempts left %d and retry in %d seconds", attempts, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				return api.readVpcInfo(instanceID, attempts, 2*sleep)
			} else {
				return nil, fmt.Errorf("ReadInfo failed, status: %v, message: %s", response.StatusCode, failed)
			}
		}
	}
	return data, nil
}

func (api *API) ReadVpcPeeringRequest(instanceID int, peeringID string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering::request instance id: %v, peering id: %v", instanceID, peeringID)
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/request/%v", instanceID, peeringID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_peering::request data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadRequest failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, nil
}

func (api *API) AcceptVpcPeering(instanceID int, peeringID string) (map[string]interface{}, error) {
	_, err := api.waitForPeeringStatus(instanceID, peeringID)

	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering::accept instance id: %v, peering id: %v", instanceID, peeringID)
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/request/%v", instanceID, peeringID)
	response, err := api.sling.New().Put(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_peering::accept data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("AcceptVpcPeering failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, nil
}

func (api *API) RemoveVpcPeering(instanceID int, peeringID string) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering::remove instance id: %v, peering id: %v", instanceID, peeringID)
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/%v", instanceID, peeringID)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return fmt.Errorf("RemoveVpcPeering failed, status: %v, message: %s", response.StatusCode, failed)
	}
	return nil
}
