package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// EnablePlugin: enable a plugin on an instance.
func (api *API) EnablePlugin(instanceID int, pluginName string, sleep, timeout int) (
	map[string]interface{}, error) {

	var (
		failed map[string]interface{}
		params = make(map[string]interface{})
		path   = fmt.Sprintf("/api/instances/%d/plugins?async=true", instanceID)
	)

	params["plugin_name"] = pluginName
	log.Printf("[DEBUG] go-api::plugin::enable instance id: %v, params: %v", instanceID, params)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilPluginChanged(instanceID, pluginName, true, 1, sleep, timeout)
	default:
		return nil,
			fmt.Errorf("enable plugin failed, status: %v, message: %s", response.StatusCode, failed)
	}
}

// ReadPlugin: reads a specific plugin from an instance.
func (api *API) ReadPlugin(instanceID int, pluginName string, sleep, timeout int) (
	map[string]interface{}, error) {

	log.Printf("[DEBUG] go-api::plugin::read instance id: %v, name: %v", instanceID, pluginName)
	data, err := api.ListPlugins(instanceID, sleep, timeout)
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

// ListPlugins: list plugins from an instance.
func (api *API) ListPlugins(instanceID, sleep, timeout int) ([]map[string]interface{}, error) {
	return api.listPluginsWithRetry(instanceID, 1, sleep, timeout)
}

// listPluginsWithRetry: list plugins from an instance, with retry if backend is busy.
func (api *API) listPluginsWithRetry(instanceID, attempt, sleep, timeout int) (
	[]map[string]interface{}, error) {

	var (
		data   []map[string]interface{}
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/plugins", instanceID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("read plugins reached timeout of %d seconds", timeout)
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[INFO] go-api::plugins::read Timeout talking to backend "+
				"attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.listPluginsWithRetry(instanceID, attempt, sleep, timeout)
		}
		return nil, fmt.Errorf("ReadWithRetry failed, status: %v, message: %s", 400, failed)
	default:
		return nil,
			fmt.Errorf("list plugin with retry failed, status: %v, message: %s",
				response.StatusCode, failed)
	}
}

// UpdatePlugin: updates a plugin from an instance.
func (api *API) UpdatePlugin(instanceID int, pluginName string, enabled bool, sleep, timeout int) (
	map[string]interface{}, error) {

	var (
		failed map[string]interface{}
		params = make(map[string]interface{})
		path   = fmt.Sprintf("/api/instances/%d/plugins?async=true", instanceID)
	)

	params["plugin_name"] = pluginName
	params["enabled"] = enabled
	log.Printf("[DEBUG] go-api::plugin::update instance ID: %v, params: %v", instanceID, params)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilPluginChanged(instanceID, pluginName, enabled, 1, sleep, timeout)
	default:
		return nil,
			fmt.Errorf("update plugin failed, status: %v, message: %s", response.StatusCode, failed)
	}
}

// DisablePlugin: disables a plugin from an instance.
func (api *API) DisablePlugin(instanceID int, pluginName string, sleep, timeout int) (
	map[string]interface{}, error) {

	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, pluginName)
	)

	log.Printf("[DEBUG] go-api::plugin::disable path: %s", path)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilPluginChanged(instanceID, pluginName, false, 1, sleep, timeout)
	default:
		return nil, fmt.Errorf("disable plugin failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

// DeletePlugin: deletes a plugin from an instance.
func (api *API) DeletePlugin(instanceID int, pluginName string, sleep, timeout int) error {
	var (
		failed map[string]interface{}
		path   = fmt.Sprintf("/api/instances/%d/plugins/%s?async=true", instanceID, pluginName)
	)

	log.Printf("[DEBUG] go-api::plugin::delete path: %s", path)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		_, err = api.waitUntilPluginChanged(instanceID, pluginName, false, 1, sleep, timeout)
		return err
	default:
		return fmt.Errorf("delete plugin failed, status: %v, message: %s",
			response.StatusCode, failed)
	}
}

// waitUntilPluginChanged: wait until plugin changed.
func (api *API) waitUntilPluginChanged(instanceID int, pluginName string, enabled bool,
	attempt, sleep, timeout int) (map[string]interface{}, error) {

	for {
		if attempt*sleep > timeout {
			return nil, fmt.Errorf("wait until plugin changed reached timeout of %d seconds", timeout)
		}

		response, err := api.ReadPlugin(instanceID, pluginName, sleep, timeout)
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
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}
