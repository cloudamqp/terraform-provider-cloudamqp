package integrations

type LogRequest struct {
	AccessKeyID       string   `json:"access_key_id,omitempty"`
	APIKey            string   `json:"api_key,omitempty"`
	Application       string   `json:"application,omitempty"`
	ApplicationID     string   `json:"application_id,omitempty"`
	ApplicationSecret string   `json:"application_secret,omitempty"`
	ClientEmail       string   `json:"client_email,omitempty"`
	DCEURI            string   `json:"dce_uri,omitempty"`
	DCRID             string   `json:"dcr_id,omitempty"`
	Endpoint          string   `json:"endpoint,omitempty"`
	Host              string   `json:"host,omitempty"`
	HostPort          string   `json:"host_port,omitempty"`
	PrivateKey        string   `json:"private_key,omitempty"`
	PrivateKeyID      string   `json:"private_key_id,omitempty"`
	ProjectID         string   `json:"project_id,omitempty"`
	Region            string   `json:"region,omitempty"`
	SecretAccessKey   string   `json:"secret_access_key,omitempty"`
	Sourcetype        string   `json:"sourcetype,omitempty"`
	Subsystem         string   `json:"subsystem,omitempty"`
	Table             string   `json:"table,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	TenantID          string   `json:"tenant_id,omitempty"`
	Token             string   `json:"token,omitempty"`
	URL               string   `json:"url,omitempty"`
}

type LogResponse struct {
	ID     int64              `json:"id"`
	Type   string             `json:"type"`
	Config *LogConfigResponse `json:"config"`
}

type LogConfigResponse struct {
	AccessKeyID       *string   `json:"access_key_id,omitempty"`
	APIKey            *string   `json:"api_key,omitempty"`
	Application       *string   `json:"application,omitempty"`
	ApplicationID     *string   `json:"application_id,omitempty"`
	ApplicationSecret *string   `json:"application_secret,omitempty"`
	ClientEmail       *string   `json:"client_email,omitempty"`
	DCEURI            *string   `json:"dce_uri,omitempty"`
	DCRID             *string   `json:"dcr_id,omitempty"`
	Endpoint          *string   `json:"endpoint,omitempty"`
	Host              *string   `json:"host,omitempty"`
	HostPort          *string   `json:"host_port,omitempty"`
	PrivateKey        *string   `json:"private_key,omitempty"`
	PrivateKeyID      *string   `json:"private_key_id,omitempty"`
	ProjectID         *string   `json:"project_id,omitempty"`
	Region            *string   `json:"region,omitempty"`
	SecretAccessKey   *string   `json:"secret_access_key,omitempty"`
	Sourcetype        *string   `json:"sourcetype,omitempty"`
	Subsystem         *string   `json:"subsystem,omitempty"`
	Table             *string   `json:"table,omitempty"`
	Tags              *[]string `json:"tags,omitempty"`
	TenantID          *string   `json:"tenant_id,omitempty"`
	Token             *string   `json:"token,omitempty"`
	URL               *string   `json:"url,omitempty"`
}
