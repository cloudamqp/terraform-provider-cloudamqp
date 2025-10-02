package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/network"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) CreateVPC(ctx context.Context, params model.VpcRequest) (model.VpcResponse, error) {
	var (
		data   model.VpcResponse
		failed map[string]any
		path   = "/api/vpcs"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateVPC",
		resourceName: "VPC",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return model.VpcResponse{}, err
	}

	err = api.pollForVpcReady(ctx, data.ID)
	if err != nil {
		return model.VpcResponse{}, fmt.Errorf("failed to poll for VPC readiness: %w", err)
	}

	name, err := api.readVpcName(ctx, data.ID)
	if err != nil {
		return model.VpcResponse{}, fmt.Errorf("failed to read VPC name: %w", err)
	}
	data.VpcName = name
	return data, nil
}

func (api *API) ReadVPC(ctx context.Context, vpcID int) (*model.VpcResponse, error) {
	var (
		data   model.VpcResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%d", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadVPC",
		resourceName: "VPC",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read VPC: %w", err)
	}

	// Handle resource drift
	if data.ID == 0 {
		return nil, nil
	}

	name, err := api.readVpcName(ctx, data.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read VPC name: %w", err)
	}
	data.VpcName = name
	return &data, nil
}

func (api *API) UpdateVPC(ctx context.Context, vpcID int, params model.VpcRequest) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%d", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v ", path, params))
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateVPC",
		resourceName: "VPC",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

func (api *API) DeleteVPC(ctx context.Context, vpcID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%d", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	return api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteVPC",
		resourceName: "VPC",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

func (api *API) pollForVpcReady(ctx context.Context, vpcID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%d/vpc-peering/info", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	// CustomRetryCode set to 400, to poll for VPC readiness (200)
	return api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName:    "PollForVpcReady",
		resourceName:    "VPC",
		attempt:         1,
		sleep:           5 * time.Second,
		data:            nil,
		failed:          &failed,
		customRetryCode: 400,
	})
}

// readVpcName - retrieves external cloud provider VPC name
func (api *API) readVpcName(ctx context.Context, vpcID int) (string, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%d/vpc-peering/info", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadVpcName",
		resourceName: "VPC",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return "", err
	}

	name, ok := data["name"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid 'name' in VPC peering info: %v", data)
	}
	return name, nil
}
