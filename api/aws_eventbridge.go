package api

import (
	"fmt"
	"log"
	"strconv"
)

func (api *API) CreateAwsEventBridge(instanceID int, params map[string]interface{}) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/eventbridges", instanceID)
	)

	log.Printf("[DEBUG] go-api::aws-eventbridge::create instance ID: %d, params: %v", instanceID, params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 201 {
		return nil, fmt.Errorf("Failed to create AWS EventBridge, status: %v, message: %s",
			response.StatusCode, failed)
	}
	if id, ok := data["id"]; ok {
		data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
		log.Printf("[DEBUG] go-api::aws-eventbridge::create EventBridge identifier: %v", data["id"])
	} else {
		return nil, fmt.Errorf("go-api::aws-eventbridge::create Invalid identifier: %v", data["id"])
	}

	return data, nil
}

func (api *API) ReadAwsEventBridge(instanceID int, eventbridgeID string) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/eventbridges/%s", instanceID, eventbridgeID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to read AWS EventBridge, status: %v, message: %s",
			response.StatusCode, failed)
	}

	return data, nil
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
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to read AWS EventBridges, status: %v, message: %s",
			response.StatusCode, failed)
	}

	return data, nil
}

func (api *API) DeleteAwsEventBridge(instanceID int, eventbridgeID string) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/eventbridges/%s", instanceID, eventbridgeID)
	)

	log.Printf("[DEBUG] go-api::aws-eventbridge::delete instance id: %d, eventbridge id: %s", instanceID, eventbridgeID)
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
	}

	return fmt.Errorf("Failed to delete AWS EventBridge, status: %v, message: %s", response.StatusCode, failed)
}
