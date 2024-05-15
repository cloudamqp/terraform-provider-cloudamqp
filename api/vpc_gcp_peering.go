package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// waitForGcpPeeringStatus: waits for the VPC peering status to be ACTIVE or until timed out
func (api *API) waitForGcpPeeringStatus(path, peerID string,
	attempt, sleep, timeout int) error {

	var (
		data map[string]interface{}
		err  error
	)

	for {
		if attempt*sleep > timeout {
			return fmt.Errorf("wait until GCP VPC peering status reached timeout of %d seconds", timeout)
		}

		attempt, data, err = api.readVpcGcpPeeringWithRetry(path, attempt, sleep, timeout)
		if err != nil {
			return err
		}

		rows := data["rows"].([]interface{})
		if len(rows) > 0 {
			for _, row := range rows {
				tempRow := row.(map[string]interface{})
				if tempRow["name"] != peerID {
					continue
				}
				if tempRow["state"] == "ACTIVE" {
					return nil
				}
			}
		}
		log.Printf("[INFO] go-api::vpc_gcp_peering::waitForGcpPeeringStatus Waiting for state = ACTIVE "+
			"attempt %d until timeout: %d", attempt, (timeout - (attempt * sleep)))
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

// RequestVpcGcpPeering: requests a VPC peering from an instance.
func (api *API) RequestVpcGcpPeering(instanceID int, params map[string]interface{},
	waitOnStatus bool, sleep, timeout int) (map[string]interface{}, error) {

	path := fmt.Sprintf("api/instances/%v/vpc-peering", instanceID)
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

// requestVpcGcpPeeringWithRetry: requests a VPC peering from a path with retry logic
func (api *API) requestVpcGcpPeeringWithRetry(path string, params map[string]interface{},
	waitOnStatus bool, attempt, sleep, timeout int) (int, map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
	)

	log.Printf("[DEBUG] go-api::vpc_gcp_peering::request path: %s, params: %v", path, params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return attempt, nil, err
	} else if attempt*sleep > timeout {
		return attempt, nil,
			fmt.Errorf("request VPC peering failed, reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 200:
		return attempt, data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[INFO] go-api::vpc_gcp_peering::request Timeout talking to backend "+
				"attempt %d until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.requestVpcGcpPeeringWithRetry(path, params, waitOnStatus, attempt, sleep, timeout)
		}
	}
	return attempt, nil, fmt.Errorf("request VPC peering failed, status: %v, message: %s",
		response.StatusCode, failed)
}

// ReadVpcGcpPeering: reads the VPC peering from the API
func (api *API) ReadVpcGcpPeering(instanceID, sleep, timeout int) (
	map[string]interface{}, error) {

	path := fmt.Sprintf("/api/instances/%v/vpc-peering", instanceID)
	_, data, err := api.readVpcGcpPeeringWithRetry(path, 1, sleep, timeout)
	return data, err
}

// readVpcGcpPeeringWithRetry: reads the VPC peering from the API with retry logic
func (api *API) readVpcGcpPeeringWithRetry(path string, attempt, sleep, timeout int) (
	int, map[string]interface{}, error) {

	var (
		data   map[string]interface{}
		failed map[string]interface{}
	)

	log.Printf("[DEBUG] go-api::vpc_gcp_peering::read path: %s", path)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return attempt, nil, err
	} else if attempt*sleep > timeout {
		return attempt, nil, fmt.Errorf("read VPC peering reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 200:
		return attempt, data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[INFO] go-api::vpc_gcp_peering::read Timeout talking to backend "+
				"attempt %d until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readVpcGcpPeeringWithRetry(path, attempt, sleep, timeout)
		}
	}
	return attempt, nil, fmt.Errorf("read VPC peering with retry failed, status: %v, message: %s",
		response.StatusCode, failed)
}

// UpdateVpcGcpPeering: updates a VPC peering from an instance.
func (api *API) UpdateVpcGcpPeering(instanceID int, sleep, timeout int) (
	map[string]interface{}, error) {

	// NOP just read out the VPC peering
	return api.ReadVpcGcpPeering(instanceID, sleep, timeout)
}

// RemoveVpcGcpPeering: removes a VPC peering from an instance.
func (api *API) RemoveVpcGcpPeering(instanceID int, peerID string) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%v/vpc-peering/%v", instanceID, peerID)
	)

	log.Printf("[DEBUG] go-api::vpc_gcp_peering::remove instance id: %v, peering id: %v", instanceID, peerID)
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

// ReadVpcGcpInfo: reads the VPC info from the API
func (api *API) ReadVpcGcpInfo(instanceID, sleep, timeout int) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/info", instanceID)
	return api.readVpcGcpInfoWithRetry(path, 1, sleep, timeout)
}

// readVpcGcpInfoWithRetry: reads the VPC info from the API with retry logic
func (api *API) readVpcGcpInfoWithRetry(path string, attempt, sleep, timeout int) (
	map[string]interface{}, error) {

	var (
		data   map[string]interface{}
		failed map[string]interface{}
	)

	log.Printf("[DEBUG] go-api::vpc_gcp_peering::info path: %s", path)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("read VPC info, reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[INFO] go-api::vpc_gcp_peering::info Timeout talking to backend "+
				"attempt %d until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readVpcGcpInfoWithRetry(path, attempt, sleep, timeout)
		}
	}
	return nil, fmt.Errorf("read VPC info failed, status: %v, message: %s",
		response.StatusCode, failed)
}
