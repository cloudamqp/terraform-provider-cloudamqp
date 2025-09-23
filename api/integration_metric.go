package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreateIntegrationMetric - create a metric integration for a specific instance
func (api *API) CreateIntegrationMetric(ctx context.Context, instanceID int64, intName string, params model.MetricRequest) (string, error) {
	var (
		data   model.MetricResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/metrics/%s", instanceID, intName)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%+v ", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateIntegrationMetric",
		resourceName: "IntegrationMetric",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return "", err
	}

	id := strconv.FormatInt(data.ID, 10)
	return id, nil
}

// ReadIntegrationMetric - retrieves a specific integration metric for an instance
func (api *API) ReadIntegrationMetric(ctx context.Context, instanceID int64, metricID string) (*model.MetricResponse, error) {

	var (
		data   model.MetricResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/metrics/%s", instanceID, metricID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadIntegrationMetric",
		resourceName: "IntegrationMetric",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s data=%+v ", path, data))

	return &data, nil
}

// UpdateIntegrationMetric - updates a specific integration metric for an instance
func (api *API) UpdateIntegrationMetric(ctx context.Context, instanceID int64, metricID string, params model.MetricRequest) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/metrics/%s", instanceID, metricID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v ", path, params))
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateIntegrationMetric",
		resourceName: "IntegrationMetric",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

// DeleteIntegrationMetric - removes a specific integration metric for an instance
func (api *API) DeleteIntegrationMetric(ctx context.Context, instanceID int64, metricID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/metrics/%s", instanceID, metricID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	return api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteIntegrationMetric",
		resourceName: "IntegrationMetric",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}
