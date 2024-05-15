package api

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func (api *API) waitUntilVpcReady(vpcID string) error {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("api/vpcs/%s/vpc-peering/info", vpcID)
	)

	log.Printf("[DEBUG] api::vpc#waitUntilVpcReady waiting")
	for {
		response, err := api.sling.New().Get(path).Receive(&data, &failed)
		if err != nil {
			return err
		}

		switch response.StatusCode {
		case 200:
			return nil
		case 400:
			log.Printf("[WARN] go-api::vpc::waitUntilVpcReady status: %d, message: %s",
				response.StatusCode, failed)
		default:
			return fmt.Errorf("waitUntilReady failed, status: %d, message: %s",
				response.StatusCode, failed)
		}
		time.Sleep(10 * time.Second)
	}
}

func (api *API) readVpcName(vpcID string) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("api/vpcs/%s/vpc-peering/info", vpcID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("readVpcName failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) CreateVpcInstance(params map[string]interface{}) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
	)

	log.Printf("[DEBUG] api::vpc#create params: %v", params)
	response, err := api.sling.New().Post("/api/vpcs").BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		if id, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
			log.Printf("[DEBUG] api::vpc#create id set: %v", data["id"])
		} else {
			return nil, fmt.Errorf("create VPC invalid instance identifier: %v", data["id"])
		}
		api.waitUntilVpcReady(data["id"].(string))
		return data, nil
	default:
		return nil, fmt.Errorf("create VPC failed, status: %d, message: %s", response.StatusCode, failed)
	}
}

func (api *API) ReadVpcInstance(vpcID string) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/vpcs/%s", vpcID)
	)

	log.Printf("[DEBUG] api::vpc#read path: %s", path)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		data_temp, _ := api.readVpcName(vpcID)
		data["vpc_name"] = data_temp["name"]
		return data, nil
	case 410:
		log.Printf("[WARN] api::vpc::read status: 410, message: The VPC has been deleted")
		return nil, nil
	default:
		return nil, fmt.Errorf("read VPC failed, status: %d, message: %s", response.StatusCode, failed)
	}
}

func (api *API) UpdateVpcInstance(vpcID string, params map[string]interface{}) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("api/vpcs/%s", vpcID)
	)

	log.Printf("[DEBUG] api::instance#update path: %s, params: %v", path, params)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	case 410:
		log.Printf("[WARN] api::vpc#update status: 410, message: The VPC has been deleted")
		return nil
	default:
		return fmt.Errorf("update VPC failed, status: %d, message: %s", response.StatusCode, failed)
	}
}

func (api *API) DeleteVpcInstance(vpcID string) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("api/vpcs/%s", vpcID)
	)

	log.Printf("[DEBUG] api::vpc#delete path: %s", path)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	case 410:
		log.Printf("[WARN] api::vpc#delete status: 410, message: The VPC has been deleted")
		return nil
	default:
		return fmt.Errorf("delete VPC failed, status: %d, message: %s", response.StatusCode, failed)
	}
}
