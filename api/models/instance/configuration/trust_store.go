package configuration

type TrustStoreRequest struct {
	RefreshInterval int64    `json:"refresh_interval"`
	Url             string   `json:"url"`
	CACert          string   `json:"cacertfile,omitempty"`
	Certs           []string `json:"certs,omitempty"`
}

type TrustStoreResponse struct {
	ConfigurationId string  `json:"id"`
	Url             *string `json:"url,omitempty"`
	RefreshInterval int64   `json:"refresh_interval"`
	Provider        string  `json:"provider"`
}

func (u TrustStoreRequest) Sanitized() TrustStoreRequest {
	sanitized := u
	if sanitized.CACert != "" {
		sanitized.CACert = "***"
	}
	if len(sanitized.Certs) > 0 {
		sanitized.Certs = []string{"***"}
	}
	return sanitized
}
