package api

import (
	"context"
	"fmt"
	"time"

	job "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/job"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// EnablePlugin: enable a plugin on an instance.
func (api *API) EnablePlugin(ctx context.Context, instanceID int, pluginName string) (*job.JobCreationResponse, error) {

	var (
		data   *job.JobCreationResponse
		failed map[string]any
		params = make(map[string]any)
		path   = fmt.Sprintf("/api/instances/%d/plugins?async=true", instanceID)
	)

	params["plugin_name"] = pluginName
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s", path), params)

	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "EnablePlugin",
		resourceName: "Plugin",
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

// ReadPlugin: reads a specific plugin from an instance.
func (api *API) ReadPlugin(ctx context.Context, instanceID int, pluginName string, sleep,
	timeout int) (map[string]any, error) {

	data, err := api.ListPlugins(ctx, instanceID, sleep, timeout)
	if err != nil {
		return nil, err
	}

	for _, plugin := range data {
		if plugin["name"] == pluginName {
			tflog.Debug(ctx, fmt.Sprintf("plugin found, %s", pluginName))
			return plugin, nil
		}
	}

	return nil, nil
}

// ListPlugins: list plugins from an instance.
func (api *API) ListPlugins(ctx context.Context, instanceID, sleep, timeout int) (
	[]map[string]any, error) {

	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d", path, sleep, timeout))

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Get(path), retryRequest{
		functionName: "ListPlugins",
		resourceName: "Plugin",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UpdatePlugin: updates a plugin from an instance.
func (api *API) UpdatePlugin(ctx context.Context, instanceID int, pluginName string, enabled bool) (*job.JobCreationResponse, error) {

	var (
		data   *job.JobCreationResponse
		failed map[string]any
		params = make(map[string]any)
		path   = fmt.Sprintf("/api/instances/%d/plugins?async=true", instanceID)
	)

	params["plugin_name"] = pluginName
	params["enabled"] = enabled
	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s", path), params)

	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdatePlugin",
		resourceName: "Plugin",
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

// DisablePlugin: disables a plugin from an instance.
func (api *API) DisablePlugin(ctx context.Context, instanceID int, pluginName string) (*job.JobCreationResponse, error) {

	var (
		data   *job.JobCreationResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, pluginName)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s", path))

	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DisablePlugin",
		resourceName: "Plugin",
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

// DeletePlugin: deletes a plugin from an instance.
func (api *API) DeletePlugin(ctx context.Context, instanceID int, pluginName string) (*job.JobCreationResponse, error) {

	var (
		data   *job.JobCreationResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, pluginName)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s", path))

	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeletePlugin",
		resourceName: "Plugin",
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
