package api

import (
	"fmt"
	"log"
	"strconv"
)

func (api *API) CreateNotification(instanceID int, params map[string]interface{}) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients", instanceID)
	)

	log.Printf("[DEBUG] go-api::notification::create path: %s", path)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::notification::create data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 201 {
		return nil, fmt.Errorf("create notification failed, status: %d, message: %s",
			response.StatusCode, failed)
	}

	if v, ok := data["id"]; ok {
		data["id"] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
	} else {
		return nil, fmt.Errorf("create notification invalid identifier: %v", data["id"])
	}

	return data, err
}

func (api *API) ReadNotification(instanceID int, recipientID string) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	log.Printf("[DEBUG] go-api::notification::read path: %s", path)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::notification::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("read notification failed, status: %v, message: %s",
			response.StatusCode, failed)
	}

	return data, err
}

func (api *API) ReadNotifications(instanceID int) ([]map[string]interface{}, error) {
	var (
		data   []map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients", instanceID)
	)

	log.Printf("[DEBUG] go-api::ReadNotifications::read path: %s", path)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::ReadNotifications::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("read notification failed, status: %d, message: %s",
			response.StatusCode, failed)
	}

	return data, err
}

func (api *API) UpdateNotification(instanceID int, recipientID string, params map[string]interface{}) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	log.Printf("[DEBUG] go-api::notification::update path: %s", path)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if response.StatusCode != 200 {
		return fmt.Errorf("update notification failed, status: %d, message: %s",
			response.StatusCode, failed)
	}

	return err
}

func (api *API) DeleteNotification(instanceID int, recipientID string) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	log.Printf("[DEBUG] go-api::notification::delete path: %s", path)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if response.StatusCode != 204 {
		return fmt.Errorf("delete notification failed, status: %d, message: %s",
			response.StatusCode, failed)
	}

	return err
}
