package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (api *API) AcceptVpcPeering(instanceID int, peeringID string, sleep, timeout int) (map[string]interface{}, error) {
	attempt, err := api.waitForPeeringStatus(instanceID, peeringID, 1, sleep, timeout)
	log.Printf("[DEBUG] go-api::vpc_peering::accept attempt: %d, sleep: %d, timeout: %d", attempt, sleep, timeout)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/request/%v", instanceID, peeringID)
	return api.retryAcceptVpcPeering(path, attempt, sleep, timeout)
}

func (api *API) ReadVpcInfo(instanceID int) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/info", instanceID)
	// Initiale values, 5 attempts and 20 second sleep
	return api.readVpcInfoWithRetry(path, 5, 20)
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
	} else if response.StatusCode != 200 {
		return nil, fmt.Errorf("readRequest failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, nil
}

func (api *API) RemoveVpcPeering(instanceID int, peeringID string, sleep, timeout int) error {
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/%v", instanceID, peeringID)
	return api.retryRemoveVpcPeering(path, 1, sleep, timeout)
}

func (api *API) retryAcceptVpcPeering(path string, attempt, sleep, timeout int) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::vpc_peering::retryRemoveVpcPeering path: %s, "+
		"attempt: %d, sleep: %d, timeout: %d", path, attempt, sleep, timeout)
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	response, err := api.sling.New().Put(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("accept VPC peering failed, reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001:
			log.Printf("[DEBUG] go-api::vpc_peering::accept firewall not finished configuring will retry "+
				"accept VPC peering, attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.retryAcceptVpcPeering(path, attempt, sleep, timeout)
		}
	}

	return nil, fmt.Errorf("accept VPC peering failed, status: %v, message: %s", response.StatusCode, failed)
}

func (api *API) readVpcInfoWithRetry(path string, attempts, sleep int) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::vpc_peering::readVpcInfoWithRetry path: %s, "+
		"attempts: %d, sleep: %d", path, attempts, sleep)
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_peering::info data: %v", data)

	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			if attempts--; attempts > 0 {
				log.Printf("[INFO] go-api::vpc_peering::info Timeout talking to backend "+
					"attempts left %d and retry in %d seconds", attempts, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				return api.readVpcInfoWithRetry(path, attempts, 2*sleep)
			}
			return nil, fmt.Errorf("readInfo failed, status: %v, message: %s", response.StatusCode, failed)
		}
	}

	return nil, fmt.Errorf("readInfo failed, status: %v, message: %s", response.StatusCode, failed)
}

func (api *API) retryRemoveVpcPeering(path string, attempt, sleep, timeout int) error {
	log.Printf("[DEBUG] go-api::vpc_peering::retryRemoveVpcPeering path: %s, "+
		"attempt: %d, sleep: %d, timeout: %d", path, attempt, sleep, timeout)
	failed := make(map[string]interface{})
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("remove VPC peering failed, reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 204:
		return nil
	case 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001:
			log.Printf("[DEBUG] go-api::vpc_peering::remove firewall not finished configuring will retry "+
				"removing VPC peering, attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.retryRemoveVpcPeering(path, attempt, sleep, timeout)
		}
	}

	return fmt.Errorf("remove VPC peering failed, status: %v, message: %s", response.StatusCode, failed)
}

func (api *API) waitForPeeringStatus(instanceID int, peeringID string, attempt, sleep, timeout int) (int, error) {
	time.Sleep(10 * time.Second)
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/status/%v", instanceID, peeringID)
	return api.waitForPeeringStatusWithRetry(path, peeringID, attempt, sleep, timeout)
}

func (api *API) waitForPeeringStatusWithRetry(path, peeringID string, attempt, sleep, timeout int) (int, error) {
	log.Printf("[DEBUG] go-api::vpc_peering::waitForPeeringStatusWithRetry path: %s "+
		"attempt: %d, sleep: %d, timeout: %d", path, attempt, sleep, timeout)
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	response, err := api.sling.New().Path(path).Receive(&data, &failed)

	if err != nil {
		return attempt, err
	} else if attempt*sleep > timeout {
		return attempt, fmt.Errorf("accept VPC peering failed, reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 200:
		switch data["status"] {
		case "active", "pending-acceptance":
			return attempt, nil
		case "deleted":
			return attempt, fmt.Errorf("peering: %s has been deleted", peeringID)
		}
	case 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40003:
			log.Printf("[DEBUG] go-api::vpc_peering::waitForPeeringStatusWithRetry %s, attempt: %d, until timeout: %d",
				failed["message"].(string), attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.waitForPeeringStatusWithRetry(path, peeringID, attempt, sleep, timeout)
		}
	}

	return attempt, fmt.Errorf("accept VPC peering failed, status: %v, message: %v", response.StatusCode, failed)
}
