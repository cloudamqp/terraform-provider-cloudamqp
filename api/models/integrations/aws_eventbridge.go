package integrations

type AwsEventBridgeRequest struct {
	AwsAccountId string `json:"aws_account_id"`
	AwsRegion    string `json:"aws_region"`
	Vhost        string `json:"vhost"`
	QueueName    string `json:"queue"`
	WithHeaders  bool   `json:"with_headers"`
	Prefetch     *int64 `json:"prefetch,omitempty"`
}

type AwsEventBridgeResponse struct {
	Id           int64   `json:"id"`
	AwsAccountId string  `json:"aws_account_id"`
	AwsRegion    string  `json:"aws_region"`
	Vhost        string  `json:"vhost"`
	QueueName    string  `json:"queue"`
	WithHeaders  bool    `json:"with_headers"`
	Prefetch     int64   `json:"prefetch"`
	Status       *string `json:"status,omitempty"`
}
