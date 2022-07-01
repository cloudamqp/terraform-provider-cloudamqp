package api

import (
	"fmt"
	"log"
	"strconv"
)

func (api *API) ResizeDisk(instanceID int, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	id := strconv.Itoa(instanceID)
	log.Printf("[DEBUG] go-api::disk::resize instance ID: %s", id)
	path := fmt.Sprintf("api/instances/%s/disk", id)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ResizeDisk failed, status: %v, message: %s", response.StatusCode, failed)
	}

	if err = api.waitUntilAllNodesReady(id); err != nil {
		return nil, err
	}
	return data, nil
}
