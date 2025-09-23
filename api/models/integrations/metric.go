package integrations

type MetricRequest struct {
	AccessKeyID     string   `json:"access_key_id,omitempty"`
	APIKey          string   `json:"api_key,omitempty"`
	ClientEmail     string   `json:"client_email,omitempty"`
	Email           string   `json:"email,omitempty"`
	IAMExternalID   string   `json:"iam_external_id,omitempty"`
	IAMRole         string   `json:"iam_role,omitempty"`
	IncludeAdQueues bool     `json:"include_ad_queues,omitempty"`
	HostPort        string   `json:"host_port,omitempty"`
	PrivateKey      string   `json:"private_key,omitempty"`
	PrivateKeyID    string   `json:"private_key_id,omitempty"`
	ProjectID       string   `json:"project_id,omitempty"`
	QueueRegex      string   `json:"queue_regex,omitempty"`
	Region          string   `json:"region,omitempty"`
	SecretAccessKey string   `json:"secret_access_key,omitempty"`
	Sourcetype      string   `json:"sourcetype,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Token           string   `json:"token,omitempty"`
	VhostRegex      string   `json:"vhost_regex,omitempty"`
}

type MetricResponse struct {
	ID     int64                 `json:"id"`
	Type   string                `json:"type"`
	Config *MetricConfigResponse `json:"config"`
}

type MetricConfigResponse struct {
	AccessKeyID     *string   `json:"access_key_id,omitempty"`
	APIKey          *string   `json:"api_key,omitempty"`
	ClientEmail     *string   `json:"client_email,omitempty"`
	Email           *string   `json:"email,omitempty"`
	IAMExternalID   *string   `json:"iam_external_id,omitempty"`
	IAMRole         *string   `json:"iam_role,omitempty"`
	IncludeAdQueues *bool     `json:"include_ad_queues,omitempty"`
	HostPort        *string   `json:"host_port,omitempty"`
	PrivateKey      *string   `json:"private_key,omitempty"`
	PrivateKeyID    *string   `json:"private_key_id,omitempty"`
	ProjectID       *string   `json:"project_id,omitempty"`
	QueueRegex      *string   `json:"queue_regex,omitempty"`
	Region          *string   `json:"region,omitempty"`
	SecretAccessKey *string   `json:"secret_access_key,omitempty"`
	Sourcetype      *string   `json:"sourcetype,omitempty"`
	Tags            *[]string `json:"tags,omitempty"`
	Token           *string   `json:"token,omitempty"`
	VhostRegex      *string   `json:"vhost_regex,omitempty"`
}
