package api

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

func (api *API) waitUntilReady(instanceID string) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%s", instanceID)
	)

	log.Printf("[DEBUG] go-api::instance::waitUntilReady waiting")

	for {
		time.Sleep(10 * time.Second)
		response, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return nil, err
		}

		switch response.StatusCode {
		case 200:
			if data["ready"] == true {
				data["id"] = instanceID
				return data, nil
			}
		default:
			return nil, fmt.Errorf("waitUntilReady failed, status: %v, message: %s",
				response.StatusCode, failed)
		}
	}
}

func (api *API) waitUntilAllNodesReady(instanceID string) error {
	var (
		data   []map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("api/instances/%v/nodes", instanceID)
	)

	for {
		time.Sleep(15 * time.Second)
		_, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] go-api::instance::waitUntilAllNodesReady numberOfNodes: %v", len(data))
		ready := true
		for _, node := range data {
			log.Printf("[DEBUG] go-api::instance::waitUntilAllNodesReady ready: %v, configured: %v",
				ready, node["configured"])
			ready = ready && node["configured"].(bool)
		}
		if ready {
			return nil
		}
	}
}

func (api *API) waitWithTimeoutUntilAllNodesConfigured(instanceID string, attempt, sleep,
	timeout int) error {

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
		return fmt.Errorf("all nodes configured timeout reached after %d seconds", timeout)
	}

	ready := true
	for _, node := range data {
		log.Printf("[DEBUG] go-api::instance::waitWithTimeoutUntilAllNodesConfigured ready: %v, configured: %v",
			ready, node["configured"])
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
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%s", instanceID)
	)

	log.Printf("[DEBUG] go-api::instance::waitUntilDeletion waiting")
	for {
		time.Sleep(10 * time.Second)
		response, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			log.Printf("[DEBUG] go-api::instance::waitUntilDeletion error: %v", err)
			return err
		}

		switch response.StatusCode {
		case 404:
			log.Print("[DEBUG] go-api::instance::waitUntilDeletion deleted")
			return nil
		case 410:
			log.Print("[DEBUG] go-api::instance::waitUntilDeletion deleted")
			return nil
		}
	}
}

func (api *API) CreateInstance(params map[string]interface{}) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
	)

	log.Printf("[DEBUG] go-api::instance::create params: %v", params)
	response, err := api.sling.New().Post("/api/instances").BodyJSON(params).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::instance::waitUntilReady data: %v", data)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		if id, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
			log.Printf("[DEBUG] go-api::instance::create id set: %v", data["id"])
		} else {
			return nil, fmt.Errorf("go-api::instance::create Invalid instance identifier: %v", data["id"])
		}
		return api.waitUntilReady(data["id"].(string))
	default:
		return nil, fmt.Errorf("create instance failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) ReadInstance(instanceID string) (map[string]interface{}, error) {
	var (
		data   map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%s", instanceID)
	)

	log.Printf("[DEBUG] go-api::instance::read instance ID: %v", instanceID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::instance::read data: %v", data)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 410:
		log.Printf("[WARN] go-api::instance::read status: 410, message: The instance has been deleted")
		return nil, nil
	default:
		return nil, fmt.Errorf("read instance failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

// TODO: Rename to ListInstances
func (api *API) ReadInstances() ([]map[string]interface{}, error) {
	var (
		data   []map[string]interface{}
		failed map[string]interface{}
	)

	response, err := api.sling.New().Get("/api/instances").Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::instance::list data: %v", data)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 410:
		log.Printf("[WARN] go-api::instance::list status: 410, message: The instance has been deleted")
		return nil, nil
	default:
		return nil, fmt.Errorf("list instances failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) UpdateInstance(instanceID string, params map[string]interface{}) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("api/instances/%v", instanceID)
	)

	log.Printf("[DEBUG] go-api::instance::update instance ID: %v, params: %v", instanceID, params)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return api.waitUntilAllNodesReady(instanceID)
	case 410:
		log.Printf("[WARN] go-api::instance::update status: 410, message: The instance has been deleted")
		return nil
	default:
		return fmt.Errorf("update instance failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) DeleteInstance(instanceID string, keep_vpc bool) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("api/instances/%s?keep_vpc=%v", instanceID, keep_vpc)
	)

	log.Printf("[DEBUG] go-api::instance::delete instance ID: %v", instanceID)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilDeletion(instanceID)
	case 410:
		log.Printf("[WARN] go-api::instance::delete status: 410, message: The instance has been deleted")
		return nil
	default:
		return fmt.Errorf("delete instance failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
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
