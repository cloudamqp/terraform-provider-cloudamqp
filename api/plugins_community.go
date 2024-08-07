package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// InstallPluginCommunity: install a community plugin on an instance.
func (api *API) InstallPluginCommunity(instanceID int, pluginName string, sleep, timeout int) (
	map[string]any, error) {

	var (
		failed map[string]any
		params = make(map[string]any)
		path   = fmt.Sprintf("/api/instances/%d/plugins/community?async=true", instanceID)
	)

	params["plugin_name"] = pluginName
	log.Printf("[DEBUG] api::plugin_community#enable path: %s", path)
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilPluginChanged(instanceID, pluginName, true, 1, sleep, timeout)
	default:
		return nil, fmt.Errorf("install community plugin failed, status: %v, message: %v",
			response.StatusCode, failed)
	}
}

// ReadPluginCommunity: reads a specific community plugin from an instance.
func (api *API) ReadPluginCommunity(instanceID int, pluginName string, sleep, timeout int) (
	map[string]any, error) {

	log.Printf("[DEBUG] api::plugin_community#read instance ID: %v, name: %v", instanceID, pluginName)
	data, err := api.ListPluginsCommunity(instanceID, sleep, timeout)
	if err != nil {
		return nil, err
	}

	for _, plugin := range data {
		if plugin["name"] == pluginName {
			log.Printf("[DEBUG] api::plugin_community#read found plugin: %v", pluginName)
			return plugin, nil
		}
	}

	return nil, nil
}

// ListPluginsCommunity: list all community plugins for an instance.
func (api *API) ListPluginsCommunity(instanceID, sleep, timeout int) ([]map[string]any, error) {
	return api.listPluginsCommunityWithRetry(instanceID, 1, sleep, timeout)
}

// listPluginsCommunityWithRetry: list all community plugins for an instance,
// with retry if the backend is busy.
func (api *API) listPluginsCommunityWithRetry(instanceID, attempt, sleep, timeout int) (
	[]map[string]any, error) {

	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/community", instanceID)
	)

	log.Printf("[DEBUG] api::plugin_community#listPluginsCommunityWithRetry path: %s", path)
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("list community plugins reached timeout of %d seconds", timeout)
	}

	log.Printf("[DEBUG] api::plugin_community#listPluginsCommunityWithRetry statusCode: %d",
		response.StatusCode)
	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			log.Printf("[INFO] api::plugins-community#read Timeout talking to backend "+
				"attempt: %d, until timeout: %d", attempt, (timeout - (attempt * sleep)))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.listPluginsCommunityWithRetry(instanceID, attempt, sleep, timeout)
		}
	}
	return nil, fmt.Errorf("list community plugins failed, status: %d, message: %v",
		response.StatusCode, failed)
}

// UpdatePluginCommunity: updates a community plugin from an instance.
func (api *API) UpdatePluginCommunity(instanceID int, pluginName string, enabled bool,
	sleep, timeout int) (map[string]any, error) {

	var (
		failed map[string]any
		params = make(map[string]any)
		path   = fmt.Sprintf("/api/instances/%d/plugins/community?async=true", instanceID)
	)

	params["plugin_name"] = pluginName
	params["enabled"] = enabled
	log.Printf("[DEBUG] api::plugin_community#update path: %s", path)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilPluginChanged(instanceID, pluginName, enabled, 1, sleep, timeout)
	default:
		return nil, fmt.Errorf("update community plugin failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

// UninstallPluginCommunity: uninstall a community plugin from an instance.
func (api *API) UninstallPluginCommunity(instanceID int, pluginName string, sleep, timeout int) (
	map[string]any, error) {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/plugins/community/%s?async=true", instanceID, pluginName)
	)

	log.Printf("[DEBUG] api::plugin_community#disable path: %s", path)
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilPluginUninstalled(instanceID, pluginName, 1, sleep, timeout)
	default:
		return nil, fmt.Errorf("disable community plugin failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

// waitUntilPluginUninstalled: wait until a community plugin been uninstalled.
func (api *API) waitUntilPluginUninstalled(instanceID int, pluginName string,
	attempt, sleep, timeout int) (map[string]any, error) {

	log.Printf("[DEBUG] api::plugin_community#waitUntilPluginUninstalled instance id: %v, name: %v",
		instanceID, pluginName)
	for {
		if attempt*sleep > timeout {
			return nil, fmt.Errorf("wait until plugin uninstalled reached timeout of %d seconds", timeout)
		}

		response, err := api.ReadPlugin(instanceID, pluginName, sleep, timeout)
		if err != nil {
			return nil, err
		}
		if len(response) == 0 {
			return response, nil
		}
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}
