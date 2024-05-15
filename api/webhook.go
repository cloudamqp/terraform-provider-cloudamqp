package api

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// CreateWebhook - create a webhook for a vhost and a specific qeueu
func (api *API) CreateWebhook(instanceID int, params map[string]any,
	sleep, timeout int) (map[string]any, error) {

	return api.createWebhookWithRetry(instanceID, params, 1, sleep, timeout)
}

// createWebhookWithRetry: create webhook with retry if backend is busy.
func (api *API) createWebhookWithRetry(instanceID int, params map[string]any,
	attempt, sleep, timeout int) (map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/webhooks", instanceID)
	)

	log.Printf("[DEBUG] api::webhook#create path: %s, params: %v", path, params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("create webhook reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 201:
		log.Printf("[DEBUG] api::webhook#create response data: %v", data)
		if v, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		} else {
			return nil, fmt.Errorf("create webhook, invalid webhook identifier: %v", data["id"])
		}
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[INFO] api::webhook#create Timeout talking to backend attempt: %d, until "+
				"timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.createWebhookWithRetry(instanceID, params, attempt, sleep, timeout)
		}
	}

	return nil, fmt.Errorf("create webhook with retry failed, status: %d, message: %s",
		response.StatusCode, failed)
}

// ReadWebhook - retrieves a specific webhook for an instance
func (api *API) ReadWebhook(instanceID int, webhookID string, sleep, timeout int) (
	map[string]any, error) {

	path := fmt.Sprintf("/api/instances/%d/webhooks/%s", instanceID, webhookID)
	return api.readWebhookWithRetry(path, 1, sleep, timeout)
}

// readWebhookWithRetry: read webhook with retry if backend is busy.
func (api *API) readWebhookWithRetry(path string, attempt, sleep, timeout int) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
	)

	log.Printf("[DEBUG] api::webhook#read path: %s", path)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("read webhook reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 200:
		log.Printf("[DEBUG] api::webhook#read response data: %v", data)
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[INFO] api::webhook#read Timeout talking to backend attempt: %d, until "+
				"timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readWebhookWithRetry(path, attempt, sleep, timeout)
		}
	}

	return nil, fmt.Errorf("read webhook with retry failed, status: %d, message: %s",
		response.StatusCode, failed)
}

// ListWebhooks - list all webhooks for an instance.
func (api *API) ListWebhooks(instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/webhooks", instanceID)
	)

	log.Printf("[DEBUG] api::webhook#list path: %s", path)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("list webhooks failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

// UpdateWebhook - updates a specific webhook for an instance
func (api *API) UpdateWebhook(instanceID int, webhookID string, params map[string]any,
	sleep, timeout int) error {

	path := fmt.Sprintf("/api/instances/%d/webhooks/%s", instanceID, webhookID)
	return api.updateWebhookWithRetry(path, params, 1, sleep, timeout)
}

// updateWebhookWithRetry: update webhook with retry if backend is busy.
func (api *API) updateWebhookWithRetry(path string, params map[string]any,
	attempt, sleep, timeout int) error {

	var (
		data   = make(map[string]any)
		failed map[string]any
	)

	log.Printf("[DEBUG] api::webhook#update path: %s, params: %v", path, params)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("update webhook reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 201:
		log.Printf("[DEBUG] api::webhook#update response data: %v", data)
		return nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[INFO] api::webhook#update Timeout talking to backend, attempt: %d, until "+
				"timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.updateWebhookWithRetry(path, params, attempt, sleep, timeout)
		}
	}

	return fmt.Errorf("update webhook with retry failed, status: %d, message: %s",
		response.StatusCode, failed)
}

// DeleteWebhook - removes a specific webhook for an instance
func (api *API) DeleteWebhook(instanceID int, webhookID string, sleep, timeout int) error {
	path := fmt.Sprintf("/api/instances/%d/webhooks/%s", instanceID, webhookID)
	return api.deleteWebhookWithRetry(path, 1, sleep, timeout)
}

// deleteWebhookWithRetry: delete webhook with retry if backend is busy.
func (api *API) deleteWebhookWithRetry(path string, attempt, sleep, timeout int) error {
	var (
		failed map[string]any
	)

	log.Printf("[DEBUG] api::webhook#delete path: %s", path)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("delete webhook reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 204:
		return nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[INFO] api::webhook#delete Timeout talking to backend attempt: %d, until "+
				"timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.deleteWebhookWithRetry(path, attempt, sleep, timeout)
		}
	}

	return fmt.Errorf("delete webhook with retry failed, status: %d, message: %s",
		response.StatusCode, failed)
}
