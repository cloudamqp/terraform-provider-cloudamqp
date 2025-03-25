package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ListNodes - list all nodes of the cluster
func (api *API) ListNodes(ctx context.Context, instanceID int) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/nodes", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("request path: %s", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("failed to list nodes, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

// ReadNode - read out node information of a single node
func (api *API) ReadNode(ctx context.Context, instanceID int, nodeName string) (
	map[string]any, error) {

	var (
		data map[string]any
	)

	response, err := api.ListNodes(ctx, instanceID)
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
func (api *API) PostAction(ctx context.Context, instanceID int, nodeName string, action string) (
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
	tflog.Debug(ctx, fmt.Sprintf("request path: %s", path))
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return api.waitOnNodeAction(ctx, instanceID, nodeName, action)
	default:

		return nil, fmt.Errorf("failed to invoke action %s, status: %d, message: %s",
			action, response.StatusCode, failed)
	}
}

func (api *API) waitOnNodeAction(ctx context.Context, instanceID int, nodeName string, action string) (
	map[string]any, error) {

	tflog.Debug(ctx, "waiting on node action")
	for {
		data, err := api.ReadNode(ctx, instanceID, nodeName)
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
