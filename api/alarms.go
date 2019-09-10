package api

import (
	"strconv"
	"time"
	"fmt"
	"errors"
	"encoding/json"
)

type AlarmQuery struct {
	AlarmId 	string 		`url:"alarm_id,omitempty"`
	AlarmType string 		`url:"alarm_type,omitempty"`
}

type Alarm struct {
	Id 								int 			`json:"id"`
	Type  						string 		`json:"type"`
	ValueThreshold 		int 			`json:"value_threshold"`
	TimeThreshold 		int 			`json:"time_threshold"`
	VhostRegex 				string 		`json:"vhost_regex"`
	QueueRegex				string 		`json:"queue_regex"`
	Notifications			[]string 	`json:"notifications"`
}

type AlarmError struct {
	Error  string  `json:"error"`
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
	alarm := new(Alarm)
	alarmError := new(AlarmError)
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.Post(path).BodyJSON(params).Receive(alarm, alarmError)
	if err != nil {
		return nil, err
	}
	if response.StatusCode > 204 {
		inrec, _ := json.Marshal(alarmError)
		return nil, errors.New(string(inrec))
	}
	string_id := strconv.Itoa(int(alarm.Id))
	api.waitForAlarm(instance_id, string_id, alarm.Type)
	inrec, _ := json.Marshal(alarm)
	json.Unmarshal(inrec, &data)
	data["id"] = string_id
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
	alarmError := new(AlarmError)
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	response, err := api.sling.Put(path).BodyJSON(params).Receive(nil, alarmError)
	if err != nil {
		return err
	}
	if response.StatusCode > 204 {
		inrec, _ := json.Marshal(alarmError)
		return errors.New(string(inrec))
	}
 	api.waitForAlarm(instance_id, params["alarm_id"].(string), params["type"].(string))
	return err
}

func (api *API) DeleteAlarm(instance_id int, params map[string]interface{}) error {
	path := fmt.Sprintf("/api/instances/%d/alarms", instance_id)
	_, err := api.sling.Delete(path).BodyJSON(params).ReceiveSuccess(nil)
	return err
}
