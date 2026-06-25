package integrations

type LogAgentRequest struct {
	// CloudWatch
	Region        string `json:"region,omitempty"`
	IAMRole       string `json:"iam_role,omitempty"`
	IAMExternalID string `json:"iam_external_id,omitempty"`
	LogGroupName  string `json:"log_group_name,omitempty"`
	LogStreamName string `json:"log_stream_name,omitempty"`
	// Uptrace
	DSN string `json:"dsn,omitempty"`
	// Splunk
	Endpoint   string `json:"endpoint,omitempty"`
	Token      string `json:"token,omitempty"`
	SourceType string `json:"sourcetype,omitempty"`
	// Coralogix
	PrivateKey  string `json:"private_key,omitempty"`
	Application string `json:"application,omitempty"`
	Subsystem   string `json:"subsystem,omitempty"`
	// Datadog
	APIKey string `json:"api_key,omitempty"`
	Tags   string `json:"tags,omitempty"`
	// Custom OTLP
	AuthType string            `json:"auth_type,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
}

type LogAgentResponse struct {
	ID     int64                  `json:"id"`
	Type   string                 `json:"type"`
	Config LogAgentConfigResponse `json:"config"`
}

type LogAgentConfigResponse struct {
	// CloudWatch
	Region        *string `json:"region,omitempty"`
	IAMRole       *string `json:"iam_role,omitempty"`
	IAMExternalID *string `json:"iam_external_id,omitempty"`
	LogGroupName  *string `json:"log_group_name,omitempty"`
	LogStreamName *string `json:"log_stream_name,omitempty"`
	// Uptrace
	DSN *string `json:"dsn,omitempty"`
	// Splunk
	Endpoint   *string `json:"endpoint,omitempty"`
	Token      *string `json:"token,omitempty"`
	SourceType *string `json:"sourcetype,omitempty"`
	// Coralogix
	PrivateKey  *string `json:"private_key,omitempty"`
	Application *string `json:"application,omitempty"`
	Subsystem   *string `json:"subsystem,omitempty"`
	// Datadog
	APIKey *string `json:"api_key,omitempty"`
	Tags   *string `json:"tags,omitempty"`
	// Custom OTLP
	AuthType *string           `json:"auth_type,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
	Username *string           `json:"username,omitempty"`
}
