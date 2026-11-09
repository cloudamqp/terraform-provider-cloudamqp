package monitoring

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type AlarmRequest struct {
	Type             string   `json:"type"`
	Enabled          bool     `json:"enabled"`
	ReminderInterval *int64   `json:"reminder_interval,omitempty"`
	ValueThreshold   *int64   `json:"value_threshold,omitempty"`
	ValueCalculation string   `json:"value_calculation,omitempty"`
	TimeThreshold    *int64   `json:"time_threshold,omitempty"`
	AllowDowntime    *bool    `json:"allow_downtime,omitempty"`
	MaxDiskSize      *int64   `json:"max_disk_size,omitempty"`
	VhostRegex       string   `json:"vhost_regex,omitempty"`
	QueueRegex       string   `json:"queue_regex,omitempty"`
	MessageType      string   `json:"message_type,omitempty"`
	Recipients       *[]int64 `json:"recipients,omitempty"`
}

type AlarmResponse struct {
	ID               int64       `json:"id"`
	Type             string      `json:"type"`
	Enabled          bool        `json:"enabled"`
	ReminderInterval *int64      `json:"reminder_interval,omitempty"`
	ValueThreshold   *int64      `json:"value_threshold,omitempty"`
	ValueCalculation *string     `json:"value_calculation,omitempty"`
	TimeThreshold    *int64      `json:"time_threshold,omitempty"`
	AllowDowntime    *bool       `json:"allow_downtime,omitempty"`
	MaxDiskSize      IntOrString `json:"max_disk_size"`
	VhostRegex       *string     `json:"vhost_regex,omitempty"`
	QueueRegex       *string     `json:"queue_regex,omitempty"`
	MessageType      *string     `json:"message_type,omitempty"`
	Recipients       *[]int64    `json:"recipients,omitempty"`
}

type IntOrString struct {
	Value int64
	Valid bool
}

func (i *IntOrString) UnmarshalJSON(data []byte) error {
	*i = IntOrString{}
	if string(data) == "null" {
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		if asString == "" {
			return nil
		}
		value, err := strconv.ParseInt(asString, 10, 64)
		if err != nil {
			return fmt.Errorf("IntOrString: invalid value %q", asString)
		}
		*i = IntOrString{Value: value, Valid: true}
		return nil
	}

	var asInt int64
	if err := json.Unmarshal(data, &asInt); err != nil {
		return fmt.Errorf("IntOrString: invalid JSON value %s", data)
	}
	*i = IntOrString{Value: asInt, Valid: true}
	return nil
}

func (i IntOrString) Int64Pointer() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Value
}
