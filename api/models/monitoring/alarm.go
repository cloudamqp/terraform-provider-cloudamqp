package monitoring

type AlarmRequest struct {
	Type             string   `json:"type"`
	Enabled          bool     `json:"enabled"`
	ReminderInterval *int64   `json:"reminder_interval,omitempty"`
	ValueThreshold   *int64   `json:"value_threshold,omitempty"`
	ValueCalculation string   `json:"value_calculation,omitempty"`
	TimeThreshold    *int64   `json:"time_threshold,omitempty"`
	VhostRegex       string   `json:"vhost_regex,omitempty"`
	QueueRegex       string   `json:"queue_regex,omitempty"`
	MessageType      string   `json:"message_type,omitempty"`
	Recipients       *[]int64 `json:"recipients,omitempty"`
}

type AlarmResponse struct {
	ID               int64    `json:"id"`
	Type             string   `json:"type"`
	Enabled          bool     `json:"enabled"`
	ReminderInterval *int64   `json:"reminder_interval,omitempty"`
	ValueThreshold   *int64   `json:"value_threshold,omitempty"`
	ValueCalculation *string  `json:"value_calculation,omitempty"`
	TimeThreshold    *int64   `json:"time_threshold,omitempty"`
	VhostRegex       *string  `json:"vhost_regex,omitempty"`
	QueueRegex       *string  `json:"queue_regex,omitempty"`
	MessageType      *string  `json:"message_type,omitempty"`
	Recipients       *[]int64 `json:"recipients,omitempty"`
}
