package api

import (
	"strconv"
	"time"
	"fmt"
)

type NotificationQuery struct {
	Id 	string 			`url:"id,omitempty"`
	Type string 		`url:"type,omitempty"`
}

func (api *API) waitForNotification(instance_id int, recipient_id, recipient_type string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	params := &NotificationQuery{ Id: recipient_id, Type: recipient_type}
	for {
		path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instance_id)
		_, err := api.sling.Get(path).QueryStruct(params).ReceiveSuccess(&data)
		if err != nil {
			return nil, err
		}
		if data["type"] == recipient_type {
			data["id"] = recipient_id
			return data, nil
		}
		time.Sleep(10 * time.Second)
	}
}

func (api *API) CreateNotification(instance_id int, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instance_id)
	_, err := api.sling.Post(path).BodyJSON(params).ReceiveSuccess(&data)
	time.Sleep(10 * time.Second)
	if err != nil {
		return nil, err
	}
	string_id := strconv.Itoa(int(data["id"].(float64)))
	api.waitForNotification(instance_id, string_id, data["type"].(string))
	data["id"] = string_id
	return data, err
}

func (api *API) ReadNotification(instance_id int, id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	params := &NotificationQuery{ Id: id }
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instance_id)
	_, err := api.sling.Path(path).QueryStruct(params).ReceiveSuccess(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (api *API) UpdateNotification(instance_id int, params map[string]interface{}) error {
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instance_id)
	_, err := api.sling.Put(path).BodyJSON(params).ReceiveSuccess(nil)
	return err
}

func (api *API) DeleteNotification(instance_id int, params map[string]interface{}) error {
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instance_id)
	_, err := api.sling.Delete(path).BodyJSON(params).ReceiveSuccess(nil)
	return err
}
