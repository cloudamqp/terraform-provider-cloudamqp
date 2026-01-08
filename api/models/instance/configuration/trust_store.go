package configuration

type TrustStoreRequest struct {
	Provider        string    `json:"provider"`
	RefreshInterval int64     `json:"refresh_interval"`
	Url             string    `json:"url,omitempty"`
	CACert          string    `json:"cacertfile,omitempty"`
	Certificates    *[]string `json:"certificates,omitempty"`
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
	if sanitized.Certificates != nil {
		sanitized.Certificates = &[]string{"***"}
	}
	return sanitized
}
