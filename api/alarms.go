package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/monitoring"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) CreateAlarm(ctx context.Context, instanceID int64, params model.AlarmRequest) (string, error) {

	var (
		data   model.AlarmResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms", instanceID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=POST path=%s params=%+v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateAlarm",
		resourceName: "Alarm",
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

func (api *API) ReadAlarm(ctx context.Context, instanceID int64, alarmID string) (*model.AlarmResponse, error) {

	var (
		data   *model.AlarmResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, alarmID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=GET path=%s", path))
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

	tflog.Info(ctx, fmt.Sprintf("method=GET path=%s data=%+v", path, data))
	// Handle resource drift
	if data == nil {
		return nil, nil
	}
	return data, nil
}

func (api *API) ListAlarms(ctx context.Context, instanceID int64) ([]model.AlarmResponse, error) {
	var (
		data   []model.AlarmResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms", instanceID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=GET path=%s", path))
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

	tflog.Debug(ctx, fmt.Sprintf("response data=%+v", data))
	return data, nil
}

func (api *API) UpdateAlarm(ctx context.Context, instanceID int64, alarmID string, params model.AlarmRequest) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, alarmID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=PUT path=%s params=%+v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateAlarm",
		resourceName: "Alarm",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}
	return nil
}

func (api *API) DeleteAlarm(ctx context.Context, instanceID int64, alarmID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, alarmID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=DELETE path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteAlarm",
		resourceName: "Alarm",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}
	return nil
}
