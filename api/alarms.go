package api

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

type AlarmQuery struct {
	AlarmId   string `url:"alarm_id,omitempty"`
	AlarmType string `url:"alarm_type,omitempty"`
}

func (api *API) CreateAlarm(instance_id int, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::alarm::create instance id: %v, params: %v", instance_id, params)
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::alarm::create data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 201 {
		return nil, errors.New(fmt.Sprintf("CreateAlarm failed, status: %v, message: %s", response.StatusCode, failed))
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

func (api *API) ReadAlarm(instance_id int, id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	params := &AlarmQuery{AlarmId: id}
	log.Printf("[DEBUG] go-api::alarm::read instance id: %v, alarm id: %v", instance_id, id)
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.New().Get(path).QueryStruct(params).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::alarm::read data : %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadAlarm failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, err
}

func (api *API) ReadAlarms(instance_id int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::alarm::read instance id: %v", instance_id)
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::alarm::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Alarms::ReadAlarms failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, err
}

func (api *API) UpdateAlarm(instance_id int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::alarm::update instance id: %v, params: %v", instance_id, params)
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 201 {
		return errors.New(fmt.Sprintf("Alarms::UpdateAlarm failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}

func (api *API) DeleteAlarm(instance_id int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::alarm::delete instance id: %v, params: %v", instance_id, params)
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.New().Delete(path).BodyJSON(params).Receive(nil, &failed)

	if response.StatusCode != 204 {
		return errors.New(fmt.Sprintf("Alarm::DeleteAlarm failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}
