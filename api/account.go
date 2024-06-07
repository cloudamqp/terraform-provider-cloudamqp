package api

import (
	"fmt"
	"log"
	"strconv"
)

func (api *API) ListInstances() ([]map[string]interface{}, error) {
	var (
		data   []map[string]interface{}
		failed map[string]interface{}
		path   = "api/instances"
	)

	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::account::list_instances data: %v", data)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ListInstances failed, status: %v, message: %s", response.StatusCode, failed)
	}
	return data, nil
}

func (api *API) ListVpcs() ([]map[string]interface{}, error) {
	var (
		data   []map[string]interface{}
		failed map[string]interface{}
		path   = "/api/vpcs"
	)

	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc::list data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ListVpcs failed, status: %v, message: %v", response.StatusCode, failed)
	}

	for k := range data {
		vpcID := strconv.FormatFloat(data[k]["id"].(float64), 'f', 0, 64)
		data_temp, _ := api.readVpcName(vpcID)
		data[k]["vpc_name"] = data_temp["name"]
	}

	return data, nil
}

func (api *API) RotatePassword(instanceID int) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("api/instances/%d/account/rotate-password", instanceID)
	)

	response, err := api.sling.New().Post(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to rotate api key, statusCode: %v, failed: %v",
			response.StatusCode, failed)
	}
}

func (api *API) RotateApiKey(instanceID int) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("api/instances/%d/account/rotate-apikey", instanceID)
	)

	response, err := api.sling.New().Post(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to rotate api key, statusCode: %v, failed: %v",
			response.StatusCode, failed)
	}
}
