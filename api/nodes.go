package api

import (
	"fmt"
	"log"
	"time"
)

// ListNodes - list all nodes of the cluster
func (api *API) ListNodes(instanceID int) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/nodes", instanceID)
	)

	log.Printf("[DEBUG] api::nodes#list_nodes path: %s", path)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("list nodes failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

// ReadNode - read out node information of a single node
func (api *API) ReadNode(instanceID int, nodeName string) (map[string]any, error) {
	var (
		data map[string]any
	)

	log.Printf("[DEBUG] go-api::nodes#read_node instance id: %d node name: %s", instanceID, nodeName)
	response, err := api.ListNodes(instanceID)
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
func (api *API) PostAction(instanceID int, nodeName string, action string) (
	map[string]any, error) {

	var (
		data          map[string]any
		failed        map[string]any
		actionAsRoute string
		params        = make(map[string][]string)
	)

	params["nodes"] = append(params["nodes"], nodeName)

	if action == "mgmt.restart" {
		actionAsRoute = "mgmt-restart"
	} else {
		actionAsRoute = action
	}
	path := fmt.Sprintf("api/instances/%d/actions/%s", instanceID, actionAsRoute)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return api.waitOnNodeAction(instanceID, nodeName, action)
	default:

		return nil, fmt.Errorf("action %s failed, status: %d, message: %s",
			action, response.StatusCode, failed)
	}
}

func (api *API) waitOnNodeAction(instanceID int, nodeName string, action string) (
	map[string]any, error) {

	log.Printf("[DEBUG] api::nodes#waitOnNodeAction waiting")
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
