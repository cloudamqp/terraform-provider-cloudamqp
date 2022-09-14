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

	switch {
	case response.StatusCode == 200:
		if err = api.waitUntilAllNodesReady(id); err != nil {
			return nil, err
		}
		return data, nil
	case response.StatusCode == 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40002:
			return nil, fmt.Errorf("Resize disk failed: %s", failed["error"].(string))
		}
	}
	return nil, fmt.Errorf("Resize disk failed, status: %v, message: %s", response.StatusCode, failed)
}
