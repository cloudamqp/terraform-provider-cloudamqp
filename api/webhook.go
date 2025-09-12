package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreateWebhook - create a webhook for a vhost and a specific queue
func (api *API) CreateWebhook(ctx context.Context, instanceID int64, params model.WebhookCreateRequest, sleep time.Duration) (
	string, error) {

	var (
		data   model.WebhookResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/webhooks", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%v ", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateWebhook",
		resourceName: "Webhook",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return "", err
	}

	id := strconv.FormatInt(data.ID, 10)
	return id, nil
}

// ReadWebhook - retrieves a specific webhook for an instance
func (api *API) ReadWebhook(ctx context.Context, instanceID int64, webhookID string, sleep time.Duration) (
	model.WebhookResponse, error) {

	var (
		data   model.WebhookResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/webhooks/%s", instanceID, webhookID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadWebhook",
		resourceName: "Webhook",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return model.WebhookResponse{}, err
	}

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s data=%v ", path, data))

	return data, nil
}

// UpdateWebhook - updates a specific webhook for an instance
func (api *API) UpdateWebhook(ctx context.Context, instanceID int64, webhookID string,
	params model.WebhookUpdateRequest, sleep time.Duration) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/webhooks/%s", instanceID, webhookID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v ", path, params))
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateWebhook",
		resourceName: "Webhook",
		attempt:      1,
		sleep:        sleep,
		data:         nil,
		failed:       &failed,
	})
}

// DeleteWebhook - removes a specific webhook for an instance
func (api *API) DeleteWebhook(ctx context.Context, instanceID int64, webhookID string, sleep time.Duration) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/webhooks/%s", instanceID, webhookID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	return api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteWebhook",
		resourceName: "Webhook",
		attempt:      1,
		sleep:        sleep,
		data:         nil,
		failed:       &failed,
	})
}
