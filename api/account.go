package api

import (
	"fmt"
	"log"
	"strconv"
)

func (api *API) ListInstances() ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Path("api/instances").Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("[DEBUG] api::account#list_instances data: %v", data)
		return data, nil
	case 410:
		log.Printf("[WARN] api::instance#list status: 410, message: The instance has been deleted")
		return nil, nil
	default:
		return nil, fmt.Errorf("list instances failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) ListVpcs() ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = "/api/vpcs"
	)

	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("[DEBUG] api::account#list_vpcs data: %v", data)
		for k := range data {
			vpcID := strconv.FormatFloat(data[k]["id"].(float64), 'f', 0, 64)
			data_temp, _ := api.readVpcName(vpcID)
			data[k]["vpc_name"] = data_temp["name"]
		}
		return data, nil
	default:
		return nil, fmt.Errorf("list VPCs failed, status: %d, message: %s", response.StatusCode, failed)
	}
}

func (api *API) RotatePassword(instanceID int) error {
	var (
		failed map[string]any
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
		return fmt.Errorf("failed to rotate api key, statusCode: %d, failed: %v",
			response.StatusCode, failed)
	}
}

func (api *API) RotateApiKey(instanceID int) error {
	var (
		failed map[string]any
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
