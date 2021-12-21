package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (api *API) waitUntilFirewallConfigured(instanceID int) ([]map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::security_firewall::waitUntilFirewallConfigured waiting")
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/security/firewall/configured", instanceID)
	for {
		response, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return nil, err
		}
		if response.StatusCode == 200 {
			return data, nil
		} else if response.StatusCode == 400 {
			log.Printf("[DEBUG] go-api::security_firewall#waitUntilFirewallConfigured: The cluster is unavailable, firewall configuring")
		} else {
			return nil, fmt.Errorf("waitUntilReady failed, status: %v, message: %s", response.StatusCode, failed)
		}

		time.Sleep(30 * time.Second)
	}
}

func (api *API) CreateFirewallSettings(instanceID int, params []map[string]interface{}) ([]map[string]interface{}, error) {
	// Initiale values, 10 attempts, 30 second sleep
	err := api.createFirewallSettingsWithReply(instanceID, params, 10, 30)
	if err != nil {
		return nil, err
	}
	api.waitUntilFirewallConfigured(instanceID)
	return api.ReadFirewallSettings(instanceID)
}

func (api *API) createFirewallSettingsWithReply(instanceID int, params []map[string]interface{}, attempts int, sleep int) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall::create instance ID: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}

	statusCode := response.StatusCode
	log.Printf("[DEBUG] go-api::security_firewall::create statusCode: %d", statusCode)
	switch {
	case statusCode == 400:
		if strings.Compare(failed["message"].(string), "Your new firewall rules have not finished "+
			"configuring yet,try again in a few minutes") == 0 {
			if attempts--; attempts > 0 {
				log.Printf("[INFO] go-api::security_firewall::create Firewall not finished configuring "+
					"attempts left %d and retry in %d seconds", attempts, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				return api.createFirewallSettingsWithReply(instanceID, params, attempts, sleep)
			} else {
				return fmt.Errorf("Create new firewall rules failed, status: %v, message: %s", response.StatusCode, failed)
			}
		}
	}
	return nil
}

func (api *API) ReadFirewallSettings(instanceID int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall#read instanceID: %v", instanceID)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
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

func (api *API) UpdateFirewallSettings(instanceID int, params []map[string]interface{}) ([]map[string]interface{}, error) {
	// Initiale values, 10 attempts, 30 second sleep
	err := api.updateFirewallSettingsWithRetry(instanceID, params, 10, 30)
	if err != nil {
		return nil, err
	}
	api.waitUntilFirewallConfigured(instanceID)
	return api.ReadFirewallSettings(instanceID)
}

func (api *API) updateFirewallSettingsWithRetry(instanceID int, params []map[string]interface{}, attempts int, sleep int) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall::update instance id: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}

	statusCode := response.StatusCode
	log.Printf("[DEBUG] go-api::security_firewall::update statusCode: %d", statusCode)
	switch {
	case statusCode == 400:
		if strings.Compare(failed["message"].(string), "Your new firewall rules have not finished "+
			"configuring yet,try again in a few minutes") == 0 {
			if attempts--; attempts > 0 {
				log.Printf("[INFO] go-api::security_firewall::update Firewall not finished configuring "+
					"attempts left %d and retry in %d seconds", attempts, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				return api.updateFirewallSettingsWithRetry(instanceID, params, attempts, sleep)
			} else {
				return fmt.Errorf("Update firewall rules failed, status: %v, message: %s", response.StatusCode, failed)
			}
		}
	case statusCode != 204:
		return fmt.Errorf("Update firewall rules failed, status: %v, message: %s", response.StatusCode, failed)
	}
	return nil
}

func (api *API) DeleteFirewallSettings(instanceID int) ([]map[string]interface{}, error) {
	var params [1]map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall::delete instance id: %v", instanceID)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)

	// Use default firewall rule and update firewall upon delete.
	params[0] = DefaultFirewallSettings()
	log.Printf("[DEBUG] go-api::security_firewall::delete default firewall: %v", params[0])
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, fmt.Errorf("DeleteNotification failed, status: %v, message: %s", response.StatusCode, failed)
	}

	api.waitUntilFirewallConfigured(instanceID)
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
