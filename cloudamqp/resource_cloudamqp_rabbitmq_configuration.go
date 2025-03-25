package cloudamqp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRabbitMqConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRabbitMqConfigurationUpdate,
		ReadContext:   resourceRabbitMqConfigurationRead,
		UpdateContext: resourceRabbitMqConfigurationUpdate,
		DeleteContext: resourceRabbitMqConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
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
				Description: "Size in bytes below which to embed messages in the queue index. 0 will turn off payload embedding in the queue index.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 10485760 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 10485760 inclusive, got: %d", key, v))
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
				ValidateDiagFunc: validateLogLevel(),
			},
			"cluster_partition_handling": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				Description:  "Set how the cluster should handle network partition.",
				ValidateFunc: validation.StringInSlice([]string{"autoheal", "pause_minority", "ignore"}, true),
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

func resourceRabbitMqConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		api           = meta.(*api.API)
		instanceID, _ = strconv.Atoi(d.Id())
		sleep         = d.Get("sleep").(int)
		timeout       = d.Get("timeout").(int)
	)

	// Set arguments during import
	if d.Get("instance_id").(int) == 0 {
		d.Set("instance_id", instanceID)
	}
	if d.Get("sleep").(int) == 0 && d.Get("timeout").(int) == 0 {
		d.Set("sleep", 60)
		d.Set("timeout", 3600)
	}

	data, err := api.ReadRabbitMqConfiguration(ctx, instanceID, sleep, timeout)
	log.Printf("[DEBUG] cloudamqp::resource::rabbitmq_configuration::read data: %v", data)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range data {
		if validateRabbitMqConfigurationJSONField(k) {
			if v == nil || v == "" {
				continue
			}
			key := strings.ReplaceAll(k, "rabbit.", "")
			if key == "connection_max" {
				if v == "infinity" {
					v = -1
				}
			} else if key == "consumer_timeout" {
				if v == "false" {
					v = -1
				}
			} else if key == "log.exchange.level" {
				key = "log_exchange_level"
			}
			d.Set(key, v)
		}
	}
	return diag.Diagnostics{}
}

func resourceRabbitMqConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		keys       = rabbitMqConfigurationWriteAttributeKeys()
		params     = make(map[string]interface{})
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	for _, k := range keys {
		if !d.HasChange(k) {
			continue
		}

		v := d.Get(k)
		if v == nil || v == "" {
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

	log.Printf("[DEBUG] RabbitMQ configuration params: %v", params)
	if len(params) > 0 {
		err := api.UpdateRabbitMqConfiguration(ctx, instanceID, params, sleep, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.Itoa(instanceID))
	return resourceRabbitMqConfigurationRead(ctx, d, meta)
}

func resourceRabbitMqConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
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
		"rabbit.log.exchange.level",
		"rabbit.cluster_partition_handling":
		return true
	}
	return false
}

func validateLogLevel() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"debug",
		"info",
		"warning",
		"error",
		"critical",
		"critical",
		"none",
	}, true))
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
		"cluster_partition_handling",
	}
}
