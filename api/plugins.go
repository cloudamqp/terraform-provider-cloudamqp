package api

import (
	"context"
	"fmt"
	"strings"
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
		return nil, fmt.Errorf("failed to enable/disable plugin, status=%d message=%s ",
			response.StatusCode, failed)
	}
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

	path := fmt.Sprintf("/api/instances/%d/plugins", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	return api.listPluginsWithRetry(ctx, path, 1, sleep, timeout)
}

// listPluginsWithRetry: list plugins from an instance, with retry if backend is busy.
func (api *API) listPluginsWithRetry(ctx context.Context, path string, attempt, sleep,
	timeout int) ([]map[string]any, error) {

	var (
		data   []map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("timeout reached after %d seconds, while reading plugins", timeout)
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
				"attempt=%d until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.listPluginsWithRetry(ctx, path, attempt, sleep, timeout)
		}
	case 404:
		// Instance not found - likely manually deleted
		tflog.Debug(ctx, fmt.Sprintf("instance not found (404), likely manually deleted: %s", path))
		return nil, fmt.Errorf("instance not found, status=404 message=%s", failed)
	case 423:
		tflog.Debug(ctx, fmt.Sprintf("resource is locked, will try again, attempt=%d ", attempt))
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
		return api.listPluginsWithRetry(ctx, path, attempt, sleep, timeout)
	case 503:
		tflog.Debug(ctx, fmt.Sprintf("service unavailable, will try again, attempt=%d ", attempt))
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
		return api.listPluginsWithRetry(ctx, path, attempt, sleep, timeout)
	}
	return nil, fmt.Errorf("failed to list plugins, status=%d message=%s ",
		response.StatusCode, failed)
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
	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d timeout=%d ", path, sleep, timeout),
		params)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilPluginChanged(ctx, instanceID, pluginName, enabled, 1, sleep, timeout)
	default:
		return nil, fmt.Errorf("failed to update plugin, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// DisablePlugin: disables a plugin from an instance.
func (api *API) DisablePlugin(ctx context.Context, instanceID int, pluginName string, sleep,
	timeout int) (map[string]any, error) {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, pluginName)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilPluginChanged(ctx, instanceID, pluginName, false, 1, sleep, timeout)
	default:
		return nil, fmt.Errorf("failed to disable plugin, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// DeletePlugin: deletes a plugin from an instance.
func (api *API) DeletePlugin(ctx context.Context, instanceID int, pluginName string,
	sleep, timeout int) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, pluginName)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		_, err = api.waitUntilPluginChanged(ctx, instanceID, pluginName, false, 1, sleep, timeout)
		return err
	case 404:
		// Instance not found - likely manually deleted
		tflog.Debug(ctx, fmt.Sprintf("instance not found (404) during plugin deletion: %s", path))
		return fmt.Errorf("instance not found, status=404 message=%s", failed)
	default:
		return fmt.Errorf("failed to delete plugin, status=%d message=%s ", response.StatusCode, failed)
	}
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
