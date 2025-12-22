package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/node"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ListNodes - list all nodes of the cluster
func (api *API) ListNodes(ctx context.Context, instanceID int64) ([]model.NodeResponse, error) {
	var (
		data   []model.NodeResponse
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/nodes", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ListNodes",
		resourceName: "Node",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return []model.NodeResponse{}, fmt.Errorf("failed to list nodes: %w", err)
	}
	return data, nil
}

// ReadNode - read out node information of a single node
func (api *API) ReadNode(ctx context.Context, instanceID int64, nodeName string) (model.NodeResponse, error) {
	response, err := api.ListNodes(ctx, instanceID)
	if err != nil {
		return model.NodeResponse{}, err
	}

	for i := range response {
		if response[i].Name == nodeName {
			return response[i], nil
		}
	}
	return model.NodeResponse{}, fmt.Errorf("node %s not found", nodeName)
}

// PostAction - request an action for the node(s) (e.g. start/stop/restart RabbitMQ)
func (api *API) PostAction(ctx context.Context, instanceID int64, nodes []string, action string, sleep time.Duration) error {
	var (
		data          map[string]any
		failed        map[string]any
		actionAsRoute string
	)

	params := model.NodeActionRequest{
		Nodes: nodes,
	}

	// Convert action to route format
	switch action {
	case "mgmt.restart":
		actionAsRoute = "mgmt-restart"
	case "cluster.restart":
		actionAsRoute = "cluster-restart"
	case "cluster.stop":
		actionAsRoute = "cluster-stop"
	case "cluster.start":
		actionAsRoute = "cluster-start"
	default:
		actionAsRoute = action
	}

	path := fmt.Sprintf("api/instances/%d/actions/%s", instanceID, actionAsRoute)
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "PostAction",
		resourceName: "NodeAction",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return fmt.Errorf("failed to invoke action %s: %w", action, err)
	}

	return api.PollNodeAction(ctx, instanceID, nodes, action, sleep)
}

func (api *API) PollNodeAction(ctx context.Context, instanceID int64, nodeNames []string, action string, sleep time.Duration) error {
	_, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context has no deadline")
	}

	// Determine expected running state based on action
	var expectedRunning bool
	switch action {
	case "start", "restart", "reboot", "mgmt.restart", "cluster.restart", "cluster.start":
		expectedRunning = true
	case "stop", "cluster.stop":
		expectedRunning = false
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	ticker := time.NewTicker(sleep)
	defer ticker.Stop()

	for {
		nodes, err := api.ListNodes(ctx, instanceID)
		if err != nil {
			return err
		}

		// Create a map for quick lookup
		nodeMap := make(map[string]model.NodeResponse)
		for _, node := range nodes {
			nodeMap[node.Name] = node
		}

		// Check if all target nodes have reached the expected state
		allReady := true
		for _, nodeName := range nodeNames {
			node, exists := nodeMap[nodeName]
			if !exists {
				return fmt.Errorf("node %s not found", nodeName)
			}
			if node.Running != expectedRunning {
				allReady = false
				tflog.Debug(ctx, fmt.Sprintf("node %s not ready: running=%t, expected=%t", nodeName, node.Running, expectedRunning))
				break
			}
		}

		if allReady {
			tflog.Debug(ctx, "all nodes reached expected state")
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while polling for node action completion")
		case <-ticker.C:
			continue
		}
	}
}
