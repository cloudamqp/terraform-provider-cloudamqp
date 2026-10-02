package instance

type PluginRequest struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type PluginResponse struct {
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	Require     string  `json:"require"`
	Description string  `json:"description"`
	Enabled     bool    `json:"enabled"`
	Recommended *string `json:"recommended,omitempty"`
	Required    *string `json:"required,omitempty"`
}
