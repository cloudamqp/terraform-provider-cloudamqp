package api

import (
	"strconv"
	"time"
	"fmt"
)

type AlarmQuery struct {
	AlarmId 	string 		`url:"alarm_id,omitempty"`
	AlarmType string 		`url:"alarm_type,omitempty"`
}

func (api *API) waitForAlarm(instance_id int, alarm_id, alarm_type string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	params := &AlarmQuery{ AlarmId: alarm_id, AlarmType: alarm_type }
	for {
		path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
		_, err := api.sling.Get(path).QueryStruct(params).ReceiveSuccess(&data)
		if err != nil {
			return nil, err
		}
		if data["type"] == alarm_type {
			data["id"] = alarm_id
			return data, nil
		}
		time.Sleep(10 * time.Second)
	}
}

func (api *API) CreateAlarm(instance_id int, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	_, err := api.sling.Post(path).BodyJSON(params).ReceiveSuccess(&data)
	if err != nil {
		return nil, err
	}
	string_id := strconv.Itoa(int(data["id"].(float64)))
	api.waitForAlarm(instance_id, string_id, data["type"].(string))
	data["id"] = string_id
	data["alarm_id"] = string_id
	return data, err
}

func (api *API) ReadAlarm(instance_id int, id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	params := &AlarmQuery{ AlarmId: id }
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	_, err := api.sling.Get(path).QueryStruct(params).ReceiveSuccess(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (api *API) UpdateAlarm(instance_id int, params map[string]interface{}) error {
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	_, err := api.sling.Put(path).BodyJSON(params).ReceiveSuccess(nil)
	if err != nil {
		return err
	}
 	api.waitForAlarm(instance_id, params["alarm_id"].(string), params["type"].(string))
	return err
}

func (api *API) DeleteAlarm(instance_id int, params map[string]interface{}) error {
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	_, err := api.sling.Delete(path).BodyJSON(params).ReceiveSuccess(nil)
	return err
}
