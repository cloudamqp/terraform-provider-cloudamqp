package integrations

type MetricsRequest struct {
	AccessKeyID     string   `json:"access_key_id,omitempty"`
	APIKey          string   `json:"api_key,omitempty"`
	ClientEmail     string   `json:"client_email,omitempty"`
	Email           string   `json:"email,omitempty"`
	IAMExternalID   string   `json:"iam_external_id,omitempty"`
	IAMRole         string   `json:"iam_role,omitempty"`
	IncludeAdQueues bool     `json:"include_ad_queues,omitempty"`
	LicenseKey      string   `json:"license_key,omitempty"`
	PrivateKey      string   `json:"private_key,omitempty"`
	PrivateKeyID    string   `json:"private_key_id,omitempty"`
	ProjectID       string   `json:"project_id,omitempty"`
	QueueRegex      string   `json:"queue_regex,omitempty"`
	Region          string   `json:"region,omitempty"`
	SecretAccessKey string   `json:"secret_access_key,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	VhostRegex      string   `json:"vhost_regex,omitempty"`
}

type MetricsResponse struct {
	ID     int64                  `json:"id"`
	Type   string                 `json:"type"`
	Config *MetricsConfigResponse `json:"config"`
}

type MetricsConfigResponse struct {
	AccessKeyID     *string   `json:"access_key_id,omitempty"`
	APIKey          *string   `json:"api_key,omitempty"`
	ClientEmail     *string   `json:"client_email,omitempty"`
	Email           *string   `json:"email,omitempty"`
	IAMExternalID   *string   `json:"iam_external_id,omitempty"`
	IAMRole         *string   `json:"iam_role,omitempty"`
	IncludeAdQueues *bool     `json:"include_ad_queues,omitempty"`
	LicenseKey      *string   `json:"license_key,omitempty"`
	PrivateKey      *string   `json:"private_key,omitempty"`
	PrivateKeyID    *string   `json:"private_key_id,omitempty"`
	ProjectID       *string   `json:"project_id,omitempty"`
	QueueRegex      *string   `json:"queue_regex,omitempty"`
	Region          *string   `json:"region,omitempty"`
	Tags            *[]string `json:"tags,omitempty"`
	SecretAccessKey *string   `json:"secret_access_key,omitempty"`
	VhostRegex      *string   `json:"vhost_regex,omitempty"`
}
