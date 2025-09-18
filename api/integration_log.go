package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreateIntegrationLog - create a log integration for a specific instance
func (api *API) CreateIntegrationLog(ctx context.Context, instanceID int64, intName string, params model.LogRequest) (string, error) {
	var (
		data   model.LogResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/logs/%s", instanceID, intName)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%+v ", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateIntegrationLog",
		resourceName: "IntegrationLog",
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

// ReadIntegrationLog - retrieves a specific integration log for an instance
func (api *API) ReadIntegrationLog(ctx context.Context, instanceID int64, logID string) (*model.LogResponse, error) {

	var (
		data   model.LogResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/logs/%s", instanceID, logID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadIntegrationLog",
		resourceName: "IntegrationLog",
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

// UpdateIntegrationLog - updates a specific integration log for an instance
func (api *API) UpdateIntegrationLog(ctx context.Context, instanceID int64, logID string, params model.LogRequest) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/logs/%s", instanceID, logID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v ", path, params))
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateIntegrationLog",
		resourceName: "IntegrationLog",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

// DeleteIntegrationLog - removes a specific integration log for an instance
func (api *API) DeleteIntegrationLog(ctx context.Context, instanceID int64, logID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/logs/%s", instanceID, logID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	return api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteIntegrationLog",
		resourceName: "IntegrationLog",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}
