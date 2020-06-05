package api

import (
	"fmt"
	"log"
)

func (api *API) CreateFirewallSettings(instanceID int, params []map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall::create instance ID: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 201 {
		return fmt.Errorf("CreateFirewallSettings failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return err
}

func (api *API) ReadFirewallSettings(instanceID int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall::read instance id: %v", instanceID)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::security_firewall::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadFirewallSettings failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, err
}

func (api *API) UpdateFirewallSettings(instanceID int, params []map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall::update instance id: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return fmt.Errorf("UpdateNotification failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return err
}

func (api *API) DeleteFirewallSettings(instanceID int) error {
	var params [1]map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall::delete instance id: %v", instanceID)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)

	// Use default firewall rule and update firewall upon delete.
	params[0] = DefaultFirewallSettings()
	log.Printf("[DEBUG] go-api::security_firewall::delete default firewall: %v", params[0])
	response, err := api.sling.New().Delete(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return fmt.Errorf("DeleteNotificaion failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return err
}

func DefaultFirewallSettings() map[string]interface{} {
	defaultRule := map[string]interface{}{
		"services":    []string{"AMQP", "AMQPS", "STOMP", "STOMPS", "MQTT", "MQTTS"},
		"ports":       []int{},
		"ip":          "0.0.0.0/0",
		"description": "Default",
	}
	return defaultRule
}
