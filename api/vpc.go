package api

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (api *API) waitUntilVpcReady(vpcID string) error {
	log.Printf("[DEBUG] go-api::vpc::waitUntilVpcReady waiting")
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		path := fmt.Sprintf("api/vpcs/%s/vpc-peering/info", vpcID)
		response, err := api.sling.New().Get(path).Receive(&data, &failed)

		if err != nil {
			return err
		}
		if response.StatusCode == 400 {
			log.Printf("[WARN] go-api::vpc::waitUntilVpcReady status: %v, message: %s", response.StatusCode, failed)
		} else if response.StatusCode != 200 {
			return fmt.Errorf("waitUntilReady failed, status: %v, message: %s", response.StatusCode, failed)
		} else if response.StatusCode == 200 {
			return nil
		}

		time.Sleep(10 * time.Second)
	}
}

func (api *API) readVpcName(vpcID string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	path := fmt.Sprintf("api/vpcs/%s/vpc-peering/info", vpcID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("readVpcName failed, status: %v, message: %s", response.StatusCode, failed)
	}
	return data, nil
}

func (api *API) CreateVpcInstance(params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc::create params: %v", params)
	response, err := api.sling.New().Post("/api/vpcs").BodyJSON(params).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("CreateVpcInstance failed, status: %v, message: %s", response.StatusCode, failed)
	}

	if id, ok := data["id"]; ok {
		data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
		log.Printf("[DEBUG] go-api::vpc::create id set: %v", data["id"])
	} else {
		msg := fmt.Sprintf("go-api::vpc::create Invalid instance identifier: %v", data["id"])
		log.Printf("[ERROR] %s", msg)
		return nil, errors.New(msg)
	}

	api.waitUntilVpcReady(data["id"].(string))
	return data, nil
}

func (api *API) ReadVpcInstance(vpcID string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc::read vpc ID: %s", vpcID)

	path := fmt.Sprintf("/api/vpcs/%s", vpcID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::vpc::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadVpcInstance failed, status: %v, message: %v", response.StatusCode, failed)
	}

	data_temp, _ := api.readVpcName(vpcID)
	data["vpc_name"] = data_temp["name"]
	return data, nil
}

func (api *API) UpdateVpcInstance(vpcID string, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::instance::update vpc ID: %s, params: %v", vpcID, params)
	path := fmt.Sprintf("api/vpcs/%s", vpcID)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("UpdateInstance failed, status: %v, message: %v", response.StatusCode, failed)
	}

	return nil
}

func (api *API) DeleteVpcInstance(vpcID string) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::vpc::delete vpc ID: %s", vpcID)
	path := fmt.Sprintf("api/vpcs/%s", vpcID)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return fmt.Errorf("DeleteVpcInstance failed, status: %v, message: %v", response.StatusCode, failed)
	}

	return nil
}
