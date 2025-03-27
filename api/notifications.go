package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) CreateNotification(ctx context.Context, instanceID int, params map[string]any) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s ", path), params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 201:
		tflog.Debug(ctx, "response data", data)
		if v, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		} else {
			return nil, fmt.Errorf("invalid identifier=%v ", data["id"])
		}
		return data, err
	default:
		return nil, fmt.Errorf("failed to create notification, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) ReadNotification(ctx context.Context, instanceID int, recipientID string) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, "response data", data)
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read notification, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) ListNotifications(ctx context.Context, instanceID int) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, fmt.Sprintf("response data=%v", data))
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read notifications, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) UpdateNotification(ctx context.Context, instanceID int, recipientID string,
	params map[string]any) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s ", path))
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	default:
		return fmt.Errorf("failed to update notification, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) DeleteNotification(ctx context.Context, instanceID int, recipientID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to delete notification, status=%d message=%s ",
			response.StatusCode, failed)
	}
}
