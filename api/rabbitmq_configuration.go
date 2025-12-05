package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ReadRabbitMqConfiguration - retrieves the RabbitMQ configuration for an instance
func (api *API) ReadRabbitMqConfiguration(ctx context.Context, instanceID, sleep int64) (*model.RabbitMqConfigResponse, error) {

	var (
		data   model.RabbitMqConfigResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/config", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadRabbitMqConfiguration",
		resourceName: "RabbitMQConfiguration",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s data=%v ", path, data))
	return &data, nil
}

// UpdateRabbitMqConfiguration - updates the RabbitMQ configuration for an instance
func (api *API) UpdateRabbitMqConfiguration(ctx context.Context, instanceID int64,
	params model.RabbitMqConfigRequest, sleep int64) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/config", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v ", path, params))
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateRabbitMqConfiguration",
		resourceName: "RabbitMQConfiguration",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         nil,
		failed:       &failed,
	})
}
