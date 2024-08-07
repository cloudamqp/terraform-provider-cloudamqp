package api

import (
	"fmt"
	"log"
	"time"
)

func (api *API) waitUntilFirewallConfigured(instanceID, attempt, sleep, timeout int) error {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall/configured", instanceID)
	)

	for {
		response, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return err
		} else if attempt*sleep > timeout {
			return fmt.Errorf("wait until firewall configured failed, reached timeout of %d seconds",
				timeout)
		}

		switch response.StatusCode {
		case 200:
			return nil
		case 400:
			log.Printf("[DEBUG] api::security_firewall#waitUntilFirewallConfigured: " +
				"The cluster is unavailable, firewall configuring")
			log.Printf("[INFO] api::security_firewall#waitUntilFirewallConfigured Attempt: %d, until "+
				"timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
		default:
			return fmt.Errorf("waitUntilReady failed, status: %d, message: %s",
				response.StatusCode, failed)
		}
	}
}

func (api *API) CreateFirewallSettings(instanceID int, params []map[string]any, sleep, timeout int) (
	[]map[string]any, error) {

	attempt, err := api.createFirewallSettingsWithRetry(instanceID, params, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}

	err = api.waitUntilFirewallConfigured(instanceID, attempt, sleep, timeout)
	if err != nil {
		return nil, err
	}

	return api.ReadFirewallSettings(instanceID)
}

func (api *API) createFirewallSettingsWithRetry(instanceID int, params []map[string]any,
	attempt, sleep, timeout int) (int, error) {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	log.Printf("[DEBUG] api::security_firewall#create path: %s", path)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return attempt, err
	} else if attempt*sleep > timeout {
		return attempt, fmt.Errorf("create firewall settings failed, reached timeout of %d seconds",
			timeout)
	}

	switch {
	case response.StatusCode == 201:
		return attempt, nil
	case response.StatusCode == 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001:
			log.Printf("[INFO] api::security_firewall#create Firewall not finished configuring "+
				"attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.createFirewallSettingsWithRetry(instanceID, params, attempt, sleep, timeout)
		case failed["error_code"].(float64) == 40002:
			return attempt, fmt.Errorf("firewall rules validation failed due to: %s",
				failed["error"].(string))
		}
	}
	return attempt, fmt.Errorf("create new firewall rules failed, status: %d, message: %s",
		response.StatusCode, failed)
}

func (api *API) ReadFirewallSettings(instanceID int) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("[DEBUG] api::security_firewall#read data: %v", data)
		return data, nil
	default:
		return nil, fmt.Errorf("ReadFirewallSettings failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) UpdateFirewallSettings(instanceID int, params []map[string]any,
	sleep, timeout int) ([]map[string]any, error) {

	log.Printf("[DEBUG] api::security_firewall#update instance id: %d, params: %v, sleep: %d, "+
		"timeout: %d", instanceID, params, sleep, timeout)
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

func (api *API) updateFirewallSettingsWithRetry(instanceID int, params []map[string]any,
	attempt, sleep, timeout int) (int, error) {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return attempt, err
	} else if attempt*sleep > timeout {
		return attempt, fmt.Errorf("update firewall settings failed, reached timeout of %d seconds",
			timeout)
	}

	switch response.StatusCode {
	case 204:
		return attempt, nil
	case 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001:
			log.Printf("[INFO] api::security_firewall#update Firewall not finished configuring "+
				"attempt: %d until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.updateFirewallSettingsWithRetry(instanceID, params, attempt, sleep, timeout)
		case failed["error_code"].(float64) == 40002:
			return attempt, fmt.Errorf("firewall rules validation failed due to: %s",
				failed["error"].(string))
		}
	}
	return attempt, fmt.Errorf("update firewall rules failed, status: %d, message: %v",
		response.StatusCode, failed)
}

func (api *API) DeleteFirewallSettings(instanceID, sleep, timeout int) ([]map[string]any, error) {
	log.Printf("[DEBUG] api::security_firewall#delete instance id: %d, sleep: %d, timeout: %d",
		instanceID, sleep, timeout)
	attempt, err := api.deleteFirewallSettingsWithRetry(instanceID, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}

	err = api.waitUntilFirewallConfigured(instanceID, attempt, sleep, timeout)
	if err != nil {
		return nil, err
	}

	return api.ReadFirewallSettings(instanceID)
}

func (api *API) deleteFirewallSettingsWithRetry(instanceID, attempt, sleep, timeout int) (
	int, error) {

	var (
		params [1]map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	// Use default firewall rule and update firewall upon delete.
	params[0] = DefaultFirewallSettings()
	log.Printf("[DEBUG] api::security_firewall#delete default firewall: %v", params[0])
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return attempt, err
	} else if attempt*sleep > timeout {
		return attempt, fmt.Errorf("delete firewall settings failed, reached timeout of %d seconds",
			timeout)
	}

	switch response.StatusCode {
	case 204:
		return attempt, nil
	case 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001:
			log.Printf("[INFO] api::security_firewall#delete Firewall not finished configuring "+
				"attempt: %d until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.deleteFirewallSettingsWithRetry(instanceID, attempt, sleep, timeout)
		case failed["error_code"].(float64) == 40002:
			return attempt, fmt.Errorf("firewall rules validation failed due to: %s",
				failed["error"].(string))
		}
	}
	return attempt, fmt.Errorf("delete firewall rules failed, status: %d, message: %s",
		response.StatusCode, failed)
}

func DefaultFirewallSettings() map[string]any {
	defaultRule := map[string]any{
		"services": []string{"AMQP", "AMQPS", "STOMP", "STOMPS", "MQTT", "MQTTS", "HTTPS", "STREAM",
			"STREAM_SSL"},
		"ports":       []int{},
		"ip":          "0.0.0.0/0",
		"description": "Default",
	}
	return defaultRule
}
