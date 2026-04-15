package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// InstallPluginCommunity: install a community plugin on an instance.
func (api *API) InstallPluginCommunity(ctx context.Context, instanceID int, pluginName string,
	sleep, timeout int) (map[string]any, error) {

	var (
		failed map[string]any
		params = make(map[string]any)
		path   = fmt.Sprintf("/api/instances/%d/plugins/community?async=true", instanceID)
	)

	params["plugin_name"] = pluginName
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d timeout=%d", path, sleep, timeout), params)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "InstallPluginCommunity",
		resourceName: "CommunityPlugin",
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

// ReadPluginCommunity: reads a specific community plugin from an instance.
func (api *API) ReadPluginCommunity(ctx context.Context, instanceID int, pluginName string, sleep,
	timeout int) (map[string]any, error) {

	data, err := api.ListPluginsCommunity(ctx, instanceID, sleep, timeout)
	if err != nil {
		return nil, err
	}

	for _, plugin := range data {
		if plugin["name"] == pluginName {
			tflog.Debug(ctx, fmt.Sprintf("plugin found, %s ", pluginName))
			return plugin, nil
		}
	}

	return nil, nil
}

// ListPluginsCommunity: list all community plugins for an instance.
func (api *API) ListPluginsCommunity(ctx context.Context, instanceID, sleep, timeout int) (
	[]map[string]any, error) {

	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/community", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d", path, sleep, timeout))

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Get(path), retryRequest{
		functionName: "ListPluginsCommunity",
		resourceName: "CommunityPlugin",
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

// UpdatePluginCommunity: updates a community plugin from an instance.
func (api *API) UpdatePluginCommunity(ctx context.Context, instanceID int, pluginName string,
	enabled bool, sleep, timeout int) (map[string]any, error) {

	var (
		failed map[string]any
		params = make(map[string]any)
		path   = fmt.Sprintf("/api/instances/%d/plugins/community?async=true", instanceID)
	)

	params["plugin_name"] = pluginName
	params["enabled"] = enabled
	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s enabled=%t sleep=%d timeout=%d", path, enabled,
		sleep, timeout), params)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdatePluginCommunity",
		resourceName: "CommunityPlugin",
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

// UninstallPluginCommunity: uninstall a community plugin from an instance.
func (api *API) UninstallPluginCommunity(ctx context.Context, instanceID int, pluginName string,
	sleep, timeout int) (map[string]any, error) {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/community/%s?async=true", instanceID, pluginName)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d", path, sleep, timeout))

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Delete(path), retryRequest{
		functionName: "UninstallPluginCommunity",
		resourceName: "CommunityPlugin",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return api.waitUntilPluginUninstalled(ctx, instanceID, pluginName, 1, sleep, timeout)
}

// waitUntilPluginUninstalled: wait until a community plugin been uninstalled.
func (api *API) waitUntilPluginUninstalled(ctx context.Context, instanceID int, pluginName string,
	attempt, sleep, timeout int) (map[string]any, error) {

	tflog.Debug(ctx, fmt.Sprintf("waiting for community plugin, %s, to be uninstalled", pluginName))
	for {
		if attempt*sleep > timeout {
			return nil, fmt.Errorf("timeout reached of %d seconds, while waiting on communit plugin "+
				"being uninstalled", timeout)
		}

		response, err := api.ReadPlugin(ctx, instanceID, pluginName, sleep, timeout)
		if err != nil {
			return nil, err
		}
		if len(response) == 0 {
			return response, nil
		}
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}
