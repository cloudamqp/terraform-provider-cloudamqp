package api

import (
	"errors"
	"fmt"
	"strconv"
)

type NotificationQuery struct {
	Id   string `url:"recipient_id,omitempty"`
	Type string `url:"type,omitempty"`
}

func (api *API) CreateNotification(instance_id int, params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instance_id)
	response, err := api.sling.Post(path).BodyJSON(params).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 201 {
		return nil, errors.New(fmt.Sprintf("CreateNotification failed, status: %v, message: %s", response.StatusCode, failed))
	}

	data["id"] = strconv.FormatFloat(data["id"].(float64), 'f', 0, 64)
	return data, err
}

func (api *API) ReadNotification(instance_id int, id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	params := &NotificationQuery{Id: id}
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instance_id)
	response, err := api.sling.Path(path).QueryStruct(params).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadNotification failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, err
}

func (api *API) UpdateNotification(instance_id int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instance_id)
	response, err := api.sling.Put(path).BodyJSON(params).Receive(nil, &failed)

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("UpdateNotification failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}

func (api *API) DeleteNotification(instance_id int, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/alarms/recipients", instance_id)
	response, err := api.sling.Delete(path).BodyJSON(params).Receive(nil, &failed)

	if response.StatusCode != 204 {
		return errors.New(fmt.Sprintf("DeleteNotificaion failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}
