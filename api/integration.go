package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreateIntegration enables integration communication, either for logs or metrics.
func (api *API) CreateIntegration(ctx context.Context, instanceID int, intType string,
	intName string, params map[string]any) (map[string]any, error) {

	var (
		data         map[string]any
		failed       map[string]any
		path         = fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, intName)
		sesnitiveCtx = tflog.MaskFieldValuesWithFieldKeys(ctx, "secret_access_key", "private_key",
			"application_secret", "api_key", "token")
	)

	tflog.Debug(sesnitiveCtx, fmt.Sprintf("method=POST path=%s ", path), params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 201:
		tflog.Debug(sesnitiveCtx, "response data", data)
		if v, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		} else {
			return nil, fmt.Errorf("invalid identifier=%v ", data["id"])
		}
		return data, err
	default:
		return nil, fmt.Errorf("failed to create integration, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// ReadIntegration retrieves a specific logs or metrics integration
func (api *API) ReadIntegration(ctx context.Context, instanceID int, intType, intID string) (
	map[string]any, error) {

	var (
		data         map[string]any
		failed       map[string]any
		path         = fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, intID)
		sesnitiveCtx = tflog.MaskFieldValuesWithFieldKeys(ctx, "access_key_id", "application_secret",
			"api_key", "credentials", "license_key", "private_key", "private_key_id", "secret_access_key",
			"token")
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(sesnitiveCtx, "response data", data)
		// Convert API response body, config part, into single map
		convertedData := make(map[string]any)
		for k, v := range data {
			if k == "id" {
				convertedData[k] = v
			} else if k == "type" {
				convertedData[k] = v
			} else if k == "config" {
				for configK, configV := range data["config"].(map[string]any) {
					convertedData[configK] = configV
				}
			}
		}
		return convertedData, err
	case 404:
		tflog.Warn(ctx, "integration not found")
		return nil, nil
	default:
		return nil, fmt.Errorf("failed to read integration, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// UpdateIntegration updated the integration with new information
func (api *API) UpdateIntegration(ctx context.Context, instanceID int, intType, intID string,
	params map[string]any) error {

	var (
		failed       map[string]any
		path         = fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, intID)
		sesnitiveCtx = tflog.MaskFieldValuesWithFieldKeys(ctx, "access_key_id", "application_secret",
			"api_key", "credentials", "license_key", "private_key", "private_key_id", "secret_access_key",
			"token")
	)

	tflog.Debug(sesnitiveCtx, fmt.Sprintf("method=PUT path=%s ", path), params)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to update integration, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// DeleteIntegration removes log or metric integration.
func (api *API) DeleteIntegration(ctx context.Context, instanceID int, intType, intID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, intID)
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
		return fmt.Errorf("failed to delete integration, status=%d message=%s ",
			response.StatusCode, failed)
	}
}
