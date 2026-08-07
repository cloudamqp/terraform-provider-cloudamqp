package api

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/cloudamqp/terraform-provider-cloudamqp/api/models/job"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreatePluginBatch: enable multiple plugins via the batch endpoint.
func (api *API) CreatePluginBatch(ctx context.Context, instanceID int64,
	params instance.PluginBatchRequest) (*job.JobCreationResponse, error) {

	var (
		data   *job.JobCreationResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/batch", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%+v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreatePluginBatch",
		resourceName: "PluginBatch",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("no data received from API")
	}

	return data, nil
}

// UpdatePluginBatch: enable and/or disable multiple plugins via the batch endpoint.
func (api *API) UpdatePluginBatch(ctx context.Context, instanceID int64,
	params instance.PluginBatchRequest) (*job.JobCreationResponse, error) {

	var (
		data   *job.JobCreationResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/batch", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%+v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdatePluginBatch",
		resourceName: "PluginBatch",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("no data received from API")
	}

	return data, nil
}

// DeletePluginBatch: disable multiple plugins via the batch endpoint.
func (api *API) DeletePluginBatch(ctx context.Context, instanceID int64,
	params instance.PluginBatchRequest) (*job.JobCreationResponse, error) {

	var (
		data   *job.JobCreationResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/batch", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s params=%+v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path).BodyJSON(params), retryRequest{
		functionName: "DeletePluginBatch",
		resourceName: "PluginBatch",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("no data received from API")
	}

	return data, nil
}
