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

	tflog.Debug(ctx, fmt.Sprintf("request path: %s", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("failed reading LavinMQ versions, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

// UpgradeLavinMQ - Upgrade to latest possible version or a specific available version
func (api *API) UpgradeLavinMQ(ctx context.Context, instanceID int, new_version string) (
	string, error) {

	tflog.Debug(ctx, fmt.Sprintf("upgrade LavinMQ, instanceID: %d, new_version: %s",
		instanceID, new_version))

	if new_version == "" {
		return api.UpgradeToLatestLavinMQVersion(ctx, instanceID)
	} else {
		return api.UpgradeToSpecificLavinMQVersion(ctx, instanceID, new_version)
	}
}

func (api *API) UpgradeToSpecificLavinMQVersion(ctx context.Context, instanceID int, version string) (
	string, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/upgrade-lavinmq", instanceID)
		params = make(map[string]any)
	)

	params["version"] = version
	tflog.Debug(ctx, fmt.Sprintf("request path: %s, params: %v", path, params))
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return "", err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, fmt.Sprintf("data: %v", data))
		return api.waitUntilLavinMQUpgraded(ctx, instanceID)
	default:
		return "", fmt.Errorf("failed to upgrade to new LavinMQ version, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) UpgradeToLatestLavinMQVersion(ctx context.Context, instanceID int) (string, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/upgrade-lavinmq", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("request path: %s", path))
	response, err := api.sling.New().Post(path).Receive(&data, &failed)
	if err != nil {
		return "", err
	}

	switch response.StatusCode {
	case 200:
		return "Already at highest possible version", nil
	case 202:
		return api.waitUntilLavinMQUpgraded(ctx, instanceID)
	default:
		return "", fmt.Errorf("failed to upgrade to latest LavinMQ version, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) waitUntilLavinMQUpgraded(ctx context.Context, instanceID int) (string, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/nodes", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("waiting until LavinMQ been upgraded, request path: %s", path))
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
