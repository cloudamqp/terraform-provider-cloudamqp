package api

import (
	"fmt"
	"log"
	"time"
)

// ReadVersions - Read versions RabbitMQ and Erlang can upgrade to
func (api *API) ReadVersions(instanceID int) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::upgrade_rabbitmq::read_versions instance id: %d", instanceID)
	path := fmt.Sprintf("api/instances/%d/actions/new-rabbitmq-erlang-versions", instanceID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadVersions failed, status: %v, message: %s", response.StatusCode, failed)
	}
	return data, nil
}

// UpgradeRabbitMQ - Upgrade to latest possible versions for both RabbitMQ and Erlang.
func (api *API) UpgradeRabbitMQ(instanceID int) (string, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	path := fmt.Sprintf("api/instances/%d/actions/upgrade-rabbitmq-erlang", instanceID)
	response, err := api.sling.New().Post(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::upgrade_rabbitmq::upgrade_rabbitmq_mq data: %v, status code: %v", data, response.StatusCode)
	if err != nil {
		return "", err
	}
	if response.StatusCode == 200 {
		return "Already at highest possible version", nil
	} else if response.StatusCode != 202 {
		return "", fmt.Errorf("upgrade RabbitMQ failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitUntilUpgraded(instanceID)
}

func (api *API) waitUntilUpgraded(instanceID int) (string, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})

	for {
		time.Sleep(30 * time.Second)
		path := fmt.Sprintf("api/instances/%v/nodes", instanceID)
		_, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			log.Printf("[ERROR] go-api::upgrade_rabbitmq::waitUntilUpgraded error: %v", err)
			return "", err
		}
		log.Printf("[DEBUG] go-api::upgrade_rabbitmq::waitUntilUpgraded numberOfNodes: %v", len(data))
		log.Printf("[DEBUG] go-api::upgrade_rabbitmq::waitUntilUpgraded data: %v", data)
		ready := true
		for _, node := range data {
			log.Printf("[DEBUG] go-api::upgrade_rabbitmq::waitUntilUpgraded ready: %v, configured: %v",
				ready, node["configured"])
			ready = ready && node["configured"].(bool)
		}
		log.Printf("[DEBUG] go-api::upgrade_rabbitmq::waitUntilUpgraded ready: %v", ready)
		if ready {
			return "", nil
		}
	}
}
