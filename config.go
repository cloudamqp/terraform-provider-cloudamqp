package main

import (
	"log"
	"net/url"

	"github.com/waveaccounting/go-cloudamqp/cloudamqp"
)

// Config is the configuration structure used to instantiate the Sentry
// provider.
type Config struct {
	APIKey  string
	BaseURL string
}

func (c *Config) Client() (interface{}, error) {
	var baseURL *url.URL
	var err error

	if c.BaseURL != "" {
		baseURL, err = url.Parse(c.BaseURL)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("[INFO] Instantiating CloudAMQP client...")
	cl := cloudamqp.NewClient(nil, baseURL, c.APIKey)

	return cl, nil
}
