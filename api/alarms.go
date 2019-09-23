package api

import (
	"strconv"
	"fmt"
	"errors"
)

type AlarmQuery struct {
	AlarmId 	string 		`url:"alarm_id,omitempty"`
	AlarmType string 		`url:"alarm_type,omitempty"`
}

func (api *API) CreateAlarm(instance_id int, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.Post(path).BodyJSON(params).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 201 {
		return nil, errors.New(fmt.Sprintf("CreateAlarm failed, status: %v, message: %s", response.StatusCode, failed))
	}

	data["id"] = strconv.FormatFloat(data["id"].(float64), 'f', 0, 64 )
	return data, err
}

func (api *API) ReadAlarm(instance_id int, id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	params := &AlarmQuery{ AlarmId: id }
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.Get(path).QueryStruct(params).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadAlarm failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, err
}

func(api *API) ReadAlarms(instance_id int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.Get(path).Receive(&data, &failed)

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
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return errors.New(fmt.Sprintf("Alarms::UpdateAlarm failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}

func (api *API) DeleteAlarm(instance_id int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.Delete(path).BodyJSON(params).Receive(nil, &failed)

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Alarm::DeleteAlarm failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}
