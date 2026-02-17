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

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s", path), params)
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateAlarm",
		resourceName: "Alarm",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	if id, ok := data["id"]; ok {
		data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
	} else {
		return nil, fmt.Errorf("invalid identifier=%v", data["id"])
	}
	return data, nil
}

func (api *API) ReadAlarm(ctx context.Context, instanceID int, alarmID string) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, alarmID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadAlarm",
		resourceName: "Alarm",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	// Handle resource drift
	if len(data) == 0 {
		return nil, nil
	}

	return data, nil
}

func (api *API) ListAlarms(ctx context.Context, instanceID int) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ListAlarms",
		resourceName: "Alarm",
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

func (api *API) UpdateAlarm(ctx context.Context, instanceID int, params map[string]any) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, params["id"].(string))
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s", path), params)
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateAlarm",
		resourceName: "Alarm",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

func (api *API) DeleteAlarm(ctx context.Context, instanceID int, alarmID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, alarmID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s", path))
	return api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteAlarm",
		resourceName: "Alarm",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}
