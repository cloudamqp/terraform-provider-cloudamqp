package api

import (
	"fmt"
	"errors"
)

type PluginParams struct {
	Name    string `url:"plugin_name,omitempty"`
	Enable    bool `url:"enable,omitempty"`
}

func (api *API) EnablePlugin(instance_id int, name string) error {
	failed := make(map[string]interface{})
	params := &PluginParams{Name: name}
	path := fmt.Sprintf("/api/instances/%d/plugins", instance_id)
	response, err := api.sling.Post(path).BodyForm(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("EnablePlugin failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return nil
}

func (api *API) ReadPlugin(instance_id int, plugin_name string) (map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/plugins", instance_id)
	response, err := api.sling.Get(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadPlugin failed, status: %v, message: %s", response.StatusCode, failed))
	}

	for index, value := range data {
		if value["name"] == plugin_name {
			return data[index], nil
		}
	}

	return nil, nil
}

func (api *API) ReadPlugins(instance_id int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/plugins", instance_id)
	response, err := api.sling.Get(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadPlugin failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) UpdatePlugin(instance_id int, name string, enable bool) error {
	failed := make(map[string]interface{})
	params := &PluginParams{Name: name, Enable: enable}
	path := fmt.Sprintf("/api/instances/%d/plugins", instance_id)
	response, err := api.sling.Put(path).BodyForm(params).Receive(nil, &failed)

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("UpdatePlugin failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}

func (api *API) DisablePlugin(instance_id int, name string) error {
	failed := make(map[string]interface{})
	params := &PluginParams{Name: name}
	path := fmt.Sprintf("/api/instances/%d/plugins", instance_id)
	response, err := api.sling.Delete(path).BodyForm(params).Receive(nil, &failed)

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("DisablePlugin failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return err
}
