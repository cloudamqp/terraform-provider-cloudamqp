package api

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (api *API) waitForPeeringStatus(instance_id int, peering_id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		path := fmt.Sprintf("/api/instances/%v/vpc-peering/status/%v", instance_id, peering_id)
		response, err := api.sling.Path(path).Receive(&data, &failed)

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
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/info", instance_id)
	response, err := api.sling.Get(path).Receive(&data, &failed)

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
	log.Printf("[DEBUG] - go-api::vpc_peering::ReadVpcPeeringRequest instance_id: %v, peering_id: %v", instance_id, peering_id)

	path := fmt.Sprintf("/api/instances/%v/vpc-peering/request/%v", instance_id, peering_id)
	response, err := api.sling.Get(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadRequest failed, status: %v, message: %s", response.StatusCode, failed))
	}

	data["id"] = strconv.FormatFloat(data["id"].(float64), 'f', 0, 64)
	return data, nil
}

func (api *API) AcceptVpcPeering(instance_id int, peering_id string) (map[string]interface{}, error) {
	_, err := api.waitForPeeringStatus(instance_id, peering_id)

	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/request/%v", instance_id, peering_id)
	response, err := api.sling.Put(path).Receive(&data, &failed)

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
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/%v", instance_id, peering_id)
	response, err := api.sling.Delete(path).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return errors.New(fmt.Sprintf("RemoveVpcPeering failed, status: %v, message: %s", response.StatusCode, failed))
	}
	return nil
}
