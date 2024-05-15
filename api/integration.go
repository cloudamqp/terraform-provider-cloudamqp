package api

import (
	"fmt"
	"log"
	"strconv"
)

// CreateIntegration enables integration communication, either for logs or metrics.
func (api *API) CreateIntegration(instanceID int, intType string, intName string,
	params map[string]interface{}) (map[string]interface{}, error) {

	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, intName)
	)

	log.Printf("[DEBUG] api::integration#create path: %s, params: %v", path, params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 201:
		log.Printf("[DEBUG] api::integration#create response data: %v", data)
		if v, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		} else {
			return nil, fmt.Errorf("create integration failed, invalid integration identifier: %v",
				data["id"])
		}
		return data, err
	default:
		return nil, fmt.Errorf("create integration failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

// ReadIntegration retrieves a specific logs or metrics integration
func (api *API) ReadIntegration(instanceID int, intType, intID string) (
	map[string]interface{}, error) {

	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, intID)
	)

	log.Printf("[DEBUG] api::integration#read path: %s", path)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("[DEBUG] api::integration#read data: %v", data)
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
		log.Printf("[DEBUG] api::integration#read convertedDatat: %v", convertedData)
		return convertedData, err
	default:
		return nil, fmt.Errorf("read integration failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

// UpdateIntegration updated the integration with new information
func (api *API) UpdateIntegration(instanceID int, intType, intID string,
	params map[string]interface{}) error {

	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, intID)
	)

	log.Printf("[DEBUG] api::integration#update path: %s", path)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("update integration failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

// DeleteIntegration removes log or metric integration.
func (api *API) DeleteIntegration(instanceID int, intType, intID string) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, intID)
	)

	log.Printf("[DEBUG] api::integration#delete path: %s", path)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("delete notification failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}
