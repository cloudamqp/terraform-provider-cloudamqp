package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (api *API) ReadRabbitConfiguration(instanceID int) (map[string]interface{}, error) {
	// Initiale values, 5 attempts and 20 second sleep
	return api.readRabbitConfigurationWithRetry(instanceID, 5, 20)
}

func (api *API) readRabbitConfigurationWithRetry(instanceID, attempts, sleep int) (map[string]interface{}, error) {
	var data map[string]interface{}
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/config", instanceID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::rabbit-configuration#readWithRetry data: %v", data)

	if err != nil {
		return nil, err
	}

	switch {
	case response.StatusCode == 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			if attempts--; attempts > 0 {
				log.Printf("[INFO] go-api::rabbit-configuration#readWithRetry Timeout talking to backend "+
					"attempts left %d and retry in %d seconds", attempts, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				return api.readRabbitConfigurationWithRetry(instanceID, attempts, 2*sleep)
			}
			return nil, fmt.Errorf("ReadWithRetry failed, status: %v, message: %s", response.StatusCode, failed)
		}
	}
	return data, nil
}

func (api *API) UpdateRabbitConfiguration(instanceID int, params map[string]interface{}) error {
	// Initiale values, 5 attempts and 20 second sleep
	return api.updateRabbitConfigurationWithRetry(instanceID, 5, 20, params)
}

func (api *API) updateRabbitConfigurationWithRetry(instanceID, attempts, sleep int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	path := fmt.Sprintf("api/instances/%d/config", instanceID)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}

	switch {
	case response.StatusCode == 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			if attempts--; attempts > 0 {
				log.Printf("[INFO] go-api::rabbit-configuration#updateWithRetry Timeout talking to backend "+
					"attempts left %d and retry in %d seconds", attempts, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				return api.updateRabbitConfigurationWithRetry(instanceID, attempts, 2*sleep, params)
			} else {
				return fmt.Errorf("UpdateWithRetry failed, status: %v, message: %s", response.StatusCode, failed)
			}
		}
	}
	return nil
}

func (api *API) DeleteRabbitConfiguration() error {
	return nil
}
