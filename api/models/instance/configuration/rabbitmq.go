package configuration

import (
	"encoding/json"
	"fmt"
)

type RabbitMqConfigRequest struct {
	Heartbeat                *int64                `json:"rabbit.heartbeat,omitempty"`
	ConnectionMax            *ConnectionMaxValue   `json:"rabbit.connection_max,omitempty"`
	ChannelMax               *int64                `json:"rabbit.channel_max,omitempty"`
	ConsumerTimeout          *ConsumerTimeoutValue `json:"rabbit.consumer_timeout,omitempty"`
	VmMemoryHighWatermark    *float64              `json:"rabbit.vm_memory_high_watermark,omitempty"`
	QueueIndexEmbedMsgsBelow *int64                `json:"rabbit.queue_index_embed_msgs_below,omitempty"`
	MaxMessageSize           *int64                `json:"rabbit.max_message_size,omitempty"`
	LogExchangeLevel         *string               `json:"rabbit.log.exchange.level,omitempty"`
	ClusterPartitionHandling *string               `json:"rabbit.cluster_partition_handling,omitempty"`
}

type RabbitMqConfigResponse struct {
	Heartbeat                int64                `json:"rabbit.heartbeat"`
	ConnectionMax            ConnectionMaxValue   `json:"rabbit.connection_max"`
	ChannelMax               int64                `json:"rabbit.channel_max"`
	ConsumerTimeout          ConsumerTimeoutValue `json:"rabbit.consumer_timeout"`
	VmMemoryHighWatermark    float64              `json:"rabbit.vm_memory_high_watermark"`
	QueueIndexEmbedMsgsBelow *int64               `json:"rabbit.queue_index_embed_msgs_below,omitempty"`
	MaxMessageSize           int64                `json:"rabbit.max_message_size"`
	LogExchangeLevel         string               `json:"rabbit.log.exchange.level"`
	ClusterPartitionHandling string               `json:"rabbit.cluster_partition_handling"`
}

// Custom type for ConnectionMax
type ConnectionMaxValue struct {
	IsInfinity bool
	Value      int64
}

func (c ConnectionMaxValue) MarshalJSON() ([]byte, error) {
	if c.IsInfinity {
		return json.Marshal("infinity")
	}
	return json.Marshal(c.Value)
}

func (c *ConnectionMaxValue) UnmarshalJSON(data []byte) error {
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		if asString == "infinity" {
			c.IsInfinity = true
			c.Value = 0
			return nil
		}
	}
	var asInt int64
	if err := json.Unmarshal(data, &asInt); err == nil {
		c.IsInfinity = false
		c.Value = asInt
		return nil
	}
	return fmt.Errorf("ConnectionMaxValue: invalid JSON value")
}

// Custom type for ConsumerTimeout
type ConsumerTimeoutValue struct {
	IsEnabled bool
	Value     int64
}

func (c ConsumerTimeoutValue) MarshalJSON() ([]byte, error) {
	if c.IsEnabled {
		return json.Marshal(c.Value)
	}
	return json.Marshal(false)
}

func (c *ConsumerTimeoutValue) UnmarshalJSON(data []byte) error {
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		if asString == "false" {
			c.IsEnabled = false
			c.Value = -1
			return nil
		}
	}
	var asInt int64
	if err := json.Unmarshal(data, &asInt); err == nil {
		c.IsEnabled = true
		c.Value = asInt
		return nil
	}
	return fmt.Errorf("ConsumerTimeoutValue: invalid JSON value")
}
