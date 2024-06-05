package api

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (api *API) CreateAlarm(instanceID int, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::alarm::create instance ID: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("/api/instances/%d/alarms", instanceID)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::alarm::create data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 201 {
		return nil, fmt.Errorf("CreateAlarm failed, status: %v, message: %s", response.StatusCode, failed)
	}

	if id, ok := data["id"]; ok {
		data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
		log.Printf("[DEBUG] go-api::alarm::create id set: %v", data["id"])
	} else {
		msg := fmt.Sprintf("go-api::instance::create Invalid alarm identifier: %v", data["id"])
		log.Printf("[ERROR] %s", msg)
		return nil, errors.New(msg)
	}

	return data, err
}

func (api *API) ReadAlarm(instanceID int, alarmID string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::alarm::read instance ID: %v, alarm ID: %v", instanceID, alarmID)
	path := fmt.Sprintf("/api/instances/%v/alarms/%v", instanceID, alarmID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::alarm::read data : %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadAlarm failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, err
}

func (api *API) ReadAlarms(instanceID int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::alarm::read instance ID: %v", instanceID)
	path := fmt.Sprintf("/api/instances/%d/alarms", instanceID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::alarm::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Alarms::ReadAlarms failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, err
}

func (api *API) UpdateAlarm(instanceID int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::alarm::update instance ID: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("/api/instances/%v/alarms/%v", instanceID, params["id"])
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 201 {
		return fmt.Errorf("Alarms::UpdateAlarm failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return err
}

func (api *API) DeleteAlarm(instanceID int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::alarm::delete instance id: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("/api/instances/%v/alarms/%v", instanceID, params["id"])
	response, _ := api.sling.New().Delete(path).BodyJSON(params).Receive(nil, &failed)

	if response.StatusCode != 204 {
		return fmt.Errorf("Alarm::DeleteAlarm failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitUntilAlarmDeletion(instanceID, params["id"].(string))
}

func (api *API) waitUntilAlarmDeletion(instanceID int, id string) error {
	log.Printf("[DEBUG] go-api::alarm::waitUntilAlarmDeletion waiting")
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		path := fmt.Sprintf("/api/instances/%v/alarms/%v", instanceID, id)
		response, err := api.sling.New().Path(path).Receive(&data, &failed)

		if err != nil {
			log.Printf("[DEBUG] go-api::alarm::waitUntilAlarmDeletion error: %v", err)
			return err
		}
		if response.StatusCode == 404 {
			log.Print("[DEBUG] go-api::alarm::waitUntilAlarmDeletion deleted")
			return nil
		}

		time.Sleep(10 * time.Second)
	}
}
