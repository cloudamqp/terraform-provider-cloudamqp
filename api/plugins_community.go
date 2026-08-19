package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// InstallPluginCommunity: install a community plugin on an instance.
func (api *API) InstallPluginCommunity(ctx context.Context, instanceID int64, params model.PluginRequest, sleep int64) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/community?async=true", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d params=%+v", path, sleep, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "InstallPluginCommunity",
		resourceName: "CommunityPlugin",
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

// ReadPluginCommunity: reads a specific community plugin from an instance.
func (api *API) ReadPluginCommunity(ctx context.Context, instanceID int64, name string, sleep int64) (*model.PluginResponse, error) {

	data, err := api.ListPluginsCommunity(ctx, instanceID, sleep)
	if err != nil {
		return nil, err
	}

	for _, plugin := range data {
		if plugin.Name == name {
			tflog.Debug(ctx, fmt.Sprintf("community plugin found, %s ", name))
			return &plugin, nil
		}
	}

	return nil, nil
}

// ListPluginsCommunity: list all community plugins for an instance.
func (api *API) ListPluginsCommunity(ctx context.Context, instanceID, sleep int64) ([]model.PluginResponse, error) {
	var (
		data   []model.PluginResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/community", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d", path, sleep))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
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
func (api *API) UpdatePluginCommunity(ctx context.Context, instanceID int64, params model.PluginRequest, sleep int64) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/community?async=true", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d params=%+v", path, sleep, params))
	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdatePluginCommunity",
		resourceName: "CommunityPlugin",
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

// UninstallPluginCommunity: uninstall a community plugin from an instance.
func (api *API) UninstallPluginCommunity(ctx context.Context, instanceID int64, pluginName string, sleep int64) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/community/%s?async=true", instanceID, pluginName)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d", path, sleep))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "UninstallPluginCommunity",
		resourceName: "CommunityPlugin",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	return api.waitUntilPluginUninstalled(ctx, instanceID, pluginName, 1, sleep)
}

// waitUntilPluginUninstalled: wait until a community plugin been uninstalled.
func (api *API) waitUntilPluginUninstalled(ctx context.Context, instanceID int64, pluginName string, attempt, sleep int64) error {

	tflog.Debug(ctx, fmt.Sprintf("waiting for community plugin to be uninstalled, instanceID=%d plugin=%s sleep=%d", instanceID, pluginName, sleep))
	for {
		if ctx.Err() != nil {
			return fmt.Errorf("timeout reached while waiting on community plugin being uninstalled: %w", ctx.Err())
		}

		tflog.Debug(ctx, fmt.Sprintf("Checking community plugin uninstall status, attempt=%d", attempt))
		response, err := api.ReadPluginCommunity(ctx, instanceID, pluginName, sleep)
		if err != nil {
			return err
		}
		if response == nil {
			return nil
		}

		tflog.Debug(ctx, fmt.Sprintf("Community plugin still installed, attempt=%d", attempt))
		attempt++
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout reached while waiting on community plugin being uninstalled: %w", ctx.Err())
		case <-time.After(time.Duration(sleep) * time.Second):
			continue
		}
	}
}
