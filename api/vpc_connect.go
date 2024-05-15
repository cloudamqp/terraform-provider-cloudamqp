package api

import (
	"fmt"
	"log"
	"time"
)

// EnableVpcConnect: Enable VPC Connect and wait until finished.
// Need to enable VPC for an instance, if no standalone VPC used.
// Wait until finished with configureable sleep and timeout.
func (api *API) EnableVpcConnect(instanceID int, params map[string][]interface{},
	sleep, timeout int) error {

	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	)

	if err := api.EnableVPC(instanceID); err != nil {
		return err
	}

	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return api.waitForEnableVpcConnectWithRetry(instanceID, 1, sleep, timeout)
	default:
		return fmt.Errorf("enable VPC Connect failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

// ReadVpcConnect: Reads VPC Connect information
func (api *API) ReadVpcConnect(instanceID int) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("read VPC Connect failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

// UpdateVpcConnect: Update allowlist for the VPC Connect
func (api *API) UpdateVpcConnect(instanceID int, params map[string][]interface{}) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	)

	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("update VPC connect failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

// DisableVpcConnect: Disable the VPC Connect feature
func (api *API) DisableVpcConnect(instanceID int) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	)

	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("disable VPC Connect failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

// waitForEnableVpcConnectWithRetry: Wait until status change from pending to enable
func (api *API) waitForEnableVpcConnectWithRetry(instanceID, attempt, sleep, timeout int) error {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("enable VPC Connect failed, reached timeout of %d seconds", timeout)
	}
	log.Printf("[DEBUG] VPC-Connect: waitForEnableVpcConnectWithRetry data: %v", data)

	switch response.StatusCode {
	case 200:
		switch data["status"].(string) {
		case "enabled":
			return nil
		case "pending":
			log.Printf("[DEBUG] go-api::vpc-connect::enable not finished and will retry, "+
				"attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.waitForEnableVpcConnectWithRetry(instanceID, attempt, sleep, timeout)
		}
	}

	return fmt.Errorf("wait for enable VPC Connect failed, status: %v, message: %s",
		response.StatusCode, failed)
}

// enableVPC: Enable VPC for an instance
// Check if the instance already have a standalone VPC
func (api *API) EnableVPC(instanceID int) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/vpc", instanceID)
	)

	data, _ := api.ReadInstance(fmt.Sprintf("%d", instanceID))
	if data["vpc"] == nil {
		response, err := api.sling.New().Put(path).Receive(nil, &failed)
		if err != nil {
			return err
		}

		switch response.StatusCode {
		case 200:
			log.Printf("[DEBUG] VPC-Connect: VPC features enabled")
			return nil
		default:
			return fmt.Errorf("enable VPC failed, status: %v, message: %s",
				response.StatusCode, failed)
		}
	}

	log.Printf("[DEBUG] VPC-Connect: VPC features already enabled")
	return nil
}
