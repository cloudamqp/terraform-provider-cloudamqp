package api

import (
	"fmt"
	"log"
	"strconv"
)

func (api *API) ListInstances() ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	path := fmt.Sprintf("api/instances")
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
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/vpcs")
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc::list data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ListVpcs failed, status: %v, message: %v", response.StatusCode, failed)
	}

	for k, _ := range data {
		vpcID := strconv.FormatFloat(data[k]["id"].(float64), 'f', 0, 64)
		data_temp, _ := api.readVpcName(vpcID)
		data[k]["vpc_name"] = data_temp["name"]
	}

	return data, nil
}
