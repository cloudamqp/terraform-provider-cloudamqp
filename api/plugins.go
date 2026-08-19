package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// EnablePlugin: enable a plugin on an instance.
func (api *API) EnablePlugin(ctx context.Context, instanceID int64, params model.PluginRequest, sleep int64) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins?async=true", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d params=%+v", path, sleep, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "EnablePlugin",
		resourceName: "Plugin",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	return api.pollPluginChanged(ctx, instanceID, params, 1, sleep)
}

// ReadPlugin: reads a specific plugin from an instance.
func (api *API) ReadPlugin(ctx context.Context, instanceID int64, pluginName string, sleep int64) (*model.PluginResponse, error) {
	data, err := api.ListPlugins(ctx, instanceID, sleep)
	if err != nil {
		return nil, err
	}

	for _, plugin := range data {
		if plugin.Name == pluginName {
			tflog.Debug(ctx, fmt.Sprintf("plugin found, %s", pluginName))
			return &plugin, nil
		}
	}

	return nil, nil
}

// ListPlugins: list plugins from an instance.
func (api *API) ListPlugins(ctx context.Context, instanceID, sleep int64) ([]model.PluginResponse, error) {
	var (
		data   []model.PluginResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d ", path, sleep))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
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
func (api *API) UpdatePlugin(ctx context.Context, instanceID int64, params model.PluginRequest, sleep int64) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins?async=true", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d params=%+v", path, sleep, params))
	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdatePlugin",
		resourceName: "Plugin",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	return api.pollPluginChanged(ctx, instanceID, params, 1, sleep)
}

// DisablePlugin: disables a plugin from an instance.
func (api *API) DisablePlugin(ctx context.Context, instanceID int64, params model.PluginRequest, sleep int64) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, params.Name)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d", path, sleep))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DisablePlugin",
		resourceName: "Plugin",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	return api.pollPluginChanged(ctx, instanceID, params, 1, sleep)
}

// DeletePlugin: deletes a plugin from an instance.
func (api *API) DeletePlugin(ctx context.Context, instanceID int64, params model.PluginRequest, sleep int64) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, params.Name)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d", path, sleep))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
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

	return api.pollPluginChanged(ctx, instanceID, params, 1, sleep)
}

// pollPluginChanged: poll plugin status until it has changed to the desired state or timeout.
func (api *API) pollPluginChanged(ctx context.Context, instanceID int64, params model.PluginRequest, attempt, sleep int64) error {

	tflog.Debug(ctx, fmt.Sprintf("waiting until plugin status has changed, instanceID=%d sleep=%d params=%+v", instanceID, sleep, params))
	for {
		if ctx.Err() != nil {
			return fmt.Errorf("timeout reached while waiting until plugin status been changed: %w", ctx.Err())
		}

		tflog.Debug(ctx, fmt.Sprintf("Checking plugin status, attempt=%d", attempt))
		response, err := api.ReadPlugin(ctx, instanceID, params.Name, sleep)
		if err != nil {
			return err
		}

		if response == nil {
			attempt++
			select {
			case <-ctx.Done():
				return fmt.Errorf("timeout reached while waiting until plugin status been changed: %w", ctx.Err())
			case <-time.After(time.Duration(sleep) * time.Second):
				continue
			}
		}

		// TODO: How to handle required plugins? If a plugin is required, it cannot be disabled.
		if response.Required != "" && response.Required != "false" {
			return nil
		}
		if status := response.Enabled; status == params.Enabled {
			return nil
		}

		tflog.Debug(ctx, fmt.Sprintf("Plugin not in desired state yet, attempt=%d", attempt))
		attempt++
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout reached while waiting until plugin status been changed: %w", ctx.Err())
		case <-time.After(time.Duration(sleep) * time.Second):
			continue
		}
	}
}
