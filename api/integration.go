package api

import (
	"fmt"
	"log"
	"strconv"
)

// CreateIntegration enables integration communication, either for logs or metrics.
func (api *API) CreateIntegration(instanceID int, intType string, intName string, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::integration::create params: %v", params)
	path := fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, intName)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::integration::create response data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 201 {
		return nil, fmt.Errorf(fmt.Sprintf("CreateIntegration failed, status: %v, message: %s", response.StatusCode, failed))
	}

	if v, ok := data["id"]; ok {
		data["id"] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
	} else {
		msg := fmt.Sprintf("go-api::integration::create Invalid integration identifier: %v", data["id"])
		log.Printf("[ERROR] %s", msg)
		return nil, fmt.Errorf(msg)
	}

	return data, err
}

// ReadIntegration retrieves a specific logs or metrics integration
func (api *API) ReadIntegration(instanceID int, intType, id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::integration::read instance id: %d, int_type: %s, id: %s", instanceID, intType, id)
	path := fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, id)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf(fmt.Sprintf("ReadIntegration failed, status: %v, message: %s", response.StatusCode, failed))
	}

	// Convert API response body, config part, into single map
	convertedData := make(map[string]interface{})
	for k, v := range data {
		if k == "id" {
			convertedData[k] = v
		} else if k == "type" {
			convertedData[k] = v
		} else if k == "config" {
			for configK, configV := range data["config"].(map[string]interface{}) {
				convertedData[configK] = configV
			}
		}
	}
	return convertedData, err
}

// UpdateIntegration updated the integration with new information
func (api *API) UpdateIntegration(instanceID int, intType, id string, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBIG] go-api::integration::update instance id: %d, intType: %s, id: %s", instanceID, intType, id)
	path := fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, id)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if response.StatusCode != 204 {
		return fmt.Errorf(fmt.Sprintf("UpdateIntegration failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}

// DeleteIntegration removes log or metric integration.
func (api *API) DeleteIntegration(instanceID int, intType, id string) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::integration::delete instance id: %d, intType: %s, id: %s", instanceID, intType, id)
	path := fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, id)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if response.StatusCode != 204 {
		return fmt.Errorf(fmt.Sprintf("DeleteNotificaion failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}
