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

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read RabbitMQ versions, status=%d message=%s ",
			response.StatusCode, failed)
	}
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
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/upgrade-rabbitmq", instanceID)
		params = make(map[string]any)
	)

	params["version"] = version
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s version=%s upgrade to specific version",
		path, version), params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return "", err
	}

	switch response.StatusCode {
	case 200:
		return api.waitUntilUpgraded(ctx, instanceID)
	default:
		return "", fmt.Errorf("failed to upgrade specific RabbitMQ version, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) UpgradeToLatestVersion(ctx context.Context, instanceID int) (string, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/upgrade-rabbitmq-erlang", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s upgrade to latest version", path))
	response, err := api.sling.New().Post(path).Receive(&data, &failed)
	if err != nil {
		return "", err
	}

	switch response.StatusCode {
	case 200:
		return "Already at highest possible version", nil
	case 202:
		return api.waitUntilUpgraded(ctx, instanceID)
	default:
		return "", fmt.Errorf("failed to upgrade to latest RabbitMQ version, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) waitUntilUpgraded(ctx context.Context, instanceID int) (string, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/nodes", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s waiting until RabbitMQ been upgraded", path))
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
