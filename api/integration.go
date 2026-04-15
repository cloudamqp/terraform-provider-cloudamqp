package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

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

	tflog.Debug(sesnitiveCtx, fmt.Sprintf("method=POST path=%s", path), params)
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateIntegration",
		resourceName: "Integration",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	tflog.Debug(sesnitiveCtx, "response data", data)
	if v, ok := data["id"]; ok {
		data["id"] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
	} else {
		return nil, fmt.Errorf("invalid identifier=%v", data["id"])
	}
	return data, nil
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

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Path(path), retryRequest{
		functionName: "ReadIntegration",
		resourceName: "Integration",
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

	tflog.Debug(sesnitiveCtx, "response data", data)
	// Convert API response body, config part, into single map
	convertedData := make(map[string]any)
	for k, v := range data {
		if k == "id" {
			convertedData[k] = v
		} else if k == "type" {
			convertedData[k] = v
		} else if k == "metrics_filter" {
			convertedData[k] = v
		} else if k == "config" {
			for configK, configV := range data["config"].(map[string]any) {
				convertedData[configK] = configV
			}
		}
	}
	return convertedData, nil
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

	tflog.Debug(sesnitiveCtx, fmt.Sprintf("method=PUT path=%s", path), params)
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateIntegration",
		resourceName: "Integration",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

// DeleteIntegration removes log or metric integration.
func (api *API) DeleteIntegration(ctx context.Context, instanceID int, intType, intID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/%s/%s", instanceID, intType, intID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s", path))
	return api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteIntegration",
		resourceName: "Integration",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

// UpdateMetricsFilter updates the metrics filter for a prometheus integration
func (api *API) UpdateMetricsFilter(ctx context.Context, instanceID int, intID string, filter []string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/integrations/metrics/%s/metrics_filter", instanceID, intID)
		params = map[string]any{"metrics_filter": filter}
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s", path), params)
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateMetricsFilter",
		resourceName: "MetricsFilter",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}
