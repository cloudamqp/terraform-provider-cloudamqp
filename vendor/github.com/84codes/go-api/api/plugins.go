package api

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type PluginParams struct {
	Name    string `json:"plugin_name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

func (api *API) waitUntilPluginChanged(instance_id int, name string, enabled bool) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::plugin::waitUntilPluginChanged instance id: %v, name: %v", instance_id, name)
	time.Sleep(10 * time.Second)
	for {
		response, err := api.ReadPlugin(instance_id, name)
		log.Printf("[DEBUG] go-api::plugin::waitUntilPluginChanged response: %v", response)
		if err != nil {
			return nil, err
		}
		if response["enabled"] == enabled {
			return response, nil
		}

		time.Sleep(10 * time.Second)
	}
}

func (api *API) EnablePlugin(instance_id int, name string) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	params := &PluginParams{Name: name}
	log.Printf("[DEBUG] go-api::plugin::enable instance id: %v, params: %v", instance_id, params)
	path := fmt.Sprintf("/api/instances/%d/plugins", instance_id)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, errors.New(fmt.Sprintf("EnablePlugin failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return api.waitUntilPluginChanged(instance_id, name, true)
}

func (api *API) ReadPlugin(instance_id int, plugin_name string) (map[string]interface{}, error) {
	var data []map[string]interface{}
	log.Printf("[DEBUG] go-api::plugin::read instance id: %v, name: %v", instance_id, plugin_name)
	data, err := api.ReadPlugins(instance_id)
	if err != nil {
		return nil, err
	}

	for _, plugin := range data {
		if plugin["name"] == plugin_name {
			log.Printf("[DEBUG] go-api::plugin::read plugin found: %v", plugin_name)
			return plugin, nil
		}
	}

	return nil, nil
}

func (api *API) ReadPlugins(instance_id int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::plugin::read instance id: %v", instance_id)
	path := fmt.Sprintf("/api/instances/%d/plugins", instance_id)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadPlugin failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) UpdatePlugin(instance_id int, params map[string]interface{}) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	pluginParams := &PluginParams{Name: params["name"].(string), Enabled: params["enabled"].(bool)}
	log.Printf("[DEBUG] go-api::plugin::update instance id: %v, params: %v", instance_id, pluginParams)
	path := fmt.Sprintf("/api/instances/%d/plugins", instance_id)
	response, err := api.sling.New().Put(path).BodyJSON(pluginParams).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, errors.New(fmt.Sprintf("UpdatePlugin failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return api.waitUntilPluginChanged(instance_id, params["name"].(string), params["enabled"].(bool))
}

func (api *API) DisablePlugin(instance_id int, name string) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::plugin::disable instance id: %v, name: %v", instance_id, name)
	path := fmt.Sprintf("/api/instances/%d/plugins/%s", instance_id, name)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, errors.New(fmt.Sprintf("DisablePlugin failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return api.waitUntilPluginChanged(instance_id, name, false)
}
