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

// UpgradeRabbitMQ - Upgrade to latest possible versions for both RabbitMQ and Erlang.
func (api *API) UpgradeRabbitMQ(instanceID int) (string, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/actions/upgrade-rabbitmq-erlang", instanceID)
	)

	log.Printf("[DEBUG] api::upgrade_rabbitmq#upgrade_rabbitmq path: %s", path)
	response, err := api.sling.New().Post(path).Receive(&data, &failed)
	if err != nil {
		return "", err
	}

	log.Printf("[DEBUG] api::upgrade_rabbitmq::upgrade_rabbitmq_mq data: %v, status code: %d",
		data, response.StatusCode)

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
		log.Printf("[DEBUG] api::upgrade_rabbitmq#waitUntilUpgraded numberOfNodes: %d", len(data))
		log.Printf("[DEBUG] api::upgrade_rabbitmq#waitUntilUpgraded data: %v", data)
		ready := true
		for _, node := range data {
			log.Printf("[DEBUG] api::upgrade_rabbitmq#waitUntilUpgraded ready: %v, configured: %v",
				ready, node["configured"])
			ready = ready && node["configured"].(bool)
		}
		log.Printf("[DEBUG] api::upgrade_rabbitmq#waitUntilUpgraded ready: %v", ready)
		if ready {
			return "", nil
		}
		time.Sleep(10 * time.Second)
	}
}
