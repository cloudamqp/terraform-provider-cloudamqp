package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// EnablePlugin: enable a plugin on an instance.
func (api *API) EnablePlugin(ctx context.Context, instanceID int, pluginName string, sleep,
	timeout int) (map[string]any, error) {

	var (
		failed map[string]any
		params = make(map[string]any)
		path   = fmt.Sprintf("/api/instances/%d/plugins?async=true", instanceID)
	)

	params["plugin_name"] = pluginName
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d timeout=%d", path, sleep, timeout), params)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "EnablePlugin",
		resourceName: "Plugin",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return api.waitUntilPluginChanged(ctx, instanceID, pluginName, true, 1, sleep, timeout)
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
func (api *API) UpdatePlugin(ctx context.Context, instanceID int, pluginName string, enabled bool,
	sleep, timeout int) (map[string]any, error) {

	var (
		failed map[string]any
		params = make(map[string]any)
		path   = fmt.Sprintf("/api/instances/%d/plugins?async=true", instanceID)
	)

	params["plugin_name"] = pluginName
	params["enabled"] = enabled
	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d timeout=%d", path, sleep, timeout), params)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdatePlugin",
		resourceName: "Plugin",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return api.waitUntilPluginChanged(ctx, instanceID, pluginName, enabled, 1, sleep, timeout)
}

// DisablePlugin: disables a plugin from an instance.
func (api *API) DisablePlugin(ctx context.Context, instanceID int, pluginName string, sleep,
	timeout int) (map[string]any, error) {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, pluginName)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d", path, sleep, timeout))

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Delete(path), retryRequest{
		functionName: "DisablePlugin",
		resourceName: "Plugin",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return api.waitUntilPluginChanged(ctx, instanceID, pluginName, false, 1, sleep, timeout)
}

// DeletePlugin: deletes a plugin from an instance.
func (api *API) DeletePlugin(ctx context.Context, instanceID int, pluginName string,
	sleep, timeout int) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, pluginName)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d", path, sleep, timeout))

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeletePlugin",
		resourceName: "Plugin",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	_, err = api.waitUntilPluginChanged(ctx, instanceID, pluginName, false, 1, sleep, timeout)
	return err
}

// waitUntilPluginChanged: wait until plugin changed.
func (api *API) waitUntilPluginChanged(ctx context.Context, instanceID int, pluginName string,
	enabled bool, attempt, sleep, timeout int) (map[string]any, error) {

	tflog.Debug(ctx, "waiting until plugin status been changed")
	for {
		if attempt*sleep > timeout {
			return nil, fmt.Errorf("timeout reached after %d seconds, while waiting until plugin status "+
				"been changed", timeout)
		}

		response, err := api.ReadPlugin(ctx, instanceID, pluginName, sleep, timeout)
		if err != nil {
			return nil, err
		}
		if response["required"] != nil && response["required"] != false {
			return response, nil
		}
		if response["enabled"] == enabled {
			return response, nil
		}
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}
