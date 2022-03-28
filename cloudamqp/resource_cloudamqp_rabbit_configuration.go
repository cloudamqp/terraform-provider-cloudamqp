package cloudamqp

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Computed:    true,
				Optional:    true,
				Description: "Set the server AMQP 0-9-1 heartbeat timeout in seconds.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be greater then 0, got: %d", key, v))
					}
					return
				},
			},
			"channel_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Set the maximum permissible number of channels per connection.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be greater then 0, got: %d", key, v))
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
					if v < 10000 || v > 25000000 {
						errs = append(errs, fmt.Errorf("%q must be between 10000 and 25000000 inclusive, got: %d", key, v))
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
						errs = append(errs, fmt.Errorf("%q must be between 0.4 and 0.9 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"queue_index_embeded_msgs_below": {
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
		},
	}
}

func resourceRabbitConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	// api := meta.(*api.API)
	// instanceID, _ := strconv.Atoi(d.Id())
	// data, err := api.ReadRabbitConfiguration(instanceID)
	// log.Printf("[DEBUG] cloudamqp::resource::rabbit_configuration::read data: %v", data)
	// if err != nil {
	// 	return err
	// }
	// d.Set("instance_id", instanceID)
	// for k, v := range data {
	// 	if v == nil {
	// 		continue
	// 	}
	// 	if validateRabbitConfigurationJsonField(k) {
	// 		key := strings.ReplaceAll(k, "rabbit.", "")
	// 		d.Set(key, v)
	// 	}
	// }
	return nil
}

func resourceRabbitConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	//api := meta.(*api.API)
	keys := []string{"heartbeat", "channel_max", "consumer_timeout", "vm_memory_high_watermark", "queue_index_embeded_msgs_below", "max_message_size"}
	params := make(map[string]interface{})
	for _, k := range keys {
		log.Printf("[DEBUG] cloudamqp::resource::rabbit_configuration::update k: %v, value: %v", k, d.Get(k))
		if v := d.Get(k); v != nil {
			params["rabbit."+k] = d.Get(k)
			if k == "queue_index_embeded_msgs_below" || k == "max_message_size" {
				if v == 0 {
					params["rabbit."+k] = nil
				}
			}
		}
	}
	log.Printf("[DEBUG] cloudamqp::resource::rabbit_configuration::update instance id: %v, params: %v", d.Get("instance_id"), params)
	// err := api.UpdateRabbitConfiguration(d.Get("instance_id").(int), params)
	// if err != nil {
	// 	return err
	// }
	return resourceRabbitConfigurationRead(d, meta)
}

func resourceRabbitConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	// api := meta.(*api.API)
	// err := api.DeleteRabbitConfiguration()
	// return err
	return nil
}

func validateRabbitConfigurationJsonField(key string) bool {
	switch key {
	case "rabbit.heartbeat",
		"rabbit.channel_max",
		"rabbit.consumer_timeout",
		"rabbit.vm_memory_high_watermark",
		"rabbit.queue_index_embeded_msgs_below",
		"rabbit.max_message_size":
		return true
	}
	return false
}
