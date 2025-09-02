package integrations

type WebhookCreateRequest struct {
	Concurrency int64  `json:"concurrency"`
	WebhookURI  string `json:"webhook_uri"`
	Vhost       string `json:"vhost"`
	Queue       string `json:"queue"`
}

type WebhookUpdateRequest struct {
	WebhookID   int64  `json:"webhook_id"`
	Concurrency int64  `json:"concurrency"`
	WebhookURI  string `json:"webhook_uri"`
	Vhost       string `json:"vhost"`
	Queue       string `json:"queue"`
}

type WebhookResponse struct {
	ID          int64  `json:"id"`
	Concurrency int64  `json:"concurrency"`
	WebhookURI  string `json:"webhook_uri"`
	Vhost       string `json:"vhost"`
	Queue       string `json:"queue"`
	LastStatus  string `json:"last_status,omitempty"`
}
