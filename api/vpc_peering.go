package api

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func (api *API) waitForPeeringStatus(instance_id int, peering_id string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::vpc_peering::waitForPeeringStatus instance id: %v, peering id: %v", instance_id, peering_id)
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		path := fmt.Sprintf("/api/instances/%v/vpc-peering/status/%v", instance_id, peering_id)
		response, err := api.sling.New().Path(path).Receive(&data, &failed)

		if err != nil {
			return nil, err
		}
		if response.StatusCode != 200 {
			return nil, errors.New(fmt.Sprintf("waitForPeeringStatus failed, status: %v, message: %s", response.StatusCode, failed))
		}
		switch data["status"] {
		case "active", "pending-acceptance":
			return data, nil
		}
		time.Sleep(10 * time.Second)
	}
}

func (api *API) ReadVpcInfo(instance_id int) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering::info instance id: %v", instance_id)
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/info", instance_id)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_peering::info data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadInfo failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) ReadVpcPeeringRequest(instance_id int, peering_id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering::request instance id: %v, peering id: %v", instance_id, peering_id)
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/request/%v", instance_id, peering_id)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_peering::request data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadRequest failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) AcceptVpcPeering(instance_id int, peering_id string) (map[string]interface{}, error) {
	_, err := api.waitForPeeringStatus(instance_id, peering_id)

	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering::accept instance id: %v, peering id: %v", instance_id, peering_id)
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/request/%v", instance_id, peering_id)
	response, err := api.sling.New().Put(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc_peering::accept data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("AcceptVpcPeering failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) RemoveVpcPeering(instance_id int, peering_id string) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc_peering::remove instance id: %v, peering id: %v", instance_id, peering_id)
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/%v", instance_id, peering_id)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return errors.New(fmt.Sprintf("RemoveVpcPeering failed, status: %v, message: %s", response.StatusCode, failed))
	}
	return nil
}
