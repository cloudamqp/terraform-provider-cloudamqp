package api

import (
	"net/http"
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
