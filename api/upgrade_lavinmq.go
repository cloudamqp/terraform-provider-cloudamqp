package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ReadVersions - Read versions LavinMQ can upgrade to
func (api *API) ReadLavinMQVersions(ctx context.Context, instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/new-lavinmq-versions", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Path(path), retryRequest{
		functionName: "ReadLavinMQVersions",
		resourceName: "LavinMQ versions",
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

// UpgradeLavinMQ - Upgrade to latest possible version or a specific available version
func (api *API) UpgradeLavinMQ(ctx context.Context, instanceID int, new_version string) (
	string, error) {

	if new_version == "" {
		return api.UpgradeToLatestLavinMQVersion(ctx, instanceID)
	} else {
		return api.UpgradeToSpecificLavinMQVersion(ctx, instanceID, new_version)
	}
}

func (api *API) UpgradeToSpecificLavinMQVersion(ctx context.Context, instanceID int, version string) (
	string, error) {

	var (
		data       map[string]any
		failed     map[string]any
		statusCode int
		path       = fmt.Sprintf("api/instances/%d/actions/upgrade-lavinmq", instanceID)
		params     = make(map[string]any)
	)

	params["version"] = version
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s version=%s upgrade to specific version",
		path, version), params)
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "UpgradeToSpecificLavinMQVersion",
		resourceName: "LavinMQ",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
		statusCode:   &statusCode,
	})
	if err != nil {
		return "", err
	}

	tflog.Debug(ctx, "response data", data)

	// Handle different success codes
	if statusCode == 200 {
		return "Already at highest possible version", nil
	}

	return api.waitUntilLavinMQUpgraded(ctx, instanceID)
}

func (api *API) UpgradeToLatestLavinMQVersion(ctx context.Context, instanceID int) (string, error) {
	var (
		data       map[string]any
		failed     map[string]any
		statusCode int
		path       = fmt.Sprintf("api/instances/%d/actions/upgrade-lavinmq", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s upgrade to latest version", path))
	err := api.callWithRetry(ctx, api.sling.New().Post(path), retryRequest{
		functionName: "UpgradeToLatestLavinMQVersion",
		resourceName: "LavinMQ",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
		statusCode:   &statusCode,
	})
	if err != nil {
		return "", err
	}

	// Handle different success codes
	if statusCode == 200 {
		return "Already at highest possible version", nil
	}

	return api.waitUntilLavinMQUpgraded(ctx, instanceID)
}

func (api *API) waitUntilLavinMQUpgraded(ctx context.Context, instanceID int) (string, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/nodes", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("waiting until LavinMQ been upgraded, method=GET path=%s ", path))
	for {
		_, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return "", err
		}

		tflog.Debug(ctx, fmt.Sprintf("data: %v", data))
		ready := true
		for _, node := range data {
			ready = ready && node["configured"].(bool)
		}
		if ready {
			return "", nil
		}
		time.Sleep(10 * time.Second)
	}
}
