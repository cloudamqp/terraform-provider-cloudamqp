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

func resourceRabbitMqConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceRabbitMqConfigurationCreate,
		Read:   resourceRabbitMqConfigurationRead,
		Update: resourceRabbitMqConfigurationUpdate,
		Delete: resourceRabbitMqConfigurationDelete,
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
				Computed:    true,
				Optional:    true,
				Description: "Set the server AMQP 0-9-1 heartbeat timeout in seconds.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be greater than or equal to 0, got: %d", key, v))
					}
					return
				},
			},
			"connection_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Set the maximum permissible number of connection, -1 means infinity.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v == -1 {
						return
					}
					if v < 1 {
						errs = append(errs, fmt.Errorf("%q must be -1 (infinity) or greater than 0, got: %d", key, v))
					}
					return
				},
			},
			"channel_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
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
				Computed:    true,
				Optional:    true,
				Description: "A consumer that has recevied a message and does not acknowledge that message within the timeout in milliseconds",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v == -1 {
						return
					}
					if v < 10000 || v > 86400000 {
						errs = append(errs, fmt.Errorf("%q must be -1 (unlimited) or between 10000 and 86400000 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"vm_memory_high_watermark": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Optional:    true,
				Description: "When the server will enter memory based flow-control as relative to the maximum available memory.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(float64)
					if v < 0.4 || v > 0.9 {
						errs = append(errs, fmt.Errorf("%q must be between 0.4 and 0.9 inclusive, got: %v", key, v))
					}
					return
				},
			},
			"queue_index_embed_msgs_below": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
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
				Computed:    true,
				Optional:    true,
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
				Computed: true,
				Optional: true,
				Description: "Log level for the logger used for log integrations and the CloudAMQP Console log view. " +
					"Does not affect the file logger. Requires a RabbitMQ restart to be applied.",
				ValidateFunc: validateLogLevel(),
			},
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Configurable sleep time in seconds between retries for RabbitMQ configuration",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "Configurable timeout time in seconds for RabbitMQ configuration",
			},
		},
	}
}

func resourceRabbitMqConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := rabbitMqConfigurationWriteAttributeKeys()
	params := make(map[string]interface{})
	for _, k := range keys {
		v := d.Get(k)
		if v == nil || v == 0 || v == 0.0 || v == "" {
			continue
		} else if k == "connection_max" {
			if v == -1 {
				v = "infinity"
			}
		} else if k == "consumer_timeout" {
			if v == -1 {
				v = "false"
			}
		} else if k == "log_exchange_level" {
			k = "log.exchange.level"
		}
		params["rabbit."+k] = v
	}
	err := api.UpdateRabbitMqConfiguration(d.Get("instance_id").(int), params, d.Get("sleep").(int), d.Get("timeout").(int))
	if err != nil {
		return err
	}
	id := strconv.Itoa(d.Get("instance_id").(int))
	d.SetId(id)
	return resourceRabbitMqConfigurationRead(d, meta)
}

func resourceRabbitMqConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())
	data, err := api.ReadRabbitMqConfiguration(instanceID, d.Get("sleep").(int), d.Get("timeout").(int))
	log.Printf("[DEBUG] cloudamqp::resource::rabbitmq_configuration::read data: %v", data)
	if err != nil {
		return err
	}

	d.Set("instance_id", instanceID)
	for k, v := range data {
		if validateRabbitMqConfigurationJSONField(k) {
			key := strings.ReplaceAll(k, "rabbit.", "")
			if key == "connection_max" {
				if v == "infinity" || v == nil {
					v = -1
				}
			} else if key == "consumer_timeout" {
				if v == "false" {
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

func resourceRabbitMqConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := rabbitMqConfigurationWriteAttributeKeys()
	params := make(map[string]interface{})
	for _, k := range keys {
		v := d.Get(k)
		if v == nil {
			continue
		} else if k == "connection_max" {
			if v == -1 {
				v = "infinity"
			}
		} else if k == "consumer_timeout" {
			if v == -1 {
				v = "false"
			}
		} else if k == "log_exchange_level" {
			k = "log.exchange.level"
		}
		params["rabbit."+k] = v
	}
	err := api.UpdateRabbitMqConfiguration(d.Get("instance_id").(int), params, d.Get("sleep").(int), d.Get("timeout").(int))
	if err != nil {
		return err
	}
	return resourceRabbitMqConfigurationRead(d, meta)
}

func resourceRabbitMqConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func validateRabbitMqConfigurationJSONField(key string) bool {
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

func rabbitMqConfigurationWriteAttributeKeys() []string {
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
