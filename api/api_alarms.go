package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dghubble/sling"
)

type AlarmAPI struct {
	sling *sling.Sling
}

func New(baseUrl, apiKey string) *AlarmAPI {
	return &AlarmAPI{
		sling: sling.New().
			Client(http.DefaultClient).
			Base(baseUrl).
			SetBasicAuth("", apiKey),
	}
}

func (api *AlarmAPI) Create(params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	_, err := api.sling.Post("/api/alarms").BodyJSON(params).ReceiveSuccess(&data)
	if err != nil {
		return nil, err
	}
	string_id := strconv.Itoa(int(data["id"].(float64)))
	return string_id//api.waitUntilReady(string_id)
}

func (api *AlarmAPI) Read(id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	_, err := api.sling.Path("/api/alarms/").Get(id).ReceiveSuccess(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// alarm_id, type
func (api *AlarmAPI) Delete(id string, alarm_type string) error {
	_, err := api.sling.Path("/api/alarms/").Delete(id).ReceiveSuccess(nil)
	return err
}
