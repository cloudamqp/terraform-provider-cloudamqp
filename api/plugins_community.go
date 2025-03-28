package api

import (
	"context"
	"fmt"
	"strings"
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
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d timeout=%d ", path, sleep, timeout),
		params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilPluginChanged(ctx, instanceID, pluginName, true, 1, sleep, timeout)
	default:
		return nil, fmt.Errorf("failed to install community plugin, status=%d message=%s ",
			response.StatusCode, failed)
	}
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

	path := fmt.Sprintf("/api/instances/%d/plugins/community", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	return api.listPluginsCommunityWithRetry(ctx, path, 1, sleep, timeout)
}

// listPluginsCommunityWithRetry: list all community plugins for an instance,
// with retry if the backend is busy.
func (api *API) listPluginsCommunityWithRetry(ctx context.Context, path string,
	attempt, sleep, timeout int) ([]map[string]any, error) {

	var (
		data   []map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("timeout reach of %d seconds while listing community plugins", timeout)
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try agian, "+
				"attempt=%d until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.listPluginsCommunityWithRetry(ctx, path, attempt, sleep, timeout)
		}
	}
	return nil, fmt.Errorf("failed to list communit plugins, status=%d message=%s ",
		response.StatusCode, failed)
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
	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s enabled=%t sleep=%d timeout=%d ", path, enabled,
		sleep, timeout), params)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilPluginChanged(ctx, instanceID, pluginName, enabled, 1, sleep, timeout)
	default:
		return nil, fmt.Errorf("failed to update community plugin, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// UninstallPluginCommunity: uninstall a community plugin from an instance.
func (api *API) UninstallPluginCommunity(ctx context.Context, instanceID int, pluginName string,
	sleep, timeout int) (map[string]any, error) {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/community/%s?async=true", instanceID, pluginName)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilPluginUninstalled(ctx, instanceID, pluginName, 1, sleep, timeout)
	default:
		return nil, fmt.Errorf("failed to disable communit plugin, status=%d message=%s ",
			response.StatusCode, failed)
	}
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
