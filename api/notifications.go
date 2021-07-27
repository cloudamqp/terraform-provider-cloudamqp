package api

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

func (api *API) CreateNotification(instanceID int, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::notification::create instance ID: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instanceID)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::notification::create data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 201 {
		return nil, fmt.Errorf("CreateNotification failed, status: %v, message: %s", response.StatusCode, failed)
	}

	if v, ok := data["id"]; ok {
		data["id"] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		log.Printf("[DEBUG] go-api::notification::create id set: %v", data["id"])
	} else {
		msg := fmt.Sprintf("go-api::notification::create Invalid notification identifier: %v", data["id"])
		log.Printf("[ERROR] %s", msg)
		return nil, errors.New(msg)
	}

	return data, err
}

func (api *API) ReadNotification(instanceID int, recipientID string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::notification::read instance ID: %v, recipient ID: %v", instanceID, recipientID)
	path := fmt.Sprintf("/api/instances/%v/alarms/recipients/%v", instanceID, recipientID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::notification::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadNotification failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, err
}

func (api *API) ReadNotifications(instanceID int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::ReadNotifications::read instance ID: %v", instanceID)
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instanceID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::ReadNotifications::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadNotifications failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, err
}

func (api *API) UpdateNotification(instanceID int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::notification::update instance id: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients/%v", instanceID, params["id"])
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if response.StatusCode != 200 {
		return fmt.Errorf("UpdateNotification failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return err
}

func (api *API) DeleteNotification(instanceID int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::notification::delete instance id: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("/api/instances/%v/alarms/recipients/%v", instanceID, params["id"])
	response, err := api.sling.New().Delete(path).BodyJSON(params).Receive(nil, &failed)

	if response.StatusCode != 204 {
		return fmt.Errorf("DeleteNotification failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return err
}
