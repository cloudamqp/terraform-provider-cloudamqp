package integrations

type LogAgentRequest struct {
	// Type string `json:"type"`
	// CloudWatch
	Region        string `json:"region,omitempty"`
	IAMRole       string `json:"iam_role,omitempty"`
	IAMExternalID string `json:"iam_external_id,omitempty"`
	LogGroupName  string `json:"log_group_name,omitempty"`
	LogStreamName string `json:"log_stream_name,omitempty"`
	// Config LogAgentConfigRequest `json:"config"`
}

type LogAgentConfigRequest struct {
	// CloudWatch
	Region        string `json:"region,omitempty"`
	IAMRole       string `json:"iam_role,omitempty"`
	IAMExternalID string `json:"iam_external_id,omitempty"`
	LogGroupName  string `json:"log_group_name,omitempty"`
	LogStreamName string `json:"log_stream_name,omitempty"`
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
}
