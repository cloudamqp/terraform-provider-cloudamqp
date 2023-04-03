package api

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

func (api *API) waitUntilReady(instanceID string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::instance::waitUntilReady waiting")
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		time.Sleep(10 * time.Second)
		response, err := api.sling.New().Path("/api/instances/").Get(instanceID).Receive(&data, &failed)

		if err != nil {
			return nil, err
		}
		if response.StatusCode != 200 {
			return nil, fmt.Errorf("waitUntilReady failed, status: %v, message: %s", response.StatusCode, failed)
		}
		if data["ready"] == true {
			data["id"] = instanceID
			return data, nil
		}
	}
}

func (api *API) waitUntilAllNodesReady(instanceID string) error {
	var data []map[string]interface{}
	failed := make(map[string]interface{})

	for {
		time.Sleep(15 * time.Second)
		path := fmt.Sprintf("api/instances/%v/nodes", instanceID)
		_, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			log.Printf("[ERROR] go-api::instance::waitUntilAllNodesReady error: %v", err)
			return err
		}
		log.Printf("[DEBUG] go-api::instance::waitUntilAllNodesReady numberOfNodes: %v", len(data))
		log.Printf("[DEBUG] go-api::instance::waitUntilAllNodesReady data: %v", data)
		ready := true
		for _, node := range data {
			log.Printf("[DEBUG] go-api::instance::waitUntilAllNodesReady ready: %v, configured: %v", ready, node["configured"])
			ready = ready && node["configured"].(bool)
		}
		log.Printf("[DEBUG] go-api::instance::waitUntilAllNodesReady ready: %v", ready)
		if ready {
			return nil
		}
	}
}

func (api *API) waitWithTimeoutUntilAllNodesConfigured(instanceID string, attempt, sleep, timeout int) error {
	var (
		data   []map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("api/instances/%v/nodes", instanceID)
	)
	log.Printf("[DEBUG] go-api::instance::waitWithTimeoutUntilAllNodesConfigured not yet ready, "+
		" will try again, attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))

	_, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("All nodes configured timeout reached after %d seconds", timeout)
	}

	ready := true
	for _, node := range data {
		log.Printf("[DEBUG] go-api::instance::waitWithTimeoutUntilAllNodesConfigured ready: %v, configured: %v", ready, node["configured"])
		ready = ready && node["configured"].(bool)
	}
	log.Printf("[DEBUG] go-api::instance::waitWithTimeoutUntilAllNodesConfigured ready: %v", ready)
	if ready {
		return nil
	}
	attempt++
	time.Sleep(time.Duration(sleep) * time.Second)
	return api.waitWithTimeoutUntilAllNodesConfigured(instanceID, attempt, sleep, timeout)
}

func (api *API) waitUntilDeletion(instanceID string) error {
	log.Printf("[DEBUG] go-api::instance::waitUntilDeletion waiting")
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		time.Sleep(10 * time.Second)
		response, err := api.sling.New().Path("/api/instances/").Get(instanceID).Receive(&data, &failed)

		if err != nil {
			log.Printf("[DEBUG] go-api::instance::waitUntilDeletion error: %v", err)
			return err
		}
		if response.StatusCode == 404 {
			log.Print("[DEBUG] go-api::instance::waitUntilDeletion deleted")
			return nil
		}
	}
}

func (api *API) numberOfNodes(instanceID string) (int, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	path := fmt.Sprintf("api/instances/%v/nodes", instanceID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::instances::numberOfNodes data: %v", data)

	if err != nil {
		fmt.Errorf("[ERROR] go-api::instances::numberOfNodes error: %v", err)
		return -1, err
	}
	if response.StatusCode != 200 {
		return -1, fmt.Errorf("go-api::instances::numberOfNodes failed, status: %v, message: %v", response.StatusCode, failed)
	}

	if data["nodes"] == nil {
		log.Printf("[ERROR] go-api::instances::numberOfNodes is nil")
		return -1, fmt.Errorf("go-api::instances::numberOfNodes is nil")
	}
	nodes, _ := strconv.Atoi(data["nodes"].(string))
	return nodes, nil
}

func (api *API) CreateInstance(params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::instance::create params: %v", params)
	response, err := api.sling.New().Post("/api/instances").BodyJSON(params).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::instance::waitUntilReady data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("CreateInstance failed, status: %v, message: %s", response.StatusCode, failed)
	}

	if id, ok := data["id"]; ok {
		data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
		log.Printf("[DEBUG] go-api::instance::create id set: %v", data["id"])
	} else {
		msg := fmt.Sprintf("go-api::instance::create Invalid instance identifier: %v", data["id"])
		log.Printf("[ERROR] %s", msg)
		return nil, errors.New(msg)
	}

	return api.waitUntilReady(data["id"].(string))
}

func (api *API) ReadInstance(instanceID string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::instance::read instance ID: %v", instanceID)
	response, err := api.sling.New().Path("/api/instances/").Get(instanceID).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::instance::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadInstance failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, nil
}

func (api *API) ReadInstances() ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	response, err := api.sling.New().Get("/api/instances").Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::instance::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ReadInstances failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return data, nil
}

func (api *API) UpdateInstance(instanceID string, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::instance::update instance ID: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("api/instances/%v", instanceID)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("UpdateInstance failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitUntilAllNodesReady(instanceID)
}

func (api *API) DeleteInstance(instanceID string, keep_vpc bool) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::instance::delete instance ID: %v", instanceID)
	path := fmt.Sprintf("api/instances/%s?keep_vpc=%v", instanceID, keep_vpc)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return fmt.Errorf("DeleteInstance failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitUntilDeletion(instanceID)
}

func (api *API) UrlInformation(url string) map[string]interface{} {
	paramsMap := make(map[string]interface{})
	r := regexp.MustCompile(`^.*:\/\/(?P<username>(.*)):(?P<password>(.*))@(?P<host>(.*))\/(?P<vhost>(.*))`)
	match := r.FindStringSubmatch(url)

	for i, value := range r.SubexpNames() {
		if value == "username" {
			paramsMap["username"] = match[i]
		}
		if value == "password" {
			paramsMap["password"] = match[i]
		}
		if value == "host" {
			paramsMap["host"] = match[i]
		}
		if value == "vhost" {
			paramsMap["vhost"] = match[i]
		}
	}

	return paramsMap
}
