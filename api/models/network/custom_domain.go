package network

type CustomDomainRequest struct {
	Hostname string `json:"hostname"`
}

type CustomDomainResponse struct {
	Hostname   string `json:"hostname"`
	Configured bool   `json:"configured"`
}
