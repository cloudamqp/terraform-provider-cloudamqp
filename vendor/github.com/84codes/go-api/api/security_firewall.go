package api

import (
	"errors"
	"fmt"
	"log"
)

func (api *API) CreateFirewallSettings(instance_id int, params []map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall::create instance id: %v, params: %v", instance_id, params)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instance_id)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 201 {
		return errors.New(fmt.Sprintf("CreateFirewallSettings failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}

func (api *API) ReadFirewallSettings(instance_id int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall::read instance id: %v", instance_id)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instance_id)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::security_firewall::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadFirewallSettings failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, err
}

func (api *API) UpdateFirewallSettings(instance_id int, params []map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall::update instance id: %v, params: %v", instance_id, params)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instance_id)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return errors.New(fmt.Sprintf("UpdateNotification failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}

func (api *API) DeleteFirewallSettings(instance_id int) error {
	var params [1]map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::security_firewall::delete instance id: %v", instance_id)
	path := fmt.Sprintf("/api/instances/%d/security/firewall", instance_id)

	// Use default firewall rule and update firewall upon delete.
	params[0] = DefaultFirewallSettings()
	log.Printf("[DEBUG] go-api::security_firewall::delete default firewall: %v", params[0])
	response, err := api.sling.New().Delete(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return errors.New(fmt.Sprintf("DeleteNotificaion failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}

func DefaultFirewallSettings() map[string]interface{} {
	defaultRule := map[string]interface{}{
		"services": []string{"AMQP", "AMQPS", "STOMP", "STOMPS", "MQTT", "MQTTS"},
		"ports":    []int{},
		"ip":       "0.0.0.0/0",
	}
	return defaultRule
}
