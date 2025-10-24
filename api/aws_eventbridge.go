package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreateAwsEventBridge - create an AWS EventBridge for a vhost and a specific queue
func (api *API) CreateAwsEventBridge(ctx context.Context, instanceID int64, params model.AwsEventBridgeRequest) (string, error) {

	var (
		data   model.AwsEventBridgeResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/eventbridges", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%v ", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateAwsEventBridge",
		resourceName: "AwsEventBridge",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return "", err
	}

	id := strconv.FormatInt(data.Id, 10)
	return id, nil
}

// ReadAwsEventBridge - read an AWS EventBridge by its and instance identifiers
func (api *API) ReadAwsEventBridge(ctx context.Context, instanceID int64, eventbridgeID string) (*model.AwsEventBridgeResponse, error) {

	var (
		data   model.AwsEventBridgeResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/eventbridges/%s", instanceID, eventbridgeID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadAwsEventBridge",
		resourceName: "AwsEventBridge",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read AWS EventBridge: %w", err)
	}

	// Handle resource drift
	if data.Id == 0 {
		return nil, nil
	}

	return &data, nil
}

// DeleteAwsEventBridge - delete an AWS EventBridge by its and instance identifiers
func (api *API) DeleteAwsEventBridge(ctx context.Context, instanceID int64, eventbridgeID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/eventbridges/%s", instanceID, eventbridgeID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	return api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteAwsEventBridge",
		resourceName: "AwsEventBridge",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}
