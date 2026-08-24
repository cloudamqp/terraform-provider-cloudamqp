package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreateIntegrationLogAgent - create a log integration for a specific instance
func (api *API) CreateIntegrationLogAgent(ctx context.Context, instanceID int64, intType string, params model.LogAgentRequest) (string, error) {
	var (
		data   model.LogAgentResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/logs/%s", instanceID, intType)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%+v ", path, params.Sanitized()))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateIntegrationLogAgent",
		resourceName: "IntegrationLogAgent",
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

// ReadIntegrationLogAgent - retrieves a specific integration log for an instance
func (api *API) ReadIntegrationLogAgent(ctx context.Context, instanceID int64, logID string) (*model.LogAgentResponse, error) {

	var (
		data   model.LogAgentResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/logs/%s", instanceID, logID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadIntegrationLogAgent",
		resourceName: "IntegrationLogAgent",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s data=%+v ", path, data.Sanitized()))

	// Handle resource drift
	if data.ID == 0 {
		return nil, nil
	}
	return &data, nil
}

// UpdateIntegrationLogAgent - updates a specific integration log for an instance
func (api *API) UpdateIntegrationLogAgent(ctx context.Context, instanceID int64, logID string, params model.LogAgentRequest) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/logs/%s", instanceID, logID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%+v ", path, params.Sanitized()))
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateIntegrationLogAgent",
		resourceName: "IntegrationLogAgent",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

// DeleteIntegrationLogAgent - removes a specific integration log for an instance
func (api *API) DeleteIntegrationLogAgent(ctx context.Context, instanceID int64, logID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/logs/%s", instanceID, logID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	return api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteIntegrationLogAgent",
		resourceName: "IntegrationLogAgent",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}
