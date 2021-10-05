package api

import (
	"fmt"
	"log"
	"time"
)

func (api *API) waitUntilCustomDomainConfigured(instanceID int, configured bool) (map[string]interface{}, error) {
	for {
		response, err := api.ReadCustomDomain(instanceID)

		if err != nil {
			return nil, err
		}

		if response["configured"] == configured {
			return response, nil
		}

		log.Printf("[DEBUG] go-api::custom_domain#waitUntilCustomDomainConfigured: still waiting, response: %s", response)
		time.Sleep(1 * time.Second)
	}
}

func (api *API) CreateCustomDomain(instanceID int, hostname string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::custom_domain::create custom domain ID: %v, hostname: %v", instanceID, hostname)

	failed := make(map[string]interface{})
	params := make(map[string]string)
	params["hostname"] = hostname
	path := fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 202 {
		return nil, fmt.Errorf("CreateCustomDomain failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitUntilCustomDomainConfigured(instanceID, true)
}

func (api *API) ReadCustomDomain(instanceID int) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::custom_domain#read instanceID: %v", instanceID)

	failed := make(map[string]interface{})
	data := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::custom_domain::read data: %v", data)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == 200 {
		return data, nil
	} else {
		return nil, fmt.Errorf("ReadCustomDomain failed, status: %v, message: %s", response.StatusCode, failed)
	}
}

func (api *API) UpdateCustomDomain(instanceID int, hostname string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::custom_domain#update instanceID: %v", instanceID)

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

func (api *API) DeleteCustomDomain(instanceID int) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::custom_domain#delete instanceID: %v", instanceID)

	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == 200 {
		// no custom domain configured
		return nil, nil
	} else if response.StatusCode == 202 {
		// wait until the remove completed successfully
		return api.waitUntilCustomDomainConfigured(instanceID, false)
	} else {
		return nil, fmt.Errorf("DeleteCustomDomain failed, status: %v, message: %s", response.StatusCode, failed)
	}
}
