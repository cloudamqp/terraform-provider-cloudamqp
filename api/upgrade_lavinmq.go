package api

import (
	"fmt"
	"log"
	"time"
)

// ReadVersions - Read versions LavinMQ can upgrade to
func (api *API) ReadLavinMQVersions(instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/new-lavinmq-versions", instanceID)
	)

	log.Printf("[DEBUG] api::upgrade_lavinmq#read_versions path: %s", path)
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

// UpgradeLavinMQ - Upgrade to latest possible version or a specific available version
func (api *API) UpgradeLavinMQ(instanceID int, new_version string) (string, error) {
	log.Printf("[DEBUG] api::upgrade_lavinmq#upgrade_lavinmq instanceID: %d"+
		", new_version: %s", instanceID, new_version)

	if new_version == "" {
		return api.UpgradeToLatestLavinMQVersion(instanceID)
	} else {
		return api.UpgradeToSpecificLavinMQVersion(instanceID, new_version)
	}
}

func (api *API) UpgradeToSpecificLavinMQVersion(instanceID int, version string) (string, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/upgrade-lavinmq", instanceID)
		params = make(map[string]any)
	)

	params["version"] = version
	log.Printf("[DEBUG] api::upgrade_lavinmq#upgrade_to_specific_version path: %s, params: %v",
		path, params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return "", err
	}

	switch response.StatusCode {
	case 200:
		return api.waitUntilLavinMQUpgraded(instanceID)
	default:
		return "", fmt.Errorf("upgrade LavinMQ failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) UpgradeToLatestLavinMQVersion(instanceID int) (string, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/upgrade-lavinmq", instanceID)
	)

	log.Printf("[DEBUG] api::upgrade_lavinmq#upgrade_to_latest_version path: %s", path)
	response, err := api.sling.New().Post(path).Receive(&data, &failed)
	if err != nil {
		return "", err
	}

	switch response.StatusCode {
	case 200:
		return "Already at highest possible version", nil
	case 202:
		return api.waitUntilLavinMQUpgraded(instanceID)
	default:
		return "", fmt.Errorf("upgrade LavinMQ failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) waitUntilLavinMQUpgraded(instanceID int) (string, error) {
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

		log.Printf("[DEBUG] api::upgrade_lavinmq#waitUntilUpgraded data: %v", data)
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
