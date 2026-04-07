package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/monitoring"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) CreateNotification(ctx context.Context, instanceID int64, params *model.RecipientRequest) (*model.RecipientResponse, error) {
	var (
		data   *model.RecipientResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s, params=%+v", path, params.Sanitized()))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateNotification",
		resourceName: "Notification",
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

func (api *API) ReadNotification(ctx context.Context, instanceID int64, recipientID string) (*model.RecipientResponse, error) {
	var (
		data   *model.RecipientResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Path(path), retryRequest{
		functionName: "ReadNotification",
		resourceName: "Notification",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("ReadNotification response: %+v", data.Sanitized()))
	// Handle resource drift
	if data == nil {
		return nil, nil
	}

	return data, nil
}

func (api *API) ReadNotificationByName(ctx context.Context, instanceID int64, name string) (*model.RecipientResponse, error) {
	notifications, err := api.ListNotifications(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, fmt.Sprintf("Looking for notification with name '%s' among %d notifications", name, len(notifications)))
	for _, notification := range notifications {
		tflog.Info(ctx, fmt.Sprintf("Checking notification with name '%s': %+v", notification.Name, notification.Sanitized()))
		if notification.Name == name {
			tflog.Info(ctx, fmt.Sprintf("Found notification with name '%s': %+v", name, notification.Sanitized()))
			return &notification, nil
		}
	}

	return nil, nil
}

func (api *API) ListNotifications(ctx context.Context, instanceID int64) ([]model.RecipientResponse, error) {
	var (
		data   []model.RecipientResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Path(path), retryRequest{
		functionName: "ListNotifications",
		resourceName: "Notification",
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

func (api *API) UpdateNotification(ctx context.Context, instanceID int64, recipientID string, params *model.RecipientRequest) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s, params=%+v", path, params.Sanitized()))
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateNotification",
		resourceName: "Notification",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

func (api *API) DeleteNotification(ctx context.Context, instanceID int64, recipientID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s", path))
	return api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteNotification",
		resourceName: "Notification",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}
