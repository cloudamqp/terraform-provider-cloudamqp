package cloudamqp

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceRabbitConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRabbitConfigurationRead,
		Update: resourceRabbitConfigurationUpdate,
		Delete: resourceRabbitConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"heartbeat": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     120,
				Description: "Set the server AMQP 0-9-1 heartbeat timeout in seconds.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be greater than 0, got: %d", key, v))
					}
					return
				},
			},
			"connection_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "Set the maximum permissible number of connection, -1 means infinity.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v == -1 {
						return
					}
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be -1 (infinity) or greater than 0, got: %d", key, v))
					}
					return
				},
			},
			"channel_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Set the maximum permissible number of channels per connection. 0 means unlimited",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be greater than or equal to 0, got: %d", key, v))
					}
					return
				},
			},
			"consumer_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     7200000,
				Description: "A consumer that has recevied a message and does not acknowledge that message within the timeout in milliseconds",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 10000 || v > 25000000 {
						errs = append(errs, fmt.Errorf("%q must be between 10000 and 25000000 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"vm_memory_high_watermark": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Default:     0.81,
				Description: "When the server will enter memory based flow-control as relative to the maximum available memory.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(float64)
					if v < 0.4 || v > 0.9 {
						errs = append(errs, fmt.Errorf("%q must be between 0.4 and 0.9 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"queue_index_embed_msgs_below": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     4096,
				Description: "Size in bytes below which to embed messages in the queue index.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 || v > 10485760 {
						errs = append(errs, fmt.Errorf("%q must be between 1 and 10485760 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"max_message_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     134217728,
				Description: "The largest allowed message payload size in bytes.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 || v > 536870912 {
						errs = append(errs, fmt.Errorf("%q must be between 1 and 536870912 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"log_exchange_level": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "error",
				Description: "Log level for the logger used for log integrations and the CloudAMQP Console log view. " +
					"Does not affect the file logger. Requires a RabbitMQ restart to be applied.",
				ValidateFunc: validateLogLevel(),
			},
		},
	}
}

func resourceRabbitConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())
	data, err := api.ReadRabbitConfiguration(instanceID)
	log.Printf("[DEBUG] cloudamqp::resource::rabbit_configuration::read data: %v", data)
	if err != nil {
		return err
	}
	d.Set("instance_id", instanceID)
	for k, v := range data {
		if validateRabbitConfigurationJSONField(k) {
			key := strings.ReplaceAll(k, "rabbit.", "")
			if key == "connection_max" {
				if v == "infinity" || v == nil {
					v = -1
				}
			} else if key == "queue_index_embed_msgs_below" {
				if v == nil {
					v = 4096
				}
			} else if key == "max_message_size" {
				if v == nil {
					v = 134217728
				}
			} else if key == "log.exchange.level" {
				key = "log_exchange_level"
			}
			d.Set(key, v)
		}
	}
	return nil
}

func resourceRabbitConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := rabbitConfigurationUpdateAttributeKeys()
	params := make(map[string]interface{})
	for _, k := range keys {
		v := d.Get(k)
		if k == "connection_max" {
			if v == -1 {
				v = "infinity"
			}
		} else if k == "log_exchange_level" {
			k = "log.exchange.level"
		}
		params["rabbit."+k] = v
	}

	err := api.UpdateRabbitConfiguration(d.Get("instance_id").(int), params)
	if err != nil {
		return err
	}
	return resourceRabbitConfigurationRead(d, meta)
}

func resourceRabbitConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func validateRabbitConfigurationJSONField(key string) bool {
	switch key {
	case "rabbit.heartbeat",
		"rabbit.connection_max",
		"rabbit.channel_max",
		"rabbit.consumer_timeout",
		"rabbit.vm_memory_high_watermark",
		"rabbit.queue_index_embed_msgs_below",
		"rabbit.max_message_size",
		"rabbit.log.exchange.level":
		return true
	}
	return false
}

func validateLogLevel() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"debug",
		"info",
		"warning",
		"error",
		"critical",
		"critical",
		"none",
	}, true)
}

func rabbitConfigurationUpdateAttributeKeys() []string {
	return []string{
		"heartbeat",
		"connection_max",
		"channel_max",
		"consumer_timeout",
		"vm_memory_high_watermark",
		"queue_index_embed_msgs_below",
		"max_message_size",
		"log_exchange_level",
	}
}
