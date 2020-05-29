package api

import (
	"fmt"
	"log"
)

// ReadNodes - read out node information of the cluster
func (api *API) ReadNodes(instanceID int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::nodes::read instance id: %d", instanceID)
	path := fmt.Sprintf("api/instances/%d/nodes", instanceID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf(fmt.Sprintf("ReadNodes failed, status: %v, message: %s", response.StatusCode, failed))
	}
	return data, nil
}
