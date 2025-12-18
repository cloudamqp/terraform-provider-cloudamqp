package configuration

type TrustStoreRequest struct {
	Url             string `json:"url"`
	RefreshInterval int64  `json:"refresh_interval"`
	Provider        string `json:"provider"`
}

type TrustStoreResponse struct {
	ConfigurationId string `json:"id"`
	Url             string `json:"url"`
	RefreshInterval int64  `json:"refresh_interval"`
	Provider        string `json:"provider"`
}
