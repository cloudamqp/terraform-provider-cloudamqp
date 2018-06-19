package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dghubble/sling"
)

type API struct {
	sling *sling.Sling
}

func New(baseUrl, apiKey string) *API {
	return &API{
		sling: sling.New().
			Client(http.DefaultClient).
			Base(baseUrl).
			SetBasicAuth("", apiKey),
	}
}

func (api *API) waitUntilReady(id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	for {
		_, err := api.sling.Path("/api/instances/").Get(id).ReceiveSuccess(&data)
		if err != nil {
			return nil, err
		}
		if data["ready"] == true {
			data["id"] = id
			return data, nil
		}
		time.Sleep(10 * time.Second)
	}
}

func (api *API) Create(params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	_, err := api.sling.Post("/api/instances").BodyJSON(params).ReceiveSuccess(&data)
	if err != nil {
		return nil, err
	}
	string_id := strconv.Itoa(int(data["id"].(float64)))
	return api.waitUntilReady(string_id)
}

func (api *API) Read(id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	_, err := api.sling.Path("/api/instances/").Get(id).ReceiveSuccess(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (api *API) Update(id string, params map[string]interface{}) error {
	_, err := api.sling.Put("/api/instances/" + id).BodyJSON(params).ReceiveSuccess(nil)
	return err
}

func (api *API) Delete(id string) error {
	_, err := api.sling.Path("/api/instances/").Delete(id).ReceiveSuccess(nil)
	return err
}
