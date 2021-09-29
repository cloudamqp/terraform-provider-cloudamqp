package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (api *API) waitUntilPluginUninstalled(instanceID int, pluginName string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::plugin_community::waitUntilPluginUninstalled instance id: %v, name: %v", instanceID, pluginName)
	time.Sleep(10 * time.Second)
	for {
		response, err := api.ReadPlugin(instanceID, pluginName)

		if err != nil {
			return nil, err
		}
		if len(response) == 0 {
			return response, nil
		}

		time.Sleep(10 * time.Second)
	}
}

func (api *API) EnablePluginCommunity(instanceID int, pluginName string) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	params := &PluginParams{Name: pluginName}
	log.Printf("[DEBUG] go-api::plugin_community::enable instance ID: %v, name: %v", instanceID, pluginName)
	path := fmt.Sprintf("/api/instances/%d/plugins/community", instanceID)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, fmt.Errorf("EnablePluginCommunity failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitUntilPluginChanged(instanceID, pluginName, true)
}

func (api *API) ReadPluginCommunity(instanceID int, pluginName string) (map[string]interface{}, error) {
	var data []map[string]interface{}
	log.Printf("[DEBUG] go-api::plugin_community::read instance ID: %v, name: %v", instanceID, pluginName)
	data, err := api.ReadPluginsCommunity(instanceID)

	if err != nil {
		return nil, err
	}

	for _, plugin := range data {
		if plugin["name"] == pluginName {
			log.Printf("[DEBUG] go-api::plugin_community::read found plugin: %v", pluginName)
			return plugin, nil
		}
	}

	return nil, nil
}

func (api *API) ReadPluginsCommunity(instanceID int) ([]map[string]interface{}, error) {
	// Initiale values, 5 attempts and 20 second sleep
	return api.readPluginsCommunityWithRetry(instanceID, 5, 20)
}

func (api *API) readPluginsCommunityWithRetry(instanceID, attempts, sleep int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::plugin_community::readPluginsCommunityWithRetry instance id: %v", instanceID)
	path := fmt.Sprintf("/api/instances/%d/plugins/community", instanceID)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::plugin_community::readPluginsCommunityWithRetry data: %v", data)

	if err != nil {
		return nil, err
	}

	statusCode := response.StatusCode
	log.Printf("[DEBUG] go-api::plugin_community::readPluginsCommunityWithRetry statusCode: %d", statusCode)
	switch {
	case statusCode == 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			if attempts--; attempts > 0 {
				log.Printf("[INFO] go-api::plugin_community::readPluginsCommunityWithRetry Timeout talking to backend "+
					"attempts left %d and retry in %d seconds", attempts, sleep)
				time.Sleep(time.Duration(sleep) * time.Second)
				return api.readPluginsCommunityWithRetry(instanceID, attempts, 2*sleep)
			} else {
				return nil, fmt.Errorf("ReadWithRetry failed, status: %v, message: %s", response.StatusCode, failed)
			}
		}
	}
	return data, nil
}

func (api *API) UpdatePluginCommunity(instanceID int, params map[string]interface{}) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	pluginParams := &PluginParams{Name: params["name"].(string), Enabled: params["enabled"].(bool)}
	log.Printf("[DEBUG] go-api::plugin_community::update instance ID: %v, params: %v", instanceID, params)
	path := fmt.Sprintf("/api/instances/%d/plugins/community", instanceID)
	response, err := api.sling.New().Put(path).BodyJSON(pluginParams).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, fmt.Errorf("UpdatePluginCommunity failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitUntilPluginChanged(instanceID, params["name"].(string), params["enabled"].(bool))
}

func (api *API) DisablePluginCommunity(instanceID int, pluginName string) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::plugin_community::disable instance ID: %v, name: %v", instanceID, pluginName)
	path := fmt.Sprintf("/api/instances/%d/plugins/community/%s", instanceID, pluginName)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, fmt.Errorf("DisablePluginCommunity failed, status: %v, message: %s", response.StatusCode, failed)
	}

	return api.waitUntilPluginUninstalled(instanceID, pluginName)
}
