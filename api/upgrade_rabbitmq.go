package api

import (
	"fmt"
	"log"
	"time"
)

// ReadVersions - Read versions RabbitMQ and Erlang can upgrade to
func (api *API) ReadVersions(instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/new-rabbitmq-erlang-versions", instanceID)
	)

	log.Printf("[DEBUG] api::upgrade_rabbitmq#read_versions path: %s", path)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("ReadVersions failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

// UpgradeRabbitMQ - Upgrade to latest possible version or a specific available version
func (api *API) UpgradeRabbitMQ(instanceID int, current_version, new_version string) (string, error) {
	log.Printf("[DEBUG] api::upgrade_rabbitmq#upgrade_rabbitmq instanceID: %d, current_version: %s"+
		", new_version: %s", instanceID, current_version, new_version)

	// Keep old behaviour
	if current_version == "" && new_version == "" {
		return api.UpgradeToLatestVersion(instanceID)
	} else if current_version != "" {
		return api.UpgradeToLatestVersion(instanceID)
	} else {
		return api.UpgradeToSpecificVersion(instanceID, new_version)
	}
}

func (api *API) UpgradeToSpecificVersion(instanceID int, version string) (string, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/upgrade-rabbitmq", instanceID)
		params = make(map[string]any)
	)

	params["version"] = version
	log.Printf("[DEBUG] api::upgrade_rabbitmq#upgrade_to_specific_version path: %s, params: %v",
		path, params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return "", err
	}

	switch response.StatusCode {
	case 200:
		return api.waitUntilUpgraded(instanceID)
	default:
		return "", fmt.Errorf("upgrade RabbitMQ failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) UpgradeToLatestVersion(instanceID int) (string, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/upgrade-rabbitmq-erlang", instanceID)
	)

	log.Printf("[DEBUG] api::upgrade_rabbitmq#upgrade_to_latest_version path: %s", path)
	response, err := api.sling.New().Post(path).Receive(&data, &failed)
	if err != nil {
		return "", err
	}

	switch response.StatusCode {
	case 200:
		return "Already at highest possible version", nil
	case 202:
		return api.waitUntilUpgraded(instanceID)
	default:
		return "", fmt.Errorf("upgrade RabbitMQ failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) waitUntilUpgraded(instanceID int) (string, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/nodes", instanceID)
	)

	for {
		_, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return "", err
		}

		log.Printf("[DEBUG] api::upgrade_rabbitmq#waitUntilUpgraded data: %v", data)
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
