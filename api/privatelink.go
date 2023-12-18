package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// EnablePrivatelink: Enable PrivateLink and wait until finished.
// Need to enable VPC for an instance, if no standalone VPC used.
// Wait until finished with configureable sleep and timeout.
func (api *API) EnablePrivatelink(instanceID int, params map[string][]interface{},
	sleep, timeout int) error {

	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
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
		return api.waitForEnablePrivatelinkWithRetry(instanceID, 1, sleep, timeout)
	default:
		return fmt.Errorf("enable PrivateLink failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

// ReadPrivatelink: Reads PrivateLink information
func (api *API) ReadPrivatelink(instanceID, sleep, timeout int) (map[string]interface{}, error) {
	return api.readPrivateLinkWithRetry(instanceID, 1, sleep, timeout)
}

func (api *API) readPrivateLinkWithRetry(instanceID, attempt, sleep, timeout int) (
	map[string]interface{}, error) {

	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("read PrivateLink failed, reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[INFO] go-api::privatelink::read Timeout talking to backend "+
				"attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readPrivateLinkWithRetry(instanceID, attempt, sleep, timeout)
		}
	}

	return nil, fmt.Errorf("read PrivateLink failed, status: %v, message: %s",
		response.StatusCode, failed)
}

// UpdatePrivatelink: Update allowed principals or subscriptions
func (api *API) UpdatePrivatelink(instanceID int, params map[string][]interface{}) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("update Privatelink failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

// DisablePrivatelink: Disable the PrivateLink feature
func (api *API) DisablePrivatelink(instanceID int) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("disable Privatelink failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

// waitForEnablePrivatelinkWithRetry: Wait until status change from pending to enable
func (api *API) waitForEnablePrivatelinkWithRetry(instanceID, attempt, sleep, timeout int) error {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("enable PrivateLink failed, reached timeout of %d seconds", timeout)
	}
	log.Printf("[DEBUG] PrivateLink: waitForEnablePrivatelinkWithRetry data: %v", data)

	switch response.StatusCode {
	case 200:
		switch data["status"].(string) {
		case "enabled":
			return nil
		case "pending":
			log.Printf("[DEBUG] go-api::privatelink::enable not finished and will retry, "+
				"attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.waitForEnablePrivatelinkWithRetry(instanceID, attempt, sleep, timeout)
		}
	}

	return fmt.Errorf("wait for enable PrivateLink failed, status: %v, message: %s",
		response.StatusCode, failed)
}
