package cloudamqp

import (
	"net/http"
	"net/url"

	"path"

	"github.com/dghubble/sling"
)

const (
	defaultBaseURL = "https://customer.cloudamqp.com/api/"
)

// Client is a sling client for instance resources
type Client struct {
	sling     *sling.Sling
	Instances *InstanceService
}

// NewClient returns a new CloudAMQP API client.
// If a nil httpClient is given, the http.DefaultClient will be used.
// If a nil baseURL is given, the defaultBaseURL will be used.
func NewClient(httpClient *http.Client, baseURL *url.URL, token string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if baseURL == nil {
		baseURL, _ = url.Parse(defaultBaseURL)
	}
	baseURL.Path = path.Join(baseURL.Path) + "/"

	base := sling.New().Base(baseURL.String()).Client(httpClient)

	if token != "" {
		base.SetBasicAuth(token, "")
	}

	c := &Client{
		sling:     base,
		Instances: newInstanceService(base.New()),
	}
	return c
}
