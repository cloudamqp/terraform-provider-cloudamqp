package api

import (
	"strconv"
	"time"
	"fmt"
	"errors"
)

func (api *API) waitUntilReady(id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		response, err := api.sling.Path("/api/instances/").Get(id).Receive(&data, &failed)

		if err != nil {
			return nil, err
		}
		if response.StatusCode != 200 {
			return nil, errors.New(fmt.Sprintf("waitUntilReady failed, status: %v, message: %s", response.StatusCode, failed))
		}
		if data["ready"] == true {
			data["id"] = id
			return data, nil
		}

		time.Sleep(10 * time.Second)
	}
}

func (api *API) CreateInstance(params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	response, err := api.sling.Post("/api/instances").BodyJSON(params).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 201 {
		return nil, errors.New(fmt.Sprintf("CreateInstance failed, status: %v, message: %s", response.StatusCode, failed))
	}

	string_id := strconv.Itoa(int(data["id"].(int)))
	return api.waitUntilReady(string_id)
}

func (api *API) ReadInstance(id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	response, err := api.sling.Path("/api/instances/").Get(id).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadInstance failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) ReadInstances() ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	response, err := api.sling.Get("/api/instances").Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadInstances failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) UpdateInstance(id string, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	response, err := api.sling.Put("/api/instances/" + id).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("UpdateInstance failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return nil
}

func (api *API) DeleteInstance(id string) error {
	failed := make(map[string]interface{})
	response, err := api.sling.Path("/api/instances/").Delete(id).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("DeleteInstance failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return nil
}
