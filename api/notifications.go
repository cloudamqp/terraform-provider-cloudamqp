package api

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

func (api *API) CreateNotification(instance_id int, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::notification::create instance id: %v, params: %v", instance_id, params)
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instance_id)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::notification::create data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 201 {
		return nil, errors.New(fmt.Sprintf("CreateNotification failed, status: %v, message: %s", response.StatusCode, failed))
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

func (api *API) ReadNotification(instance_id int, id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::notification::read instance id: %v, recipient id: %v", instance_id, id)
	path := fmt.Sprintf("/api/instances/%v/alarms/recipients/%v", instance_id, id)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::notification::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadNotification failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, err
}

func (api *API) ReadNotifications(instance_id int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::ReadNotifications::read instance id: %v", instance_id)
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instance_id)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::ReadNotifications::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadNotifications failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, err
}

func (api *API) UpdateNotification(instance_id int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::notification::update instance id: %v, params: %v", instance_id, params)
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients/%v", instance_id, params["id"])
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("UpdateNotification failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}

func (api *API) DeleteNotification(instance_id int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::notification::delete instance id: %v, params: %v", instance_id, params)
	path := fmt.Sprintf("/api/instances/%v/alarms/recipients/%v", instance_id, params["id"])
	response, err := api.sling.New().Delete(path).BodyJSON(params).Receive(nil, &failed)

	if response.StatusCode != 204 {
		return errors.New(fmt.Sprintf("DeleteNotificaion failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}
