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
	err := CallWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), "VPC", 1, 10*time.Second, &data, &failed)
	if err != nil {
		return model.VpcResponse{}, err
	}

	name, err := api.readVpcName(ctx, data.ID)
	if err != nil {
		return model.VpcResponse{}, fmt.Errorf("failed to read VPC name: %w", err)
	}
	data.VpcName = name
	return data, nil
}

func (api *API) ReadVPC(ctx context.Context, vpcID int) (model.VpcResponse, error) {
	var (
		data   model.VpcResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%d", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := CallWithRetry(ctx, api.sling.New().Get(path), "VPC", 1, 10*time.Second, &data, &failed)
	if err != nil {
		return model.VpcResponse{}, fmt.Errorf("failed to read VPC: %w", err)
	}

	name, err := api.readVpcName(ctx, data.ID)
	if err != nil {
		return model.VpcResponse{}, fmt.Errorf("failed to read VPC name: %w", err)
	}
	data.VpcName = name
	return data, nil
}

func (api *API) readVpcName(ctx context.Context, vpcID int) (string, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/vpcs/%d/vpc-peering/info", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := CallWithRetry(ctx, api.sling.New().Get(path), "VPC", 1, 10*time.Second, &data, &failed)
	if err != nil {
		return "", err
	}

	return data["name"].(string), nil
}

func (api *API) UpdateVPC(ctx context.Context, vpcID string, params model.VpcRequest) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/vpcs/%s", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v ", path, params))
	return CallWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), "VPC", 1, 10*time.Second, nil, &failed)
}

func (api *API) DeleteVPC(ctx context.Context, vpcID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/vpcs/%s", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	return CallWithRetry(ctx, api.sling.New().Delete(path), "VPC", 1, 10*time.Second, nil, &failed)
}
