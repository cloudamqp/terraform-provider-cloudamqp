package api

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func (api *API) CreateAlarm(instanceID int, params map[string]any) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms", instanceID)
	)

	log.Printf("[DEBUG] api::alarms#create path: %s", path)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 201:
		log.Printf("[DEBUG] api::alarms#create data: %v", data)
		if id, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
		} else {
			return nil, fmt.Errorf("create alarm failed, invalid alarm identifier: %v", data["id"])
		}
		return data, err
	default:
		return nil,
			fmt.Errorf("create alarm failed, status: %d, message: %s", response.StatusCode, failed)
	}
}

func (api *API) ReadAlarm(instanceID int, alarmID string) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, alarmID)
	)

	log.Printf("[DEBUG] api::alarms#read path: %s", path)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("[DEBUG] api::alarms#read data : %v", data)
		return data, err
	default:
		return nil,
			fmt.Errorf("read alarm failed, status: %d, message: %s", response.StatusCode, failed)
	}
}

func (api *API) ListAlarms(instanceID int) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms", instanceID)
	)

	log.Printf("[DEBUG] api::alarms#list path: %s", path)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("[DEBUG] api::alarms#list data: %v", data)
		return data, err
	default:
		return nil,
			fmt.Errorf("list alarms failed, status: %d, message: %s", response.StatusCode, failed)
	}
}

func (api *API) UpdateAlarm(instanceID int, params map[string]any) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, params["id"].(string))
	)

	log.Printf("[DEBUG] api::alarms#update path: %s", path)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 201:
		return nil
	default:
		return fmt.Errorf("update alarm failed, status: %d, message: %s", response.StatusCode, failed)
	}
}

func (api *API) DeleteAlarm(instanceID int, alarmID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, alarmID)
	)

	log.Printf("[DEBUG] api::alarms::delete path: %s", path)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilAlarmDeletion(instanceID, alarmID)
	default:
		return fmt.Errorf("delete alarm failed, status: %d, message: %s", response.StatusCode, failed)
	}
}

func (api *API) waitUntilAlarmDeletion(instanceID int, alarmID string) error {
	var (
		data   map[string]any
		failed map[string]any
	)

	log.Printf("[DEBUG] api::alarms#waitUntilAlarmDeletion waiting")
	for {
		path := fmt.Sprintf("/api/instances/%d/alarms/%s", instanceID, alarmID)
		response, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			log.Printf("[DEBUG] api::alarms#waitUntilAlarmDeletion error: %v", err)
			return err
		}

		switch response.StatusCode {
		case 404:
			log.Print("[DEBUG] api::alarms#waitUntilAlarmDeletion deleted")
			return nil
		}

		time.Sleep(10 * time.Second)
	}
}
