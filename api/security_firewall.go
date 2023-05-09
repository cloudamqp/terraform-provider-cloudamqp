package api

import (
	"fmt"
	"log"
	"time"
)

func (api *API) waitUntilFirewallConfigured(instanceID, attempt, sleep, timeout int) error {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/security/firewall/configured", instanceID)
	)

	for {
		response, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return err
		} else if attempt*sleep > timeout {
			return fmt.Errorf("Wait until firewall configured failed, reached timeout of %d seconds", timeout)
		}

		switch response.StatusCode {
		case 200:
			return nil
		case 400:
			log.Printf("[DEBUG] go-api::security_firewall#waitUntilFirewallConfigured: The cluster is unavailable, firewall configuring")
		default:
			return fmt.Errorf("waitUntilReady failed, status: %v, message: %s", response.StatusCode, failed)
		}

		log.Printf("[INFO] go-api::security_firewall::waitUntilFirewallConfigured The cluster is unavailable, "+
			"firewall configuring. Attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func (api *API) CreateFirewallSettings(instanceID int, params []map[string]interface{}, sleep,
	timeout int) ([]map[string]interface{}, error) {
	attempt, err := api.createFirewallSettingsWithReply(instanceID, params, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}
	err = api.waitUntilFirewallConfigured(instanceID, attempt, sleep, timeout)
	if err != nil {
		return nil, err
	}
	return api.ReadFirewallSettings(instanceID)
}

func (api *API) createFirewallSettingsWithReply(instanceID int, params []map[string]interface{},
	attempt, sleep, timeout int) (int, error) {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)
	log.Printf("[DEBUG] go-api::security_firewall::create instance ID: %v, params: %v", instanceID, params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return attempt, err
	} else if attempt*sleep > timeout {
		return attempt, fmt.Errorf("Create firewall settings failed, reached timeout of %d seconds", timeout)
	}

	switch {
	case response.StatusCode == 201:
		return attempt, nil
	case response.StatusCode == 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001:
			log.Printf("[INFO] go-api::security_firewall::create Firewall not finished configuring "+
				"attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.createFirewallSettingsWithReply(instanceID, params, attempt, sleep, timeout)
		case failed["error_code"].(float64) == 40002:
			return attempt, fmt.Errorf("Firewall rules validation failed due to: %s", failed["error"].(string))
		}
	}
	return attempt, fmt.Errorf("Create new firewall rules failed, status: %v, message: %s", response.StatusCode, failed)
}

func (api *API) ReadFirewallSettings(instanceID int) ([]map[string]interface{}, error) {
	var (
		data   []map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::security_firewall::read data: %v", data)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == 200 {
		return data, nil
	}
	return nil, fmt.Errorf("ReadFirewallSettings failed, status: %v, message: %s", response.StatusCode, failed)
}

func (api *API) UpdateFirewallSettings(instanceID int, params []map[string]interface{},
	sleep, timeout int) ([]map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::security_firewall::update instance id: %v, params: %v, sleep: %d, timeout: %d",
		instanceID, params, sleep, timeout)
	attempt, err := api.updateFirewallSettingsWithRetry(instanceID, params, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}
	err = api.waitUntilFirewallConfigured(instanceID, attempt, sleep, timeout)
	if err != nil {
		return nil, err
	}
	return api.ReadFirewallSettings(instanceID)
}

func (api *API) updateFirewallSettingsWithRetry(instanceID int, params []map[string]interface{},
	attempt, sleep, timeout int) (int, error) {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return attempt, err
	} else if attempt*sleep > timeout {
		return attempt, fmt.Errorf("Update firewall settings failed, reached timeout of %d seconds", timeout)
	}

	switch {
	case response.StatusCode == 204:
		return attempt, nil
	case response.StatusCode == 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001:
			log.Printf("[INFO] go-api::security_firewall::update Firewall not finished configuring "+
				"attempt: %d until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.updateFirewallSettingsWithRetry(instanceID, params, attempt, sleep, timeout)
		case failed["error_code"].(float64) == 40002:
			return attempt, fmt.Errorf("Firewall rules validation failed due to: %s", failed["error"].(string))
		}
	}
	return attempt, fmt.Errorf("Update firewall rules failed, status: %v, message: %v", response.StatusCode, failed)
}

func (api *API) DeleteFirewallSettings(instanceID, sleep, timeout int) ([]map[string]interface{}, error) {
	var (
		params [1]map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)
	log.Printf("[DEBUG] go-api::security_firewall::delete instance id: %v", instanceID)

	// Use default firewall rule and update firewall upon delete.
	params[0] = DefaultFirewallSettings()
	log.Printf("[DEBUG] go-api::security_firewall::delete default firewall: %v", params[0])
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, fmt.Errorf("DeleteFirewallSettings failed, status: %v, message: %s", response.StatusCode, failed)
	}

	err = api.waitUntilFirewallConfigured(instanceID, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}
	return api.ReadFirewallSettings(instanceID)
}

func DefaultFirewallSettings() map[string]interface{} {
	defaultRule := map[string]interface{}{
		"services":    []string{"AMQP", "AMQPS", "STOMP", "STOMPS", "MQTT", "MQTTS", "HTTPS", "STREAM", "STREAM_SSL"},
		"ports":       []int{},
		"ip":          "0.0.0.0/0",
		"description": "Default",
	}
	return defaultRule
}
