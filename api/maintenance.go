package api

import (
	"fmt"
	"log"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
)

func (api *API) SetMaintenance(instanceID int, data model.Maintenance) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/maintenance/settings", instanceID)
	)

	log.Printf("[DEBUG] api::maintenance#set data: %v", data)

	response, err := api.sling.New().Post(path).BodyJSON(data).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	default:
		return fmt.Errorf("failed to update maintenance window, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) ReadMaintenance(instanceID int) (model.Maintenance, error) {
	var (
		data   model.Maintenance
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/maintenance/settings", instanceID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return model.Maintenance{}, err
	}

	log.Printf("[DEBUG] api::maintenance#read data: %v", data)

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return model.Maintenance{},
			fmt.Errorf("read maintenance settings failed, status: %d, message: %s",
				response.StatusCode, failed)
	}
}
