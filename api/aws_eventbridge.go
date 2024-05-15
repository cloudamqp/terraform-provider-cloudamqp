package api

import (
	"fmt"
	"log"
	"strconv"
)

func (api *API) CreateAwsEventBridge(instanceID int, params map[string]interface{}) (
	map[string]interface{}, error) {

	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/eventbridges", instanceID)
	)

	log.Printf("[DEBUG] api::aws-eventbridge#create path: %s", path)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 201:
		if id, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
			log.Printf("[DEBUG] api::aws-eventbridge#create AWS EventBridge identifier: %v", data["id"])
		} else {
			return nil, fmt.Errorf("failed to create AWS EventBridge, invalid identifier: %v", data["id"])
		}
		return data, nil
	default:
		return nil, fmt.Errorf("failed to create AWS EventBridge, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) ReadAwsEventBridge(instanceID int, eventbridgeID string) (
	map[string]interface{}, error) {

	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/eventbridges/%s", instanceID, eventbridgeID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read AWS EventBridge, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) ReadAwsEventBridges(instanceID int) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/eventbridges", instanceID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read AWS EventBridges, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) DeleteAwsEventBridge(instanceID int, eventbridgeID string) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/eventbridges/%s", instanceID, eventbridgeID)
	)

	log.Printf("[DEBUG] api::aws-eventbridge#delete path: %s", path)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	case 404:
		// AWS EventBridge not found in the backend. Silent let the resource be deleted.
		return nil
	default:
		return fmt.Errorf("failed to delete AWS EventBridge, status: %d, message: %s",
			response.StatusCode, failed)
	}
}
