package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type PluginParams struct {
	Name    string `json:"plugin_name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

func (api *API) waitUntilPluginChanged(instanceID int, pluginName string, enabled bool) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::plugin::waitUntilPluginChanged instance id: %v, name: %v", instanceID, pluginName)
	for {
		time.Sleep(10 * time.Second)
		response, err := api.ReadPlugin(instanceID, pluginName)
		log.Printf("[DEBUG] go-api::plugin::waitUntilPluginChanged response: %v", response)
		if err != nil {
			return nil, err
		}
		if response["required"] != nil && response["required"] != false {
			return response, nil
		}
		if response["enabled"] == enabled {
			return response, nil
		}
	}
}

func (api *API) EnablePlugin(instanceID int, pluginName string) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	params := &PluginParams{Name: pluginName}
	log.Printf("[DEBUG] go-api::plugin::enable instance id: %v, params: %v", instanceID, pluginName)
	path := fmt.Sprintf("/api/instances/%d/plugins?async=true", instanceID)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, fmt.Errorf("EnablePlugin failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitUntilPluginChanged(instanceID, pluginName, true)
}

func (api *API) ReadPlugin(instanceID int, pluginName string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::plugin::read instance id: %v, name: %v", instanceID, pluginName)
	data, err := api.ReadPlugins(instanceID)
	if err != nil {
		return nil, err
	}

	for _, plugin := range data {
		if plugin["name"] == pluginName {
			log.Printf("[DEBUG] go-api::plugin::read plugin found: %v", pluginName)
			return plugin, nil
		}
	}

	return nil, nil
}

func (api *API) ReadPlugins(instanceID int) ([]map[string]interface{}, error) {
	// Initiale values, 5 attempts and 20 second sleep
	return api.readPluginsWithRetry(instanceID, 5, 20)
}

func (api *API) readPluginsWithRetry(instanceID, attempts, sleep int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::plugin::readWithRetry instance id: %v", instanceID)
	path := fmt.Sprintf("/api/instances/%d/plugins", instanceID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}

	statusCode := response.StatusCode
	log.Printf("[DEBUG] go-api::plugins::readWithRetry statusCode: %d", statusCode)
	switch {
	case statusCode == 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			if attempts--; attempts > 0 {
				log.Printf("[INFO] go-api::plugin::readWithRetry Timeout talking to backend "+
					"attempts left %d and retry in %d seconds", attempts, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				return api.readPluginsWithRetry(instanceID, attempts, 2*sleep)
			} else {
				return nil, fmt.Errorf("ReadWithRetry failed, status: %v, message: %s", response.StatusCode, failed)
			}
		}
	}
	return data, nil
}

func (api *API) UpdatePlugin(instanceID int, params map[string]interface{}) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	pluginParams := &PluginParams{Name: params["name"].(string), Enabled: params["enabled"].(bool)}
	log.Printf("[DEBUG] go-api::plugin::update instance ID: %v, params: %v", instanceID, pluginParams)
	path := fmt.Sprintf("/api/instances/%d/plugins?async=true", instanceID)
	response, err := api.sling.New().Put(path).BodyJSON(pluginParams).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, fmt.Errorf("UpdatePlugin failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitUntilPluginChanged(instanceID, params["name"].(string), params["enabled"].(bool))
}

func (api *API) DisablePlugin(instanceID int, pluginName string) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::plugin::disable instance id: %v, name: %v", instanceID, pluginName)
	path := fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, pluginName)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, fmt.Errorf("DisablePlugin failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitUntilPluginChanged(instanceID, pluginName, false)
}

func (api *API) DeletePlugin(instanceID int, pluginName string) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::plugin::delete instance: %v, name: %v", instanceID, pluginName)
	path := fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, pluginName)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return fmt.Errorf("DeletePlugin failed, status: %v, message: %s", response.StatusCode, failed)
	}

	_, err = api.waitUntilPluginChanged(instanceID, pluginName, false)
	return err
}
