package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) CreateAlarm(ctx context.Context, instanceID int, params map[string]any) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s ", path), params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 201:
		tflog.Debug(ctx, "response data", data)
		if id, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
		} else {
			return nil, fmt.Errorf("invalid identifier=%v ", data["id"])
		}
		return data, err
	default:
		return nil, fmt.Errorf("failed to create alarm, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) ReadAlarm(ctx context.Context, instanceID int, alarmID string) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, alarmID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, "response data", data)
		return data, err
	default:
		return nil,
			fmt.Errorf("failed to read alarm, status=%d message=%s ", response.StatusCode, failed)
	}
}

func (api *API) ListAlarms(ctx context.Context, instanceID int) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, fmt.Sprintf("response data=%v ", data))
		return data, err
	default:
		return nil,
			fmt.Errorf("failed to list alarms, status=%d message=%s ", response.StatusCode, failed)
	}
}

func (api *API) UpdateAlarm(ctx context.Context, instanceID int, params map[string]any) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, params["id"].(string))
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s ", path), params)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 201:
		return nil
	default:
		return fmt.Errorf("failed to update alarm, status=%d message=%s ", response.StatusCode, failed)
	}
}

func (api *API) DeleteAlarm(ctx context.Context, instanceID int, alarmID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, alarmID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilAlarmDeletion(ctx, instanceID, alarmID)
	default:
		return fmt.Errorf("failed to delete alarm, status=%d message=%s ", response.StatusCode, failed)
	}
}

// TODO: Should not be needed. Removed from backend and should be instant.
func (api *API) waitUntilAlarmDeletion(ctx context.Context, instanceID int, alarmID string) error {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, alarmID)
	)

	tflog.Debug(ctx, "waiting on deletion")
	for {
		response, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return fmt.Errorf("failed to delete alarm, status=%d message=%s ",
				response.StatusCode, failed)
		}

		switch response.StatusCode {
		case 404:
			tflog.Debug(ctx, fmt.Sprintf("alarm with identifier %s deleted", alarmID))
			return nil
		}

		time.Sleep(10 * time.Second)
	}
}
