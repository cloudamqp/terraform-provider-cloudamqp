package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (api *API) ReadRabbitMqConfiguration(instanceID, sleep, timeout int) (map[string]interface{}, error) {
	return api.readRabbitMqConfigurationWithRetry(instanceID, 1, sleep, timeout)
}

func (api *API) readRabbitMqConfigurationWithRetry(instanceID, attempt, sleep, timeout int) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/config", instanceID)
	)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::rabbitmq-configuration#readWithRetry data: %v", data)

	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[DEBUG] go-api::rabbitmq-configuration::readWithRetry Timeout talking to backend, "+
				"attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readRabbitMqConfigurationWithRetry(instanceID, attempt, sleep, timeout)
		} else {
			break
		}
	case 503:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[DEBUG] go-api::rabbitmq-configuration::readWithRetry Timeout talking to backend, "+
				"attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readRabbitMqConfigurationWithRetry(instanceID, attempt, sleep, timeout)
		}
	}
	return nil, fmt.Errorf("read RabbitMQ configuration failed, status: %v, message: %s", response.StatusCode, failed)
}

func (api *API) UpdateRabbitMqConfiguration(instanceID int, params map[string]interface{},
	sleep, timeout int) error {
	return api.updateRabbitMqConfigurationWithRetry(instanceID, params, 1, sleep, timeout)
}

func (api *API) updateRabbitMqConfigurationWithRetry(instanceID int, params map[string]interface{},
	attempt, sleep, timeout int) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("api/instances/%d/config", instanceID)
	)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("update RabbitMQ configuraiton failed, reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 200:
		return nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[DEBUG] go-api::rabbitmq-configuration::updateWithRetry Timeout talking to backend, "+
				"attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.updateRabbitMqConfigurationWithRetry(instanceID, params, attempt, sleep, timeout)
		} else {
			break
		}
	case 503:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[DEBUG] go-api::rabbitmq-configuration::updateWithRetry Timeout talking to backend, "+
				"attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.updateRabbitMqConfigurationWithRetry(instanceID, params, attempt, sleep, timeout)
		} else {
			break
		}
	}
	return fmt.Errorf("update RabbitMQ configuration failed, status: %v, message: %s", response.StatusCode, failed)
}

func (api *API) DeleteRabbitMqConfiguration() error {
	return nil
}
