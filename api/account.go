package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/network"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) ListInstances(ctx context.Context) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = "api/instances"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, fmt.Sprintf("data: %v", data))
		return data, nil
	case 410: // TODO: Remove should only be needed for a single instance or VPC.
		tflog.Warn(ctx, "status=410 message=\"the instance has been deleted\" ")
		return nil, nil
	default:
		return nil, fmt.Errorf("failed to list instaces, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) ListVpcs(ctx context.Context) ([]model.VpcResponse, error) {
	var (
		data   []model.VpcResponse
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
		return []model.VpcResponse{}, fmt.Errorf("failed to read VPC: %w", err)
	}

	for i := range data {
		name, err := api.readVpcName(ctx, data[i].ID)
		if err != nil {
			return []model.VpcResponse{}, fmt.Errorf("failed to read VPC name: %w", err)
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

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s ", path))
	response, err := api.sling.New().Post(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to rotate api key, status=%d failed=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) RotateApiKey(ctx context.Context, instanceID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/account/rotate-apikey", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s ", path))
	response, err := api.sling.New().Post(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to rotate api key, status=%d failed=%s ",
			response.StatusCode, failed)
	}
}
