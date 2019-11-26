package cloudamqp

import (
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceSecurityFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityFirewallCreate,
		Read:   resourceSecurityFirewallRead,
		Update: resourceSecurityFirewallUpdate,
		Delete: resourceSecurityFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"services": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateServices(),
							},
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntBetween(0, 65554),
							},
						},
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSecurityFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	var params []map[string]interface{}
	localFirewalls := d.Get("rules").(*schema.Set).List()

	for _, k := range localFirewalls {
		params = append(params, k.(map[string]interface{}))
	}

	instance_id := d.Get("instance_id").(int)
	err := api.CreateFirewallSettings(instance_id, params)
	d.SetId(strconv.Itoa(instance_id))
	return err
}

func resourceSecurityFirewallRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	instance_id, _ := strconv.Atoi(d.Id())
	data, err := api.ReadFirewallSettings(instance_id)
	if err != nil {
		return err
	}
	d.Set("instance_id", instance_id)
	d.Set("rules", data)

	return nil
}

func resourceSecurityFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	var params []map[string]interface{}
	localFirewalls := d.Get("rules").(*schema.Set).List()
	for _, k := range localFirewalls {
		params = append(params, k.(map[string]interface{}))
	}
	err := api.UpdateFirewallSettings(d.Get("instance_id").(int), params)
	return err
}

func resourceSecurityFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	err := api.DeleteFirewallSettings(d.Get("instance_id").(int))
	return err
}

func validateServices() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"AMQP",
		"AMQPS",
		"MQTT",
		"MQTTS",
		"STOMP",
		"STOMPS",
	}, true)
}
