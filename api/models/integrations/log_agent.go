package integrations

type LogAgentRequest struct {
	// CloudWatch
	Region        string `json:"region,omitempty"`
	IAMRole       string `json:"iam_role,omitempty"`
	IAMExternalID string `json:"iam_external_id,omitempty"`
	LogGroup      string `json:"log_group,omitempty"`
	LogStream     string `json:"log_stream,omitempty"`
	// Uptrace
	DSN string `json:"dsn,omitempty"`
	// Splunk
	Endpoint   string `json:"endpoint,omitempty"`
	Token      string `json:"token,omitempty"`
	SourceType string `json:"sourcetype,omitempty"`
	// Coralogix
	Domain      string `json:"domain,omitempty"`
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
	// Google Cloud
	CredentialType string `json:"type,omitempty"`
	ProjectID      string `json:"project_id,omitempty"`
	ClientEmail    string `json:"client_email,omitempty"`
	PrivateKeyID   string `json:"private_key_id,omitempty"`
	// Grafana
	GrafanaInstanceID string `json:"instance_id,omitempty"`
	APIToken          string `json:"api_token,omitempty"`
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
	LogGroup      *string `json:"log_group,omitempty"`
	LogStream     *string `json:"log_stream,omitempty"`
	// Uptrace
	DSN *string `json:"dsn,omitempty"`
	// Splunk
	Endpoint   *string `json:"endpoint,omitempty"`
	Token      *string `json:"token,omitempty"`
	SourceType *string `json:"sourcetype,omitempty"`
	// Coralogix
	Domain      *string `json:"domain,omitempty"`
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
	// Google Cloud
	ProjectID    *string `json:"project_id,omitempty"`
	ClientEmail  *string `json:"client_email,omitempty"`
	PrivateKeyID *string `json:"private_key_id,omitempty"`
	// Grafana
	GrafanaInstanceID *string `json:"instance_id,omitempty"`
}

func redactedString(s string) string {
	if s != "" {
		return "***"
	}
	return ""
}

func redactedStringPtr(s *string) *string {
	if s != nil && *s != "" {
		v := "***"
		return &v
	}
	return s
}

func (r LogAgentRequest) Sanitized() LogAgentRequest {
	sanitized := r
	sanitized.DSN = redactedString(r.DSN)
	sanitized.Token = redactedString(r.Token)
	sanitized.PrivateKey = redactedString(r.PrivateKey)
	sanitized.PrivateKeyID = redactedString(r.PrivateKeyID)
	sanitized.APIKey = redactedString(r.APIKey)
	sanitized.Password = redactedString(r.Password)
	sanitized.APIToken = redactedString(r.APIToken)
	if len(r.Headers) > 0 {
		sanitized.Headers = make(map[string]string, len(r.Headers))
		for k := range r.Headers {
			sanitized.Headers[k] = "***"
		}
	}
	return sanitized
}

func (r LogAgentResponse) Sanitized() LogAgentResponse {
	sanitized := r
	sanitized.Config = r.Config.Sanitized()
	return sanitized
}

func (c LogAgentConfigResponse) Sanitized() LogAgentConfigResponse {
	sanitized := c
	sanitized.DSN = redactedStringPtr(c.DSN)
	sanitized.Token = redactedStringPtr(c.Token)
	sanitized.PrivateKey = redactedStringPtr(c.PrivateKey)
	sanitized.APIKey = redactedStringPtr(c.APIKey)
	if len(c.Headers) > 0 {
		sanitized.Headers = make(map[string]string, len(c.Headers))
		for k := range c.Headers {
			sanitized.Headers[k] = "***"
		}
	}
	return sanitized
}
