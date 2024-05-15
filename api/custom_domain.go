package api

import (
	"fmt"
	"log"
	"time"
)

func (api *API) waitUntilCustomDomainConfigured(instanceID int, configured bool) (
	map[string]any, error) {

	for {
		response, err := api.ReadCustomDomain(instanceID)

		if err != nil {
			return nil, err
		}

		if response["configured"] == configured {
			return response, nil
		}

		log.Printf("[DEBUG] api::custom_domain#waitUntilCustomDomainConfigured: still waiting, "+
			"response: %s", response)
		time.Sleep(1 * time.Second)
	}
}

func (api *API) CreateCustomDomain(instanceID int, hostname string) (map[string]any, error) {
	var (
		failed map[string]any
		params = make(map[string]string)
		path   = fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	)

	log.Printf("[DEBUG] api::custom_domain#create custom domain ID: %d, hostname: %s",
		instanceID, hostname)
	params["hostname"] = hostname
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 202:
		return api.waitUntilCustomDomainConfigured(instanceID, true)
	default:
		return nil, fmt.Errorf("create custom domain failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) ReadCustomDomain(instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	)

	log.Printf("[DEBUG] api::custom_domain#read instanceID: %d", instanceID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		log.Printf("[DEBUG] go-api::custom_domain::read data: %v", data)
		return data, nil
	default:
		return nil, fmt.Errorf("read custom domain failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) UpdateCustomDomain(instanceID int, hostname string) (map[string]any, error) {
	log.Printf("[DEBUG] go-api::custom_domain#update instanceID: %d, hostname: %s",
		instanceID, hostname)

	// delete and wait
	_, err := api.DeleteCustomDomain(instanceID)
	if err != nil {
		return nil, err
	}
	_, err = api.waitUntilCustomDomainConfigured(instanceID, false)
	if err != nil {
		return nil, err
	}

	// create and wait
	_, err = api.CreateCustomDomain(instanceID, hostname)
	if err != nil {
		return nil, err
	}
	return api.waitUntilCustomDomainConfigured(instanceID, true)
}

func (api *API) DeleteCustomDomain(instanceID int) (map[string]any, error) {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	)

	log.Printf("[DEBUG] go-api::custom_domain#delete instanceID: %v", instanceID)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		// no custom domain configured
		return nil, nil
	case 202:
		// wait until the remove completed successfully
		return api.waitUntilCustomDomainConfigured(instanceID, false)
	default:
		return nil, fmt.Errorf("delete custom domain failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}
