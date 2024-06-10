package api

import (
	"fmt"
	"log"
	"strconv"
)

func (api *API) CreateNotification(instanceID int, params map[string]any) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients", instanceID)
	)

	log.Printf("[DEBUG] api::notification#create path: %s", path)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 201:
		log.Printf("[DEBUG] api::notification#create data: %v", data)
		if v, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		} else {
			return nil, fmt.Errorf("create notification invalid identifier: %v", data["id"])
		}
		return data, err
	default:
		return nil, fmt.Errorf("create notification failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) ReadNotification(instanceID int, recipientID string) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	log.Printf("[DEBUG] api::notification#read path: %s", path)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("[DEBUG] api::notification#read data: %v", data)
		return data, nil
	default:
		return nil, fmt.Errorf("read notification failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) ListNotifications(instanceID int) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients", instanceID)
	)

	log.Printf("[DEBUG] api::ReadNotifications#read path: %s", path)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("[DEBUG] api::ReadNotifications#read data: %v", data)
		return data, nil
	default:
		return nil, fmt.Errorf("read notification failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) UpdateNotification(instanceID int, recipientID string,
	params map[string]any) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	log.Printf("[DEBUG] api::notification#update path: %s", path)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	default:
		return fmt.Errorf("update notification failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) DeleteNotification(instanceID int, recipientID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/alarms/recipients/%s", instanceID, recipientID)
	)

	log.Printf("[DEBUG] api::notification#delete path: %s", path)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("delete notification failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}
