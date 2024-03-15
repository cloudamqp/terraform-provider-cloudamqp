package api

import (
	"net/http"

	"github.com/dghubble/sling"
)

type API struct {
	sling  *sling.Sling
	client *http.Client
}

func (api *API) DefaultRmqVersion() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	_, err := api.sling.New().Get("/api/default_rabbitmq_version").Receive(&data, &failed)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func New(baseUrl, apiKey string, useragent string, client *http.Client) *API {
	if len(useragent) == 0 {
		useragent = "84codes go-api"
	}
	return &API{
		sling: sling.New().
			Client(client).
			Base(baseUrl).
			SetBasicAuth("", apiKey).
			Set("User-Agent", useragent),
		client: client,
	}
}
