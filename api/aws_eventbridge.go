package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) CreateAwsEventBridge(ctx context.Context, instanceID int, params map[string]any) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/eventbridges", instanceID)
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
		return data, nil
	default:
		return nil, fmt.Errorf("failed to create AWS EventBridge, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) ReadAwsEventBridge(ctx context.Context, instanceID int, eventbridgeID string) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/eventbridges/%s", instanceID, eventbridgeID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, "response data", data)
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read AWS EventBridge, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) ReadAwsEventBridges(ctx context.Context, instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/eventbridges", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, "response data", data)
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read AWS EventBridges, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) DeleteAwsEventBridge(ctx context.Context, instanceID int, eventbridgeID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/eventbridges/%s", instanceID, eventbridgeID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	case 404:
		// AWS EventBridge not found in the backend. Silent let the resource be deleted.
		return nil
	default:
		return fmt.Errorf("failed to delete AWS EventBridge, status=%d message=%s ",
			response.StatusCode, failed)
	}
}
