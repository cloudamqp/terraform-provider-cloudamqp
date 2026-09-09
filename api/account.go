package api

import (
	"context"
	"fmt"
	"time"

	instanceModel "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	networkModel "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/network"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) ListInstances(ctx context.Context) ([]instanceModel.InstanceResponse, error) {
	var (
		data   []instanceModel.InstanceResponse
		failed map[string]any
		path   = "api/instances"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Path(path), retryRequest{
		functionName: "ListInstances",
		resourceName: "Instance",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (api *API) ListVpcs(ctx context.Context) ([]networkModel.VpcResponse, error) {
	var (
		data   []networkModel.VpcResponse
		failed map[string]any
		path   = "/api/vpcs"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ListVpcs",
		resourceName: "VPC",
		attempt:      1,
		sleep:        10 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return []networkModel.VpcResponse{}, fmt.Errorf("failed to read VPC: %w", err)
	}

	for i := range data {
		name, err := api.readVpcName(ctx, data[i].ID)
		if err != nil {
			return []networkModel.VpcResponse{}, fmt.Errorf("failed to read VPC name: %w", err)
		}
		data[i].VpcName = name
	}
	return data, nil
}

func (api *API) RotatePassword(ctx context.Context, instanceID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/account/rotate-password", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s", path))
	return api.callWithRetry(ctx, api.sling.New().Post(path), retryRequest{
		functionName: "RotatePassword",
		resourceName: "Password",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

func (api *API) RotateApiKey(ctx context.Context, instanceID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/account/rotate-apikey", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s", path))
	return api.callWithRetry(ctx, api.sling.New().Post(path), retryRequest{
		functionName: "RotateApiKey",
		resourceName: "ApiKey",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}
