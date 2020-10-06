package api

import (
	"fmt"
	"log"
	"strconv"
)

// CreateWebhook - create a webhook for a vhost and a specific qeueu
func (api *API) CreateWebhook(instanceID int, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::webhook::create params: %v", params)
	path := fmt.Sprintf("/api/instances/%d/webhooks", instanceID)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::webhook::create response data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 201 {
		return nil, fmt.Errorf(fmt.Sprintf("CreateWebhook failed, status: %v, message: %s", response.StatusCode, failed))
	}

	if v, ok := data["id"]; ok {
		data["id"] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
	} else {
		msg := fmt.Sprintf("go-api::webhook::create Invalid webhook identifier: %v", data["id"])
		log.Printf("[ERROR] %s", msg)
		return nil, fmt.Errorf(msg)
	}

	return data, err
}

// ReadWebhook - retrieves a specific webhook for an instance
func (api *API) ReadWebhook(instanceID int, webhookID string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::webhook::read instance ID: %d, webhookID: %s", instanceID, webhookID)
	path := fmt.Sprintf("/api/instances/%d/webhooks/%s", instanceID, webhookID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadWebhook failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, err
}

// ReadWebhooks - retrieves all webhooks for an instance.
func (api *API) ReadWebhooks(instanceID int) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::webhook::read instance ID: %d", instanceID)
	path := fmt.Sprintf("/api/instances/%d/webhooks", instanceID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadWebhooks failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, err
}

// DeleteWebhook - removes a specific webhook for an instance
func (api *API) DeleteWebhook(instanceID int, webhookID string) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::webhook::delete instance ID: %d, webhookID: %s", instanceID, webhookID)
	path := fmt.Sprintf("/api/instances/%d/webhooks/%s", instanceID, webhookID)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if response.StatusCode != 204 {
		return fmt.Errorf("DeleteWebhook failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return err
}
