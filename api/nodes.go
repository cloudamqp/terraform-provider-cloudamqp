package api

import (
	"fmt"
	"log"
	"time"
)

// ReadNodes - read out node information of the cluster
func (api *API) ReadNodes(instanceID int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::nodes::read_nodes instance id: %d", instanceID)
	path := fmt.Sprintf("api/instances/%d/nodes", instanceID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadNodes failed, status: %v, message: %s", response.StatusCode, failed)
	}
	return data, nil
}

// ReadNode - read out node information of a single node
func (api *API) ReadNode(instanceID, nodeID int) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::nodes::read_node instance id: %d, node id: %d", instanceID, nodeID)
	path := fmt.Sprintf("api/instances/%d/nodes/%d", instanceID, nodeID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadNode failed, status: %v, message: %s", response.StatusCode, failed)
	}
	return data, nil
}

// PostAction - request an action for the node (e.g. start/stop/restart RabbitMQ)
func (api *API) PostAction(instanceID, nodeID int, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	path := fmt.Sprintf("api/instances/%d/nodes/%d/action", instanceID, nodeID)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Action failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitOnNodeAction(instanceID, nodeID, params["action"].(string))
}

func (api *API) waitOnNodeAction(instanceID, nodeID int, action string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::nodes::waitOnNodeAction waiting")
	for {
		time.Sleep(20 * time.Second)
		data, err := api.ReadNode(instanceID, nodeID)

		if err != nil {
			return nil, err
		}

		switch action {
		case "start", "restart", "reboot", "mgmt.restart":
			if data["running"] == true {
				return data, nil
			}
		case "stop":
			if data["running"] == false {
				return data, nil
			}
		}
	}
}
