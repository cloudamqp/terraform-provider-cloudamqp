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
func (api *API) ReadNode(instanceID int, nodeName string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::nodes::read_node instance id: %d node name: %s", instanceID, nodeName)
	response, err := api.ReadNodes(instanceID)
	if err != nil {
		return nil, err
	}
	for i := range response {
		if response[i]["name"] == nodeName {
			data = response[i]
			break
		}
	}
	return data, nil
}

// PostAction - request an action for the node (e.g. start/stop/restart RabbitMQ)
func (api *API) PostAction(instanceID int, nodeName string, action string) (map[string]interface{}, error) {
	var actionAsRoute string
	params := make(map[string][]string)
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	if action == "mgmt.restart" {
		actionAsRoute = "mgmt-restart"
	} else {
		actionAsRoute = action
	}
	params["nodes"] = append(params["nodes"], nodeName)
	path := fmt.Sprintf("api/instances/%d/actions/%s", instanceID, actionAsRoute)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("action failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitOnNodeAction(instanceID, nodeName, action)
}

func (api *API) waitOnNodeAction(instanceID int, nodeName string, action string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::nodes::waitOnNodeAction waiting")
	for {
		data, err := api.ReadNode(instanceID, nodeName)

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
		time.Sleep(20 * time.Second)
	}
}
