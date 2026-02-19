package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ReadVersions - Read versions RabbitMQ and Erlang can upgrade to
func (api *API) ReadVersions(ctx context.Context, instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/new-rabbitmq-erlang-versions", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Path(path), retryRequest{
		functionName: "ReadVersions",
		resourceName: "RabbitMQ versions",
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

// UpgradeRabbitMQ - Upgrade to latest possible version or a specific available version
func (api *API) UpgradeRabbitMQ(ctx context.Context, instanceID int, current_version,
	new_version string) (string, error) {

	tflog.Debug(ctx, fmt.Sprintf("instanceID=%d current_version=%s new_version=%s "+
		"upgrade RabbitMQ version", instanceID, current_version, new_version))
	// Keep old behaviour
	if current_version == "" && new_version == "" {
		return api.UpgradeToLatestVersion(ctx, instanceID)
	} else if current_version != "" {
		return api.UpgradeToLatestVersion(ctx, instanceID)
	} else {
		return api.UpgradeToSpecificVersion(ctx, instanceID, new_version)
	}
}

func (api *API) UpgradeToSpecificVersion(ctx context.Context, instanceID int, version string) (
	string, error) {

	var (
		data       map[string]any
		failed     map[string]any
		statusCode int
		path       = fmt.Sprintf("api/instances/%d/actions/upgrade-rabbitmq", instanceID)
		params     = make(map[string]any)
	)

	params["version"] = version
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s version=%s upgrade to specific version",
		path, version), params)
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "UpgradeToSpecificVersion",
		resourceName: "RabbitMQ",
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

	return api.waitUntilUpgraded(ctx, instanceID)
}

func (api *API) UpgradeToLatestVersion(ctx context.Context, instanceID int) (string, error) {
	var (
		data       map[string]any
		failed     map[string]any
		statusCode int
		path       = fmt.Sprintf("api/instances/%d/actions/upgrade-rabbitmq-erlang", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s upgrade to latest version", path))
	err := api.callWithRetry(ctx, api.sling.New().Post(path), retryRequest{
		functionName: "UpgradeToLatestVersion",
		resourceName: "RabbitMQ",
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

	return api.waitUntilUpgraded(ctx, instanceID)
}

func (api *API) waitUntilUpgraded(ctx context.Context, instanceID int) (string, error) {
	var path = fmt.Sprintf("api/instances/%d/nodes", instanceID)

	tflog.Debug(ctx, fmt.Sprintf("waiting until RabbitMQ been upgraded, method=GET path=%s", path))

	for {
		var (
			data   []map[string]any
			failed map[string]any
		)

		err := api.callWithRetry(ctx, api.sling.New().Path(path), retryRequest{
			functionName: "waitUntilUpgraded",
			resourceName: "RabbitMQ nodes",
			attempt:      1,
			sleep:        5 * time.Second,
			data:         &data,
			failed:       &failed,
		})
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

		// Check context before sleeping
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(10 * time.Second):
			// Continue polling
		}
	}
}
