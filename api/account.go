package api

import (
	"fmt"
	"log"
)

// ReadAccount - read out node information of the cluster
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
